package handler

import (
	"github.com/Ankr-network/refactor/app_dccn_account/db_account"
	"github.com/Ankr-network/refactor/app_dccn_account/proto"
	"github.com/Ankr-network/refactor/app_dccn_account/token"
)

type AccountHandler struct{
	accountDB dbaccount.AccountDBService
	token token.TokenService
}

func New(db dbaccount.AccountDBService, token token.TokenService) *AccountHandler {
	return &AccountHandler{
		accountDB:db,
		token:token,
	}
}

func (p *AccountHandler) New(ctx context.Context, ac *accountmgr.Account, rsp *accountmgr.Response) error {
	rsp.Error =  p.accountDB.New(ac)
	return nil
}

func (p *AccountHandler) Get(ctx context.Context, name *accountmgr.Name, rsp *accountmgr.Response) error {
	rsp.Account, rsp.Error = p.accountDB.Get(name.Name)
	return nil
}

func (p *AccountHandler) Auth(ctx context.Context, ac *accountmgr.Account, rsp *accountmgr.Token) error {
	return nil
}

func (p *AccountHandler) ValidateToken(ctx context.Context, t *accountmgr.Token, rsp *accountmgr.Token) error {
	return nil
}

