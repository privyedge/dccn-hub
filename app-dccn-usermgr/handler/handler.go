package handler

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"

	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/email/v1/micro"
	"github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	ankr_util "github.com/Ankr-network/dccn-common/util"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/token"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/util"
)

type UserHandler struct {
	db        dbservice.DBService // db
	token     token.Token        // token interface
	pubEmail  micro.Publisher
	blacklist *Blacklist // used for logout
}

func New(dbService dbservice.DBService, tokenService token.Token, pubEmail micro.Publisher) *UserHandler {
	return &UserHandler{
		db:        dbService,
		token:     tokenService,
		pubEmail:  pubEmail,
		blacklist: NewBlacklist(),
	}
}


func getIdFromToken(refreshToken string) (string, error) {
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



func (p *UserHandler) Register(ctx context.Context, req *usermgr.RegisterRequest, rsp *common_proto.Empty) error {

	log.Println("Debug new Register")
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
		log.Printf("you should see this ")

	} else {
		_, confirmRegistrationCode, err := p.token.NewToken(user.Id, false)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		log.Printf(">>>>>confirmRegistrationCode for %s   %s \n", user.Email ,confirmRegistrationCode)

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
		return ankr_default.ErrEmailFormat
	}


	user, err := p.db.GetUserByEmail(req.Email)
	if err != nil {
		log.Println(err.Error())
		return ankr_default.ErrEmailNoExit
	}



	// verify code if is expired
	playload, err := p.token.Verify(req.ConfirmationCode)
	if err != nil {
		log.Println(err.Error())
		return err
	}else{
		if playload.Id != user.ID {
			return ankr_default.ErrEmailNoMatch
		}
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

//func (p *UserHandler) ConfirmRegistration(ctx context.Context, req *usermgr.ConfirmRegistrationRequest, rsp *common_proto.Empty) error {
//
//	log.Println("Debug into ConfirmRegistration")
//
//	if !user_util.MatchPattern(user_util.OpEmailMatch, req.Email) {
//		log.Println(ankr_default.ErrEmailFormat)
//		return ankr_default.ErrEmailFormat
//	}
//
//	// verify code if is expired
//	_, err := p.token.Verify(req.ConfirmationCode)
//	if err != nil {
//		log.Println(err.Error())
//		return err
//	}
//
//	attr := []*usermgr.UserAttribute{
//		{
//			Key: "Status", Value: &usermgr.UserAttribute_IntValue{
//				IntValue: int64(usermgr.UserStatus_CONFIRMED),
//			},
//		},
//	}
//	// update password. if not exist, db return not found
//	if err := p.db.UpdateUserByEmail(req.Email, attr); err != nil {
//		log.Println(err.Error())
//		return err
//	}
//
//	return nil
//}

//func (p *UserHandler) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {
//
//	log.Println("Debug Login")
//	user, err := p.db.GetUserByEmail(strings.ToLower(req.Email))
//	if err != nil {
//		log.Println(err.Error())
//		return err
//	}
//
//	// Compares our given password against the hashed password
//	// stored in the database
//	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
//		err = ankr_default.ErrPasswordError
//		log.Println(err.Error())
//		return ankr_default.ErrPasswordError
//	}
//
//	log.Printf("user userToken %+v", user.Token)
//
//	expired, userToken, err2 := p.token.NewToken(user.ID, false)
//
//	if err2 != nil {
//		log.Println(err2.Error())
//		return err2
//	}
//
//	rsp.AuthenticationResult = &usermgr.AuthenticationResult{}
//	rsp.User = &usermgr.User{}
//	rsp.User.Attributes = &usermgr.UserAttributes{}
//	rsp.User.Id = user.ID
//	rsp.User.Email = user.Email
//	rsp.User.Attributes.Name = user.Name
//	rsp.User.Attributes.CreationDate = user.CreationDate
//	rsp.User.Attributes.LastModifiedDate = user.LastModifiedDate
//
//	rsp.AuthenticationResult.AccessToken = userToken
//	rsp.AuthenticationResult.Expiration = uint64(expired)
//	rsp.AuthenticationResult.IssuedAt = uint64(time.Now().Unix())
//	_, refreshToken, _ := p.token.NewToken(user.ID, true)
//	rsp.AuthenticationResult.RefreshToken = refreshToken
//
//	attr := []*usermgr.UserAttribute{
//		{
//			Key:   "Token",
//			Value: &usermgr.UserAttribute_StringValue{StringValue: refreshToken},
//		},
//	}
//
//	if err := p.db.UpdateUser(user.ID, attr); err != nil {
//		log.Println(err.Error())
//		return err
//	}
//
//	return nil
//}


func (p *UserHandler) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {

	req.Email = strings.ToLower(req.Email)
	log.Println("Debug Login")

	user, err := p.db.GetUserByEmail(req.Email)
	if err != nil {
		log.Println(err.Error())
		return ankr_default.ErrEmailNoExit
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.Password)); err != nil {
		log.Println(ankr_default.ErrPasswordError.Error())
		return ankr_default.ErrPasswordError
	}

	if user.Status ==  usermgr.UserStatus_CONFIRMING {
         return ankr_default.ErrUserNotVariyEmail
	}

	if user.Status ==  usermgr.UserStatus_DEACTIVATED {
		return ankr_default.ErrUserDeactive
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
	rsp.User.Attributes.PubKey = user.PubKey

	attr2 := []*usermgr.UserAttribute{
		{
			Key:   "AvatarBackgroundColor",
			Value: &usermgr.UserAttribute_IntValue{IntValue: int64(user.AvatarBackgroundColor)},
		},
	}


	rsp.User.Attributes.ExtraFields = attr2

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

	uid, err := getIdFromToken(req.RefreshToken)
	if err != nil {
		return err
	}

	user, err := p.db.GetUser(uid)

	if user.Token != req.RefreshToken {
		return ankr_default.ErrRefreshToken
	}

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
	uid, err := getIdFromToken(req.RefreshToken)
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

	// varify signature
	if _, err := p.token.Verify(accessToken); err != nil {
		log.Println(err.Error())
		return ankr_default.ErrAccessTokenExpired
	}
	return nil

}



func (p *UserHandler) Destroy() {
	p.blacklist.destroy()
}

// forgetPassword = > confirm password
func (p *UserHandler) ForgotPassword(ctx context.Context, req *usermgr.ForgotPasswordRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ForgetPassword")

	_, err := p.db.GetUserByEmail(req.Email)
	if err != nil {
		log.Println("ForgotPassword does not exit " + err.Error())
		return ankr_default.ErrEmailNoExit
	}

	// generate new authorization token for reset password   input email as uid
	forgetPasswordCode, err := p.token.NewTokenWithoutExpired(req.Email)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	log.Printf("ForgetPassword %s for %s \n", forgetPasswordCode , req.Email)

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

func getSha256(value string) string{
	h := sha1.New()
	h.Write([]byte(value))
	sha := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return sha
}

func (p *UserHandler) ConfirmPassword(ctx context.Context, req *usermgr.ConfirmPasswordRequest, rsp *common_proto.Empty) error {
	log.Println("Debug ConfirmPassword")

	if err := user_util.CheckPassword(req.NewPassword); err != nil {
		log.Println(err.Error())
		return err
	}

    var email string

	if playload, err := p.token.Verify(req.ConfirmationCode); err != nil {
		log.Println(err.Error())
		return err
	}else{
		email = playload.Id
		log.Printf("find email %s  new email %s \n", email, req.Email)
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
	}else { // check user status
		if record.Status == usermgr.UserStatus_CONFIRMING {
			// update status to UserStatus_CONFIRMED
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
		}
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
	uid := ankr_util.GetUserID(ctx)
	log.Println("Debug ChangePassword")

	user, err := p.db.GetUser(uid)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.OldPassword)); err != nil {
		log.Println(ankr_default.ErrPasswordError.Error())
		return ankr_default.ErrOldPassword
	}


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
	uid := ankr_util.GetUserID(ctx)
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


	attr2 := []*usermgr.UserAttribute{
		{
			Key:   "AvatarBackgroundColor",
			Value: &usermgr.UserAttribute_IntValue{IntValue: int64(userRecord.AvatarBackgroundColor)},
		},
	}


	rsp.Attributes.ExtraFields = attr2



	rsp.Status = userRecord.Status
	return nil
}

func (p *UserHandler) ChangeEmail(ctx context.Context, req *usermgr.ChangeEmailRequest, rsp *common_proto.Empty) error {
	uid := ankr_util.GetUserID(ctx)
	req.NewEmail = strings.ToLower(req.NewEmail)
	log.Println("Debug ChangeEmail")



	if err := user_util.CheckEmail(req.NewEmail); err != nil {
		log.Println(ankr_default.ErrEmailFormat)
		return ankr_default.ErrEmailFormat
	}


	// new password should not same as before
	if _, err := p.db.GetUserByEmail(strings.ToLower(req.NewEmail)); err != nil {
       // can not find record, it is ok
	}else{
		log.Println("new email have been used")
		return ankr_default.ErrEmailExit
	}


	if userRecord, err := p.db.GetUser(uid); err != nil {
		log.Println(err.Error())
		return err
	} else if userRecord.Email == req.NewEmail {
		log.Println(ankr_default.ErrEmailSame)
		return ankr_default.ErrEmailSame
	}


    // use email as uid
    changeEmailCode, err := p.token.NewTokenWithoutExpired(req.NewEmail)
	log.Printf("------>comfirm code %s  for email %s\n", changeEmailCode, req.NewEmail)
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

	uid := ankr_util.GetUserID(ctx)
	log.Println("Debug ChangeEmail")


	if playload, err := p.token.Verify(req.ConfirmationCode); err != nil {
		log.Println(err.Error())
		return err
	}else{
		    // id as email previous
			if playload.Id != req.NewEmail {
			return ankr_default.ErrEmailNoMatch
		}
	}

	attr := []*usermgr.UserAttribute{
		{
			Key:   "Email",
			Value: &usermgr.UserAttribute_StringValue{StringValue: strings.ToLower(req.NewEmail)},
		},
	}


	if err := p.db.UpdateUser(uid, attr); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
