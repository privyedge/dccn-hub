package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/micro/go-micro/metadata"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/token"
)

type UserHandler struct {
	db        dbservice.DBService // db
	token     token.IToken        // token interface
	blacklist *Blacklist          // used for logout
}

func New(dbService dbservice.DBService, tokenService token.IToken) *UserHandler {
	return &UserHandler{
		db:        dbService,
		token:     tokenService,
		blacklist: NewBlacklist(),
	}
}

type Token struct {
	Exp int64
	Jti string
	Iss string
}


func getUserIDByRefreshToken(refresh_token string) (string, error){
	parts := strings.Split(refresh_token, ".")

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	var dat Token

	if err := json.Unmarshal(decoded, &dat); err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	return string(dat.Jti), nil
}


func VarifyAccessToken(refresh_token string) (string, error){
	parts := strings.Split(refresh_token, ".")

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		return "", ankr_default.ErrTokenParseFailed
	}

	var dat Token

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
	email_error := ValidateEmailFormat(user.Email)
	if email_error != nil {
		log.Println(email_error.Error())
		return email_error
	}

	if len(req.Password) < 6 {
		return ankr_default.ErrPasswordLength
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return errors.New("GenerateFromPassword error")
	}

	_, dbErr := p.db.Get(strings.ToLower(user.Email))
	if dbErr == nil {
		log.Println("email exist")
		return ankr_default.ErrEmailExit
	}

	hashPassword := string(hashedPwd)
	user.Email = strings.ToLower(user.Email)
	user.Id = uuid.New().String()
	if err := p.db.Create(user, hashPassword); err != nil {
		log.Println(err.Error())
		return errors.New("data add user error")
	}
	return nil
}





func (p *UserHandler) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {

	log.Println("Debug Login")
	user, err := p.db.Get(strings.ToLower(req.Email))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Compares our given password against the hashed password
	// stored in the database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		err = ankr_default.ErrPasswordError
		log.Println(err.Error())
		return ankr_default.ErrPasswordError
	}

	log.Printf("user token %+v", user.Tokens)

	if len(user.Tokens)  > 100 { //todo just for test
		return ankr_default.ErrTokenPassedMax
	}



	expired, token, err2 := p.token.NewToken(user.ID)

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
	rsp.User.Attributes.CreationDate = user.Creation_date
	rsp.User.Attributes.LastModifiedDate = user.Last_modified_date


	rsp.AuthenticationResult.AccessToken = token
	rsp.AuthenticationResult.Expiration = uint64(expired)
	rsp.AuthenticationResult.IssuedAt = uint64(time.Now().Unix())
	_, refresh_token , _ := p.token.NewToken(user.ID)
	rsp.AuthenticationResult.RefreshToken = refresh_token

    tokens := append(user.Tokens, refresh_token)

    p.db.UpdateToken(user.ID, tokens)

	//
	////for token reflesh
	//p.blacklist.Add(rsp.Token)
	return nil
}


func (p *UserHandler) Logout(ctx context.Context, req *usermgr.RefreshToken , out *common_proto.Empty) error {

	log.Println("Debug into Logout")
	uid , err := getUserIDByRefreshToken(req.RefreshToken)
	if err != nil {
		return err
	}

	user, err := p.db.GetUserByID(uid)

	if !containsToken(user.Tokens, req.RefreshToken){
		return ankr_default.ErrRefreshToken
	}

	index := indexOfToken(user.Tokens, req.RefreshToken)

	newTokens := append(user.Tokens[:index], user.Tokens[index+1:]...)
	p.db.UpdateToken(user.ID, newTokens)

	return nil
}

// func (p *UserHandler) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.NewTokenResponse) error {

// 	log.Println("Debug into NewToken")
// 	req, err := p.db.Get(strings.ToLower(req.Email))
// 	if err != nil {
// 		log.Println(err.Error())
// 		return err
// 	}

// 	rsp.Token, err = p.token.New(req)
// 	if err != nil {VerifyToken
// 		log.Println(err.Error())
// 		return err
// 	}

// 	return nil
// }

