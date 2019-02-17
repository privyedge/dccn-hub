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

func (p *ApiUser) Register(ctx context.Context, req *usermgr.User, rsp *common_proto.Empty) error {

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

func (p *ApiUser) Logout(ctx context.Context, req *common_proto.Empty, rsp *common_proto.Empty) error {

	log.Println("Debug into Logout")
	if _, err := p.api.Logout(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// func (p *ApiUser) NewToken(ctx context.Context, req *usermgr.User, rsp *usermgr.NewTokenResponse) error {

// 	log.Println("Debug into NewToken")
// 	if out, err := p.api.NewToken(ctx, req); err != nil {
// 		log.Println(err.Error())
// 		return err
// 	} else {
// 		*rsp = *out
// 	}
// 	return nil
// }

// func (p *ApiUser) VerifyToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

// 	log.Println("Debug into VerifyToken")
// 	if out, err := p.api.VerifyToken(ctx, req); err != nil {
// 		log.Println(err.Error())
// 		return err
// 	} else {
// 		*rsp = *out
// 	}
// 	return nil
// }

// func (p *ApiUser) VerifyAndRefreshToken(ctx context.Context, req *usermgr.Token, rsp *common_proto.Error) error {

// 	log.Println("Debug into VerifyAndRefreshToken")
// 	if out, err := p.api.VerifyAndRefreshToken(ctx, req); err != nil {
// 		log.Println(err.Error())
// 		return err
// 	} else {
// 		*rsp = *out
// 	}
// 	return nil
// }

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

func NewApiUser(c client.Client) *ApiUser {
	return &ApiUser{
		api: usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, c),
	}
}
