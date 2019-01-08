package handler

import (
	"context"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/db_service"
	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/token"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct{
	db dbservice.DBService
	token token.TokenService
}

func New(dbService dbservice.DBService, tokenService token.TokenService) *UserHandler {
	return &UserHandler{db: dbService, token: tokenService}
}

func (p *UserHandler) Get(ctx context.Context, id *pb.ID, user *pb.User) error {
	var err error
	user, err = p.db.Get(id.Id)
	return err
}

func (p *UserHandler) Add(ctx context.Context, user *pb.User, rsp *pb.Response) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPwd)
	return p.db.Add(*user)
}

func (p *UserHandler) Update(ctx context.Context, user *pb.User, rsp *pb.Response) error {
	return p.db.Update(user)
}

func (p *UserHandler) NewToken(ctx context.Context, user *pb.User, rsp *pb.TokenString) error {
	u, err := p.db.Get(user.Id.Id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return err
	}

	rsp.TokenString, err = p.token.New(user)

	return err
}

func (p *UserHandler) VerifyToken(ctx context.Context, token *pb.TokenString, rsp *pb.Response) error {
	return p.token.Verify(token.TokenString)
}

