package apihandler

import (
	"context"
	"log"

	"github.com/micro/go-micro/client"

	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
)

type ApiUser struct {
	api usermgr.UserMgrService
}

func (p *ApiUser) Register(ctx context.Context, req *usermgr.RegisterRequest, rsp *common_proto.Empty) error {

	log.Println("Debug into Register")
	if _, err := p.api.Register(ctx, req); err != nil {
		log.Println(err.Error())
		return err
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

func (p *ApiUser) Logout(ctx context.Context,  req *usermgr.RefreshToken, rsp *common_proto.Empty) error {

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



func (p *ApiUser)  ConfirmRegistration(ctx context.Context, req *usermgr.ConfirmRegistrationRequst, rsp *common_proto.Empty) error{
	//todo
	return nil
}

func (p *ApiUser)  ForgotPassword(ctx context.Context, req *usermgr.ForgotPasswordRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *ApiUser)  ConfirmPassword(ctx context.Context, req *usermgr.ConfirmPasswordRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *ApiUser) ChangePasword(ctx context.Context, req *usermgr.ChangePasswordRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *ApiUser) UpdateAttributes(ctx context.Context, req *usermgr.UpdateAttributesRequest, rsp *usermgr.User) error {
	//todo
	return nil
}


func (p *ApiUser) ChangeEmail(ctx context.Context, req *usermgr.ChangeEmailRequst, rsp *common_proto.Empty) error {
	//todo
	return nil
}

func (p *ApiUser) VerifyEmail(ctx context.Context, req *usermgr.VerifyEmailRequst, rsp *usermgr.User) error {

		out, err := p.api.VerifyEmail(ctx, req)
		if err != nil {
		log.Println(err.Error())
		return err
	}
		*rsp = *out
		return nil

}


func (p *ApiUser) VerifyAccessToken(ctx context.Context, req *common_proto.Empty, rsp *common_proto.Empty) error {
	return nil

}







func NewApiUser(c client.Client) *ApiUser {
	return &ApiUser{
		api: usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, c),
	}


}
