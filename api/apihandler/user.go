package apihandler

import (
	"context"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/usermgr/v1"
	"github.com/micro/go-micro/client"
)

type ApiUser struct {
	api usermgr.UserMgrService
}

func NewApiUser(c client.Client) *ApiUser {
	return &ApiUser{
		api: usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, c),
	}
}

func (p *ApiUser) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.User) error {
	return p.Login(ctx, req, rsp)
}

func (p *ApiUser) Create(ctx context.Context, req *usermgr.User, rsp *usermgr.Response) error {
	return p.Create(ctx, req, rsp)
}

func (p *ApiUser) Get(ctx context.Context, req *usermgr.Email, rsp *usermgr.User) error {
	return p.Get(ctx, req, rsp)
}

func (p *ApiUser) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.Token) error {
	return p.NewToken(ctx, req, rsp)
}

func (p *ApiUser) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *usermgr.Response) error {
	return p.VerifyToken(ctx, req, rsp)
}
