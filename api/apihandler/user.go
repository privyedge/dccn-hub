package apihandler

import (
	"context"

	"github.com/micro/go-micro/client"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/usermgr/v1"
)

type ApiUser struct {
	api usermgr.UserMgrService
}

func (p *ApiUser) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.User) error {
	out, _ := p.api.Login(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) Create(ctx context.Context, req *usermgr.User, rsp *usermgr.Response) error {
	out, _ := p.api.Login(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) Get(ctx context.Context, req *usermgr.Email, rsp *usermgr.User) error {
	out, _ := p.api.Login(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.Token) error {
	out, _ := p.api.Login(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *usermgr.Response) error {
	out, _ := p.api.Login(ctx, req)
	*rsp = *out
	return nil
}

func NewApiUser(c client.Client) *ApiUser {
	return &ApiUser{
		api: usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, c),
	}
}
