package handler

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"

	"github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	"github.com/google/uuid"
	"github.com/micro/go-micro/metadata"
	"golang.org/x/crypto/bcrypt"

	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-usermgr/db_service"
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

func (p *UserHandler) Register(ctx context.Context, user *usermgr.User, rsp *common_proto.Error) error {

	log.Println("Debug Register")
	email_error := ValidateEmailFormat(user.Email)
	if email_error != nil {
		log.Println(email_error.Error())
		rsp.Status = common_proto.Status_FAILURE
		rsp.Details = email_error.Error()
		return nil
	}

	if len(user.Password) < 6 {
		rsp.Status = common_proto.Status_FAILURE
		rsp.Details = ankr_default.ErrPasswordLength.Error()
		return nil
	}

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err.Error())
		rsp.Status = common_proto.Status_FAILURE
		return nil
	}

	_, dbErr := p.db.Get(strings.ToLower(user.Email))
	if dbErr == nil {
		log.Println("email exist")
		rsp.Status = common_proto.Status_FAILURE
		rsp.Details = ankr_default.ErrEmailExit.Error()
		return nil
	}

	user.Password = string(hashedPwd)
	user.Email = strings.ToLower(user.Email)
	user.Id = uuid.New().String()
	if err := p.db.Create(user); err != nil {
		log.Println(err.Error())
		rsp.Status = common_proto.Status_FAILURE
		return nil
	}

	rsp.Status = common_proto.Status_SUCCESS
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
		rsp.Error.Status = common_proto.Status_FAILURE
		rsp.Error.Details = ankr_default.ErrPasswordError.Error()
		return err
	}

	rsp.Token, err = p.token.New(user)
	if err != nil {
		log.Println(err.Error())
		rsp.Error.Status = common_proto.Status_FAILURE
		return err
	}
	rsp.UserId = user.Id

	//for token reflesh
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

	rsp.Token, err = p.token.New(req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *UserHandler) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

	log.Println("Debug into VerifyToken: ", req.Token)
	if !p.blacklist.Available(req.Token) {
		err := ankr_default.ErrTokenNeedRefresh
		rsp.Status = common_proto.Status_FAILURE
		rsp.Details = ankr_default.ErrTokenNeedRefresh.Error()
		log.Println(err.Error())
		return err
	}

	if _, err := p.token.Verify(req.Token); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *UserHandler) VerifyAndRefreshToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

	log.Println("Debug into VerifyAndRefreshToken: ", req.Token)

	//for token reflesh
	if !p.blacklist.Available(req.Token) {
		rsp.Status = common_proto.Status_FAILURE
		rsp.Details = ankr_default.ErrTokenNeedRefresh.Error()
		err := ankr_default.ErrTokenNeedRefresh
		log.Println(err.Error())
		return err
	}

	_, err := p.token.Verify(req.Token)
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
