package handler

import (
	"context"
	"log"
	"strings"
	"time"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	user_util "github.com/Ankr-network/dccn-hub/app-dccn-usermgr/util"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/token"
)

type UserHandler struct {
	db        dbservice.DBService // db
	token     token.IToken        // token interface
	pubEmail  micro.Publisher
	blacklist *Blacklist // used for logout
}

func New(dbService dbservice.DBService, tokenService token.IToken, pubEmail micro.Publisher) *UserHandler {
	return &UserHandler{
		db:        dbService,
		token:     tokenService,
		pubEmail:  pubEmail,
		blacklist: NewBlacklist(),
	}
}

func (p *UserHandler) Register(ctx context.Context, req *usermgr.RegisterRequest, rsp *common_proto.Empty) error {

	log.Println("Debug Register")
	user := req.User

	// verify email and password
	if err := user_util.CheckRegister(user.Attributes.Name, user.Email, req.Password); err != nil {
		log.Println(err.Error())
		return err
	}

	// we store the hashed password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// check if email exists already
	_, dbErr := p.db.GetUserByEmail(strings.ToLower(user.Email))
	if dbErr == nil {
		log.Println(ankr_default.ErrEmailExit)
		return ankr_default.ErrEmailExit
	}

	hashPassword := string(hashedPwd)
	user.Email = strings.ToLower(user.Email)
	user.Id = uuid.New().String()
	user.Status = usermgr.UserStatus_CONFIRMING

	if user.Attributes.Name == "ankrtest" { // for debug
		user.Status = usermgr.UserStatus_CONFIRMED

	} else {
		_, confirmRegistrationCode, err := p.token.NewToken(user.Id, false)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		e := &mail.MailEvent{
			Type: mail.EmailType_CONFIRM_REGISTRATION,
			From: ankr_default.NoReplyEmailAddress,
			To:   []string{user.Email},
			OpMail: &mail.MailEvent_ConfirmRegistration{
				ConfirmRegistration: &mail.ConfirmRegistration{
					UserName: user.Attributes.Name,
					UserId:   user.Id,
					Code:     confirmRegistrationCode,
				},
			},
		}

		if err := p.pubEmail.Publish(context.Background(), e); err != nil {
			log.Println(err.Error())
			return err
		}

	}

	if err := p.db.CreateUser(user, hashPassword); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) ConfirmRegistration(ctx context.Context, req *usermgr.ConfirmRegistrationRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ConfirmRegistration")

	if err := user_util.CheckEmail(req.Email); err != nil {
		log.Println(err.Error())
		return err
	}

	// verify code if is expired
	_, err := p.token.Verify(req.ConfirmationCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	attr := []*usermgr.UserAttribute{
		{
			Key: "Status", Value: &usermgr.UserAttribute_IntValue{
				IntValue: int64(usermgr.UserStatus_CONFIRMED),
			},
		},
	}
	// update password. if not exist, db return not found
	if err := p.db.UpdateUserByEmail(req.Email, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {

	req.Email = strings.ToLower(req.Email)
	log.Println("Debug Login")

	user, err := p.db.GetUserByEmail(req.Email)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		log.Println(ankr_default.ErrPasswordError.Error())
		return ankr_default.ErrPasswordError
	}

	expired, userToken, err := p.token.NewToken(user.ID, false)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	rsp.AuthenticationResult = &usermgr.AuthenticationResult{}
	rsp.User = &usermgr.User{}
	rsp.User.Attributes = &usermgr.UserAttributes{}
	rsp.User.Id = user.ID
	rsp.User.Email = user.Email
	rsp.User.Attributes.Name = user.Name
	rsp.User.Attributes.CreationDate = user.CreationDate
	rsp.User.Attributes.LastModifiedDate = user.LastModifiedDate

	rsp.AuthenticationResult.AccessToken = userToken
	rsp.AuthenticationResult.Expiration = uint64(expired)
	rsp.AuthenticationResult.IssuedAt = uint64(time.Now().Unix())
	_, refreshToken, _ := p.token.NewToken(user.ID, true)
	rsp.AuthenticationResult.RefreshToken = refreshToken

	attr := []*usermgr.UserAttribute{
		{
			Key:   "Token",
			Value: &usermgr.UserAttribute_StringValue{StringValue: refreshToken},
		},
	}

	if err := p.db.UpdateUser(user.ID, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) Logout(ctx context.Context, req *usermgr.RefreshToken, out *common_proto.Empty) error {

	uid, err := user_util.GetIdFromToken(req.RefreshToken)
	if err != nil {
		return err
	}

	user, err := p.db.GetUser(uid)

	attr := []*usermgr.UserAttribute{
		{
			Key:   "Token",
			Value: &usermgr.UserAttribute_StringValue{StringValue: ""},
		},
	}

	if err := p.db.UpdateUser(user.ID, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) RefreshSession(ctx context.Context, req *usermgr.RefreshToken, rsp *usermgr.AuthenticationResult) error {
	uid, err := user_util.GetIdFromToken(req.RefreshToken)
	if err != nil {
		return err
	}

	user, err := p.db.GetUser(uid)

	if user.Token != req.RefreshToken {
		return ankr_default.ErrRefreshToken
	}

	expired, newToken, err2 := p.token.NewToken(user.ID, false)

	if err2 != nil {
		log.Println(err2.Error())
		return err2
	}

	rsp.AccessToken = newToken
	rsp.Expiration = uint64(expired)
	rsp.IssuedAt = uint64(time.Now().Unix())
	_, refreshToken, _ := p.token.NewToken(user.ID, true)
	rsp.RefreshToken = refreshToken

	attr := []*usermgr.UserAttribute{
		{
			Key:   "Token",
			Value: &usermgr.UserAttribute_StringValue{StringValue: refreshToken},
		},
	}

	if err := p.db.UpdateUser(user.ID, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) VerifyAccessToken(ctx context.Context, req *common_proto.Empty, rsp *common_proto.Empty) error {
	meta, ok := metadata.FromContext(ctx)
	// Note this is now uppercase (not entirely sure why this is...)
	var accessToken string
	if ok {
		accessToken = meta["token"]
	}

	log.Printf("find token %s \n", accessToken)

	if len(accessToken) == 0 {
		return ankr_default.ErrTokenParseFailed
	}
	_, err := user_util.VerifyAccessToken(accessToken)

	if err != nil {
		return err
	}
	return nil

}

func (p *UserHandler) Destroy() {
	p.blacklist.destroy()
}

func (p *UserHandler) ForgotPassword(ctx context.Context, req *usermgr.ForgotPasswordRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ForgetPassword")
	// generate new authorization token for reset password
	_, forgetPasswordCode, err := p.token.NewToken(req.Email, false)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// some logic to examine and verify user here

	e := &mail.MailEvent{
		Type: mail.EmailType_FORGET_PASSWORD,
		From: ankr_default.NoReplyEmailAddress,
		To:   []string{req.Email},
		OpMail: &mail.MailEvent_ForgetPassword{
			ForgetPassword: &mail.ForgetPassword{
				Email: req.Email,
				Code:  forgetPasswordCode,
			},
		},
	}
	// sends activate email
	if err := p.pubEmail.Publish(context.Background(), e); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) ConfirmPassword(ctx context.Context, req *usermgr.ConfirmPasswordRequest, rsp *common_proto.Empty) error {
	log.Println("Debug ConfirmPassword")

	if err := user_util.CheckPassword(req.NewPassword); err != nil {
		log.Println(err.Error())
		return err
	}

	// verify code if is expired
	_, err := p.token.Verify(req.ConfirmationCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// check email auth
	email, err := user_util.GetIdFromToken(req.ConfirmationCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if strings.ToLower(email) != strings.ToLower(req.Email) || email == "" {
		log.Println(ankr_default.ErrAuthNotAllowed)
		return ankr_default.ErrAuthNotAllowed
	}

	// new password should not same as before
	if record, err := p.db.GetUserByEmail(strings.ToLower(email)); err != nil {
		log.Println(err.Error())
		return err
	} else if err := bcrypt.CompareHashAndPassword([]byte(record.HashedPassword), []byte(req.NewPassword)); err == nil {
		log.Println(ankr_default.ErrPasswordSame)
		return ankr_default.ErrPasswordSame
	}

	// hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// update password. if not exist, db return not found
	attr := []*usermgr.UserAttribute{
		{
			Key:   "HashedPassword",
			Value: &usermgr.UserAttribute_StringValue{StringValue: string(hashedPwd)},
		},
	}

	if err := p.db.UpdateUserByEmail(strings.ToLower(req.Email), attr); err != nil {
		log.Println(err.Error())
		return err
	}

	// TODO: remove user cache token here
	// change secret, otherwise token in caches
	// need redesign token

	return nil
}

func (p *UserHandler) ChangePassword(ctx context.Context, req *usermgr.ChangePasswordRequest, rsp *common_proto.Empty) error {
	uid := user_util.GetUserID(ctx)
	log.Println("Debug ChangePassword")

	if err := user_util.CheckPassword(req.NewPassword); err != nil {
		log.Println(err.Error())
		return err
	}

	if req.NewPassword == req.OldPassword {
		log.Println(ankr_default.ErrPasswordSame)
		return ankr_default.ErrPasswordSame
	}

	// hash password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	attr := []*usermgr.UserAttribute{
		{
			Key:   "HashedPassword",
			Value: &usermgr.UserAttribute_StringValue{StringValue: string(hashedPwd)},
		},
	}

	// update password. if not exist, db return not found
	if err := p.db.UpdateUser(uid, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) UpdateAttributes(ctx context.Context, req *usermgr.UpdateAttributesRequest, rsp *usermgr.User) error {
	uid := user_util.GetUserID(ctx)
	log.Println("Debug UpdateAttributes")

	if err := p.db.UpdateUser(uid, req.UserAttributes); err != nil {
		log.Println(err.Error())
		return err
	}

	userRecord, err := p.db.GetUser(uid)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	rsp.Id = userRecord.ID
	rsp.Email = userRecord.Email
	rsp.Attributes = &usermgr.UserAttributes{
		Name:             userRecord.Name,
		CreationDate:     userRecord.CreationDate,
		LastModifiedDate: userRecord.LastModifiedDate,
		PubKey:           userRecord.PubKey,
	}
	rsp.Status = userRecord.Status

	return nil
}

func (p *UserHandler) ChangeEmail(ctx context.Context, req *usermgr.ChangeEmailRequest, rsp *common_proto.Empty) error {

	uid := user_util.GetUserID(ctx)
	req.NewEmail = strings.ToLower(req.NewEmail)
	log.Println("Debug ChangeEmail")

	if err := user_util.CheckEmail(req.NewEmail); err != nil {
		log.Println(ankr_default.ErrEmailFormat)
		return ankr_default.ErrEmailFormat
	}

	if userRecord, err := p.db.GetUser(uid); err != nil {
		log.Println(err.Error())
		return err
	} else if userRecord.Email == req.NewEmail {
		log.Println(ankr_default.ErrEmailSame)
		return ankr_default.ErrEmailSame
	}

	_, changeEmailCode, err := p.token.NewToken(uid, false)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err := p.pubEmail.Publish(context.TODO(),
		&mail.MailEvent{
			Type: mail.EmailType_CONFIRM_EMAIL,
			From: ankr_default.NoReplyEmailAddress,
			To:   []string{req.NewEmail},
			OpMail: &mail.MailEvent_ChangeEmail{
				ChangeEmail: &mail.ChangeEmail{
					UserId:   uid,
					NewEmail: req.NewEmail,
					Code:     changeEmailCode,
				},
			},
		}); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) ConfirmEmail(ctx context.Context, req *usermgr.ConfirmEmailRequest, rsp *common_proto.Empty) error {

	uid := user_util.GetUserID(ctx)
	log.Println("Debug ChangeEmail")

	if _, err := p.token.Verify(req.ConfirmationCode); err != nil {
		log.Println(err.Error())
		return err
	}

	attr := []*usermgr.UserAttribute{
		{
			Key:   "Email",
			Value: &usermgr.UserAttribute_StringValue{StringValue: strings.ToLower(req.NewEmail)},
		},
	}

	tokenUid, err := user_util.GetIdFromToken(req.ConfirmationCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if uid != tokenUid {
		log.Println(ankr_default.ErrAuthNotAllowed)
		return ankr_default.ErrAuthNotAllowed
	}

	if err := p.db.UpdateUser(uid, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
