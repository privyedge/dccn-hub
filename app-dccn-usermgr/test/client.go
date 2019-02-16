package main

import (
	"context"
	"log"

	pb "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	grpc "github.com/micro/go-grpc"
)

func mockUser() *pb.User {
	return &pb.User{
		Password: "1234567",
		Email:    "123@gmail.com",
	}
}

func createUser(cli pb.UserMgrService) error {
	if _, err := cli.Register(context.Background(), mockUser()); err != nil {
		return err
	}
	log.Println("Create User success")
	return nil
}

func newToken(cli pb.UserMgrService, user *pb.User) (string, error) {
	pbToken, err := cli.NewToken(context.Background(), user)
	if err != nil {
		return "", err
	}
	return pbToken.Token, nil
}

func verifyToken(cli pb.UserMgrService, token string) error {
	if _, err := cli.VerifyToken(context.TODO(), &pb.Token{Token: token}); err != nil {
		return err
	}
	log.Println("Verify Token OK.")
	return nil
}

func main() {
	srv := grpc.NewService()

	srv.Init()

	cli := pb.NewUserMgrService("go.micro.srv.v1", srv.Client())
	if cli == nil {
		panic("NIL")
	}

	if err := createUser(cli); err != nil {
		log.Println(err.Error())
		return
	}

	tokenString, err := newToken(cli, mockUser())
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(tokenString)

	if err = verifyToken(cli, tokenString); err != nil {
		log.Println(err.Error())
		return
	}
}
