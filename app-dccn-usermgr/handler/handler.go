package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"regexp"
	"strings"
	"time"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	ankr_util "github.com/Ankr-network/dccn-common/util"
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

func getUserIDByRefreshToken(refreshToken string) (string, error) {
	parts := strings.Split(refreshToken, ".")
	if len(parts) != 3 {
		return "", ankr_default.ErrTokenParseFailed
	}

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	var dat ankr_util.Token

	if err := json.Unmarshal(decoded, &dat); err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	return string(dat.Jti), nil
}

func VerifyAccessToken(refreshToken string) (string, error) {
	parts := strings.Split(refreshToken, ".")

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	var dat ankr_util.Token

	if err := json.Unmarshal(decoded, &dat); err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	now := time.Now().Unix()

	if now > int64(dat.Exp) {
		return "", ankr_default.ErrTokenParseFailed
	}

	return string(dat.Jti), nil
}

func (p *UserHandler) Register(ctx context.Context, req *usermgr.RegisterRequest, rsp *common_proto.Empty) error {

	log.Println("Debug Register")
	user := req.User
	emailError := ValidateEmailFormat(user.Email)
	if emailError != nil {
		log.Println(emailError.Error())
		return emailError
	}

	if len(req.Password) < 6 {
		return ankr_default.ErrPasswordLength
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.New("GenerateFromPassword error")
	}

	_, dbErr := p.db.GetUserByEmail(strings.ToLower(user.Email))
	if dbErr == nil {
		log.Println("email exist")
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

		if err := p.pubEmail.Publish(context.TODO(), e); err != nil {
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

	log.Println("Debug Login")
	user, err := p.db.GetUserByEmail(strings.ToLower(req.Email))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		err = ankr_default.ErrPasswordError
		log.Println(err.Error())
		return ankr_default.ErrPasswordError
	}

	log.Printf("user token %+v", user.Token)

	//if len(user.Tokens) > 100 { //todo just for test
	//	return ankr_default.ErrTokenPassedMax
	//}

	expired, token, err2 := p.token.NewToken(user.ID, false)

	if err2 != nil {
		log.Println(err2.Error())
		return err2
	}

	rsp.AuthenticationResult = &usermgr.AuthenticationResult{}
	rsp.User = &usermgr.User{}
	rsp.User.Attributes = &usermgr.UserAttributes{}
	rsp.User.Id = user.ID
	rsp.User.Email = user.Email
	rsp.User.Attributes.Name = user.Name
	rsp.User.Attributes.CreationDate = user.CreationDate
	rsp.User.Attributes.LastModifiedDate = user.LastModifiedDate

	rsp.AuthenticationResult.AccessToken = token
	rsp.AuthenticationResult.Expiration = uint64(expired)
	rsp.AuthenticationResult.IssuedAt = uint64(time.Now().Unix())
	_, refreshToken, _ := p.token.NewToken(user.ID, true)
	rsp.AuthenticationResult.RefreshToken = refreshToken

	//tokens := append(user.Tokens, refreshToken)

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

	//
	////for token reflesh
	//p.blacklist.Add(rsp.Token)
	return nil
}

func (p *UserHandler) Logout(ctx context.Context, req *usermgr.RefreshToken, out *common_proto.Empty) error {
	
	uid, err := getUserIDByRefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	user, err := p.db.GetUser(uid)

	//if !containsToken(user.Tokens, req.RefreshToken) {
	//	return ankr_default.ErrRefreshToken
	//}
	//
	//index := indexOfToken(user.Tokens, req.RefreshToken)
	//
	//newTokens := append(user.Tokens[:index], user.Tokens[index+1:]...)
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

func containsToken(tokens []string, token string) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}

func indexOfToken(tokens []string, token string) int {
	for k, v := range tokens {
		if token == v {
			return k
		}
	}
	return -1 //not found.
}

func (p *UserHandler) RefreshSession(ctx context.Context, req *usermgr.RefreshToken, rsp *usermgr.AuthenticationResult) error {
	uid, err := getUserIDByRefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	user, err := p.db.GetUser(uid)

	if user.Token != req.RefreshToken {
		return ankr_default.ErrRefreshToken
	}

	expired, token, err2 := p.token.NewToken(user.ID, false)

	if err2 != nil {
		log.Println(err2.Error())
		return err2
	}

	//rsp = usermgr.AuthenticationResult{}

	rsp.AccessToken = token
	rsp.Expiration = uint64(expired)
	rsp.IssuedAt = uint64(time.Now().Unix())
	_, refreshToken, _ := p.token.NewToken(user.ID, true)
	rsp.RefreshToken = refreshToken

	//tokens := append(user.Tokens, refreshToken)
	//
	//index := indexOfToken(user.Tokens, req.RefreshToken)
	//user.Tokens[index] = refreshToken

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
	var access_token string
	if ok {
		access_token = meta["token"]
	}

	log.Printf("find token %s \n", access_token)

	if len(access_token) == 0 {
		return ankr_default.ErrTokenParseFailed
	}
	_, err := VerifyAccessToken(access_token)

	if err != nil {
		return err
	}
	return nil

}

func ValidateEmailFormat(email string) error {
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return ankr_default.ErrEmailFormat
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

	// verify code if is expired
	_, err := p.token.Verify(req.ConfirmationCode)
	if err != nil {
		log.Println(err.Error())
		return err
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
			Key:   "hashedpassword",
			Value: &usermgr.UserAttribute_StringValue{StringValue: string(hashedPwd)},
		},
	}

	if err := p.db.UpdateUserByEmail(req.Email, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	// TODO: remove user cache token here
	// change secret, otherwise token in caches
	// need redesign token

	return nil
}

func (p *UserHandler) ChangePassword(ctx context.Context, req *usermgr.ChangePasswordRequest, rsp *common_proto.Empty) error {
	uid := ankr_util.GetUserID(ctx)
	log.Println("Debug ChangePassword")

	if len(req.NewPassword) < 6 {
		log.Println("password len invalid")
		return errors.New("password len invalid")
	}

	// hash password, TODO: equal return err
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

func (p *UserHandler) UpdateAttributes(ctx context.Context, req *usermgr.UpdateAttributesRequest, rsq *usermgr.User) error {
	uid := ankr_util.GetUserID(ctx)
	log.Println("Debug UpdateAttributes")

	if err := p.db.UpdateUser(uid, req.UserAttributes); err != nil {
		log.Println(err.Error())
		return err
	}

	// if userRecord, err := p.db.GetUser(uid); err != nil {
	// rsp.Id = userRecord.ID
	// userRecord.Name
	// userRecord.Email
	// userRecord.Status
	// }

	return nil
}

func (p *UserHandler) ChangeEmail(ctx context.Context, req *usermgr.ChangeEmailRequest, rsp *common_proto.Empty) error {
	uid := ankr_util.GetUserID(ctx)
	log.Println("Debug ChangeEmail")

	email_error := ValidateEmailFormat(req.NewEmail)
	if email_error != nil {
		log.Println(email_error.Error())
		return email_error
	}

	attr := []*usermgr.UserAttribute{
		{
			Key:   "email",
			Value: &usermgr.UserAttribute_StringValue{StringValue: string(req.NewEmail)},
		},
	}

	if err := p.db.UpdateUser(uid, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
