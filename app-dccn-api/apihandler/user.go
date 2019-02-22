package apihandler

import (
	"context"
	"log"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
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

func (p *ApiUser) Register(ctx context.Context, req *usermgr.RegisterRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into Register")
	if out, err := p.api.Register(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}

	return nil
}

func (p *ApiUser) Login(ctx context.Context, req *usermgr.LoginRequest, rsp *usermgr.LoginResponse) error {

	log.Println("Debug into Login")
	out, err := p.api.Login(ctx, req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*rsp = *out
	return nil
}

func (p *ApiUser) Logout(ctx context.Context, req *usermgr.RefreshToken, rsp *common_proto.Empty) error {

	log.Println("Debug into Logout")
	if _, err := p.api.Logout(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *ApiUser) RefreshSession(
	ctx context.Context, req *usermgr.RefreshToken, rsp *usermgr.AuthenticationResult) error {
	out, err := p.api.RefreshSession(ctx, req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*rsp = *out
	return nil
}

func (p *ApiUser) ConfirmRegistration(ctx context.Context, req *usermgr.ConfirmRegistrationRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ConfirmRegistration")
	if _, err := p.api.ConfirmRegistration(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *ApiUser) ForgotPassword(ctx context.Context, req *usermgr.ForgotPasswordRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ForgotPassword")
	if _, err := p.api.ForgotPassword(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *ApiUser) ChangePassword(ctx context.Context, req *usermgr.ChangePasswordRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ChangePassword")
	if _, err := p.api.ChangePassword(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *ApiUser) ConfirmPassword(ctx context.Context, req *usermgr.ConfirmPasswordRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ConfirmPassword")
	if _, err := p.api.ConfirmPassword(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *ApiUser) UpdateAttributes(ctx context.Context, req *usermgr.UpdateAttributesRequest, rsp *usermgr.User) error {

	log.Println("Debug into UpdateAttributes")
	if out, err := p.api.UpdateAttributes(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}

	return nil
}

func (p *ApiUser) ChangeEmail(ctx context.Context, req *usermgr.ChangeEmailRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ChangeEmail")
	if _, err := p.api.ChangeEmail(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *ApiUser) VerifyAccessToken(ctx context.Context, req *common_proto.Empty, rsp *common_proto.Empty) error {

	log.Println("Debug into VerifyAccessToken")
	if _, err := p.api.VerifyAccessToken(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil

}

func (p *ApiUser) ConfirmEmail(ctx context.Context, req *usermgr.ConfirmEmailRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into ConfirmEmail")
	if out, err := p.api.ConfirmEmail(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	} else {
		*rsp = *out
	}

	return nil
}