///
//func (p *UserHandler) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {
//
//	log.Println("Debug into VerifyToken: ", req.Token)
//	if !p.blacklist.Available(req.Token) {
//		err := ankr_default.ErrTokenNeedRefresh
//		rsp.Status = common_proto.Status_FAILURE
//		rsp.Details = ankr_default.ErrTokenNeedRefresh.Error()
//		log.Println(err.Error())
//		return err
//	}
//
//	if _, err := p.token.Verify(req.Token); err != nil {
//		log.Println(err.Error())
//		return err
//	}
//	return nil
//}

//func (p *UserHandler) VerifyAndRefreshToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {
//
//	log.Println("Debug into VerifyAndRefreshToken: ", req.Token)
//
//	//for token reflesh
//	if !p.blacklist.Available(req.Token) {
//		rsp.Status = common_proto.Status_FAILURE
//		rsp.Details = ankr_default.ErrTokenNeedRefresh.Error()
//		err := ankr_default.ErrTokenNeedRefresh
//		log.Println(err.Error())
//		return err
//	}
//
//	_, err := p.token.Verify(req.Token)
//	if err == nil || (err != nil && !p.blacklist.Available(req.Token)) {
//		p.blacklist.Refresh(req.Token)
//		return nil
//	}
//
//	p.blacklist.Remove(req.Token)
//	log.Println(err.Error())
//	return err
//}

func containsToken(tokens []string, token string) bool {
	for _, t := range tokens {
		if t == token {
			return true
		}
	}
	return false
}

func indexOfToken(tokens []string, token string) (int) {
	for k, v := range tokens {
		if token == v {
			return k
		}
	}
	return -1    //not found.
}


func (p *UserHandler) RefreshSession(ctx context.Context, req *usermgr.RefreshToken, rsp *usermgr.AuthenticationResult) error {
    uid , err := getUserIDByRefreshToken(req.RefreshToken)
    if err != nil {
    	return err
	}

	user, err := p.db.GetUserByID(uid)

	if !containsToken(user.Tokens, req.RefreshToken){
        return ankr_default.ErrRefreshToken
	}

	expired, token, err2 := p.token.NewToken(user.ID)

	if err2 != nil {
		log.Println(err2.Error())
		return err2
	}




	//rsp = usermgr.AuthenticationResult{}

	rsp.AccessToken = token
	rsp.Expiration = uint64(expired)
	rsp.IssuedAt = uint64(time.Now().Unix())
	_, refresh_token , _ := p.token.NewToken(user.ID)
	rsp.RefreshToken = refresh_token

	tokens := append(user.Tokens, refresh_token)

	index := indexOfToken(user.Tokens, req.RefreshToken)
	user.Tokens[index] = refresh_token
	p.db.UpdateToken(user.ID, tokens)

	return nil
}


func (p *UserHandler)  ConfirmRegistration(ctx context.Context, req *usermgr.ConfirmRegistrationRequst, rsp *common_proto.Empty) error{
	//todo
	return nil
}

func (p *UserHandler)  ForgotPassword(ctx context.Context, req *usermgr.ForgotPasswordRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *UserHandler)  ConfirmPassword(ctx context.Context, req *usermgr.ConfirmPasswordRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *UserHandler) ChangePasword(ctx context.Context, req *usermgr.ChangePasswordRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *UserHandler) UpdateAttributes(ctx context.Context, req *usermgr.UpdateAttributesRequest, rsp *usermgr.User) error {
	//todo
	return nil
}


func (p *UserHandler) ChangeEmail(ctx context.Context, req *usermgr.ChangeEmailRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *UserHandler) VerifyEmail(ctx context.Context, req *usermgr.VerifyEmailRequst, rsp *usermgr.User) error {

	return errors.New("acccess_token ok")
}


func (p *UserHandler) VerifyAccessToken(ctx context.Context, req *common_proto.Empty, rsp *common_proto.Empty) error {
	meta, ok := metadata.FromContext(ctx)
	// Note this is now uppercase (not entirely sure why this is...)
	var access_token string
	if ok  {
		access_token = meta["token"]
	}

	log.Printf("find token %s \n", access_token)

	if len(access_token) == 0 {
		return ankr_default.ErrTokenParseFailed
	}
	_ , err := VarifyAccessToken(access_token)

	if err != nil {
		return err
	}
    return nil

}


func ValidateEmailFormat(email string) error{
	emailRegexp := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegexp.MatchString(email) {
		return ankr_default.ErrEmailFormat
	}
	return nil
}





func (p *UserHandler) Destroy() {
	p.blacklist.destroy()
}
