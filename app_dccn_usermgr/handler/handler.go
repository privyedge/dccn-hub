package handler

import (
	"context"
	"strings"

	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/db_service"
	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/token"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	db    dbservice.DBService
	token token.IToken
}

func New(dbService dbservice.DBService, tokenService token.IToken) *UserHandler {
	return &UserHandler{db: dbService, token: tokenService}
}

func (p *UserHandler) Get(ctx context.Context, email *pb.Email, user *pb.User) error {
	var err error
	user, err = p.db.Get(strings.ToLower(email.Email))
	return err
}

func (p *UserHandler) Create(ctx context.Context, user *pb.User, rsp *pb.Response) error {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPwd)
	user.Email = strings.ToLower(user.Email)
	return p.db.Create(user)
}

func (p *UserHandler) Login(ctx context.Context, req *pb.LoginRequest, user *pb.User) error {
	user, err := p.db.Get(req.Email)
	if err != nil {
		return err
	}
	user.Token, err = p.token.New(user)
	if err != nil {
		return err
	}
	return nil
}

func (p *UserHandler) NewToken(ctx context.Context, user *pb.User, rsp *pb.Token) error {
	var err error
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(user.Password)); err != nil {
		return err
	}

	rsp.Token, err = p.token.New(user)

	return err
}

func (p *UserHandler) VerifyToken(ctx context.Context, token *pb.Token, rsp *pb.Response) error {
	return p.token.Verify(token.Token)
}
