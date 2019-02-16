package apihandler

import (
	"context"
	"log"

	"github.com/micro/go-micro/client"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
)

type ApiUser struct {
	api usermgr.UserMgrService
}

func (p *ApiUser) Register(ctx context.Context, req *usermgr.User, rsp *common_proto.Error) error {

	log.Println("Debug into Register")
	if out, err := p.api.Register(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) ForgetPassword(ctx context.Context, req *usermgr.AskResetPasswordRequest, rsp *common_proto.Error) error {

	log.Println("Debug AskResetPassword")
	if out, err := p.api.AskResetPassword(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) ConfirmPassword(ctx context.Context, req *usermgr.ResetPasswordRequest, rsp *common_proto.Error) error {

	log.Println("Debug ResetPassword")
	if out, err := p.api.ResetPassword(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) ConfirmRegistration(ctx context.Context, req *usermgr.ActivateRequest, rsp *common_proto.Error) error {

	log.Println("Debug Activate")
	if out, err := p.api.Activate(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {

	log.Println("Debug into Login")
	if out, err := p.api.Login(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) Logout(ctx context.Context, req *usermgr.LogoutRequest, rsp *common_proto.Error) error {

	log.Println("Debug into Logout")
	if out, err := p.api.Logout(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.NewTokenResponse) error {

	log.Println("Debug into NewToken")
	if out, err := p.api.NewToken(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

	log.Println("Debug into VerifyToken")
	if out, err := p.api.VerifyToken(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) VerifyAndRefreshToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

	log.Println("Debug into VerifyAndRefreshToken")
	if out, err := p.api.VerifyAndRefreshToken(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func (p *ApiUser) RefreshToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {
	if out, err := p.api.RefreshToken(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}
	return nil
}

func NewApiUser(c client.Client) *ApiUser {
	return &ApiUser{
		api: usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, c),
	}
}
