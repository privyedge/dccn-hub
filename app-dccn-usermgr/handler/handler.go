package handler

import (
	"context"
	"errors"
	"log"
	"strings"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	"github.com/google/uuid"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
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

func (p *UserHandler) Register(ctx context.Context, user *usermgr.User, rsp *common_proto.Error) error {

	log.Println("Debug Register")
	// verify email and password
	if !matchPattern(OpUserNameMatch, user.Name) || !matchPattern(OpPasswordMatch, user.Password) || !matchPattern(OpEmailMatch, user.Email) {
		err := errors.New("name or email or password invalid")
		log.Println(err.Error())
		return err
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	_, dbErr := p.db.Get(strings.ToLower(user.Email))
	if dbErr == nil {
		log.Println(dbErr.Error())
		return dbErr
	}

	authorizationToken, err := p.token.NewAuthorizationToken(user.Email)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// sends activate email
	if err := p.pubEmail.Publish(context.Background(), &common_proto.MailEvent{ToAddress: user.Email, Body: authorizationToken}); err != nil {
		log.Println(err.Error())
		return err
	}

	user.Password = string(hashedPwd)
	user.Email = strings.ToLower(user.Email)
	user.Id = uuid.New().String()
	if err := p.db.Create(user); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) ForgetPassword(ctx context.Context, req *usermgr.AskResetPasswordRequest, rsp *common_proto.Error) error {
	log.Println("Debug AskResetPassword")

	authorizationToken, err := p.token.NewAuthorizationToken(req.Email)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// some logic to examine and verify user here

	// sends activate email
	if err := p.pubEmail.Publish(context.Background(), &common_proto.MailEvent{ToAddress: req.Email, Body: authorizationToken}); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *UserHandler) ConfirmPassword(ctx context.Context, req *usermgr.ResetPasswordRequest, rsp *common_proto.Error) error {

	log.Println("Debug ResetPassword")

	// verify code if is expired
	payload, err := p.token.VerifyAuthorizationToken(req.Token)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if payload.Email != req.Email {
		err = errors.New("user invalid")
		log.Println(err.Error())
		return err
	}

	// encrypt password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// update password. if not exist, db return not found
	if err := p.db.UpdatePassword(req.Email, hashedPwd); err != nil {
		log.Println(err.Error())
		return err
	}

	// TODO: remove user cache token here
	// change secret, otherwise token in caches
	// need redesign token

	return nil
}

func (p *UserHandler) ConfirmRegistration(ctx context.Context, req *usermgr.ActivateRequest, rsp *common_proto.Error) error {

	log.Println("Debug Activate")

	// verify code if is expired
	payload, err := p.token.VerifyAuthorizationToken(req.Token)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// update password. if not exist, db return not found
	if err := p.db.UpdateActivateStatus(payload.Email); err != nil {
		log.Println(err.Error())
		return err
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
		err = errors.New("invalid password")
		log.Println(err.Error())
		return err
	}

	rsp.Token, err = p.token.NewAuthenticationToken(user)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	*rsp.User = *user

	// used to refresh token
	p.blacklist.Add(rsp.Token)
	return nil
}

func (p *UserHandler) Logout(ctx context.Context, in *usermgr.LogoutRequest, out *common_proto.Error) error {

	log.Println("Debug into Logout")
	md, ok := metadata.FromContext(ctx)
	if !ok {
		log.Println("no auth meta-data found in request")
		return errors.New("no auth meta-data found in request")
	}

	p.blacklist.Remove(md["token"])
	return nil
}

func (p *UserHandler) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.NewTokenResponse) error {

	log.Println("Debug into NewToken")
	req, err := p.db.Get(strings.ToLower(req.Email))
	if err != nil {
		log.Println(err.Error())
		return err
	}

	rsp.Token, err = p.token.NewAuthenticationToken(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

	log.Println("Debug into VerifyToken: ", req.Token)
	if !p.blacklist.Available(req.Token) {
		err := errors.New("token is unavailable")
		log.Println(err.Error())
		return err
	}

	if _, err := p.token.VerifyAuthenticationToken(req.Token); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *UserHandler) VerifyAndRefreshToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

	log.Println("Debug into VerifyAndRefreshToken: ", req.Token)

	//for token reflesh
	if !p.blacklist.Available(req.Token) {
		err := errors.New("token is unavailable")
		log.Println(err.Error())
		return err
	}

	_, err := p.token.VerifyAuthenticationToken(req.Token)
	if err == nil || (err != nil && !p.blacklist.Available(req.Token)) {
		p.blacklist.Refresh(req.Token)
		return nil
	}

	p.blacklist.Remove(req.Token)
	log.Println(err.Error())
	return err
}

func (p *UserHandler) RefreshToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {
	p.blacklist.Refresh(req.Token)
	return nil
}

func (p *UserHandler) Destroy() {
	p.blacklist.destroy()
}
