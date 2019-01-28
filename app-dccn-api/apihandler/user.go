package apihandler

import (
	"context"

	"github.com/micro/go-micro/client"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
)

type ApiUser struct {
	api usermgr.UserMgrService
}

func (p *ApiUser) Register(ctx context.Context, req *usermgr.User, rsp *common_proto.Error) error {
	out, _ := p.api.Register(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {
	out, _ := p.api.Login(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) Logout(ctx context.Context, req *usermgr.LogoutRequest, rsp *common_proto.Error) error {
	out, _ := p.api.Logout(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.NewTokenResponse) error {
	out, _ := p.api.NewToken(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiUser) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {
	out, _ := p.api.VerifyToken(ctx, req)
	*rsp = *out
	return nil
}

func NewApiUser(c client.Client) *ApiUser {
	return &ApiUser{
		api: usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, c),
	}
}
