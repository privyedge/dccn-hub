package handler

import (
	"context"
	"log"
	"strings"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	dccnwrapper "github.com/Ankr-network/dccn-common/wrapper"
	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/token"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db    dbservice.DBService
	token token.IToken
}

func New(dbService dbservice.DBService, tokenService token.IToken) *UserHandler {
	return &UserHandler{db: dbService, token: tokenService}
}

func (p *UserHandler) Register(ctx context.Context, user *usermgr.User, rsp *common_proto.Error) error {

	log.Println("Debug Register")
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		dccnwrapper.PbError(&rsp, err)
		log.Println(rsp.Details)
		return nil
	}
	user.Password = string(hashedPwd)
	user.Email = strings.ToLower(user.Email)
	user.Id = uuid.New().String()
	if err := p.db.Create(user); err != nil {
		dccnwrapper.PbError(&rsp, err)
		log.Println(rsp.Details)
		return nil
	}
	return nil
}

func (p *UserHandler) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {

	log.Println("Debug Login")
	user, err := p.db.Get(strings.ToLower(req.Email))
	if err != nil {
		dccnwrapper.PbError(&rsp.Error, err)
		log.Println(rsp.Error.Details)
		return nil
	}
	rsp.Token, err = p.token.New(user)
	if err != nil {
		dccnwrapper.PbError(&rsp.Error, err)
		log.Println(rsp.Error.Details)
		return nil
	}
	return nil
}

func (p *UserHandler) Logout(ctx context.Context, in *usermgr.LogoutRequest, out *common_proto.Error) error {
	log.Println("Debug into Logout")
	return nil
}

func (p *UserHandler) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.NewTokenResponse) error {

	log.Println("Debug into NewToken")
	req, err := p.db.Get(strings.ToLower(req.Email))
	if err != nil {
		dccnwrapper.PbError(&rsp.Error, err)
		log.Println(rsp.Error.Details)
		return nil
	}

	rsp.Token, err = p.token.New(req)
	if err != nil {
		dccnwrapper.PbError(&rsp.Error, err)
		log.Println(rsp.Error.Details)
		return nil
	}

	return nil
}

func (p *UserHandler) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

	log.Println("Debug into VerifyToken")
	if err := p.token.Verify(req.Token); err != nil {
		dccnwrapper.PbError(&rsp, err)
		log.Println(rsp.Details)
		return nil
	}
	return nil
}
