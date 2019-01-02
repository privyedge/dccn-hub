package handler

import (
	"context"

	"github.com/Ankr-network/refactor/proto/usermgr"
	"github.com/Ankr-network/refactor/app_dccn_usermgr/db_user"
	"github.com/Ankr-network/refactor/app_dccn_usermgr/token"
)

type UserHandler struct{
	dbuser.DBService
	token.TokenService
}

func New(dbService dbuser.DBService, tokenService token.TokenService) *UserHandler {
	return &UserHandler{DBService: dbService, TokenService: tokenService}
}

func (p *UserHandler) New(ctx context.Context, user *usermgr.User, rsp *usermgr.Response) error {
	err := p.DBService.New(user)
	if err != nil {
		rsp.Error.Code = 101
		rsp.Error.Description = err.Error()
	}
	return nil
}

func (p *UserHandler) Get(ctx context.Context, name *usermgr.Name, rsp *usermgr.Response) error {
	user, err := p.DBService.Get(name.Name)
	if err != nil {
		rsp.Error.Code = 102
		rsp.Error.Description = err.Error()
	} else {
		rsp.User = user
	}
	return nil
}

func (p *UserHandler) NewToken(ctx context.Context, user *usermgr.User, rsp *usermgr.Token) error {
	tokenString, err := p.TokenService.New(user)
	if err != nil {
		rsp.Error.Code = 103
		rsp.Error.Description = "New Token Error"
	} else {
		rsp.TokenString = tokenString
	}

	return nil
}

func (p *UserHandler) VerifyToken(ctx context.Context, token *usermgr.Token, rsp *usermgr.Token) error {
	if err := p.TokenService.Verify(token.TokenString); err != nil {
		rsp.Error.Code = 104
		rsp.Error.Description = err.Error()
	} else {
		rsp.Valid = true
	}

	return nil
}

