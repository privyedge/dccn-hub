package main

import (
	"context"
	"log"

	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	micro "github.com/micro/go-micro"
)

func CreateUser(cli pb.UserMgrService) {
	if _, err := cli.Create(context.Background(), &pb.User{
		Nickname: "xiaoming",
		Password: "1234567",
		Email:    "123@gmail.com",
	}); err != nil {
		panic(err.Error())
	}
	log.Println("Create User success")
}

func GetUser(cli pb.UserMgrService, id *pb.ID) *pb.User {
	user, err := cli.Get(context.TODO(), id)
	if err != nil {
		panic(err.Error())
	}
	return user
}

func GetUserByEmail(cli pb.UserMgrService) *pb.User {
	user, err := cli.GetByEmail(context.Background(), &pb.Email{Email: "123@gmail.com"})
	if err != nil {
		panic(err.Error())
	}
	return user
}

func NewToken(cli pb.UserMgrService, user *pb.User) string {
	pbToken, err := cli.NewToken(context.Background(), user)
	if err != nil {
		panic(err.Error())
	}
	return pbToken.TokenString
}

func VerifyToken(cli pb.UserMgrService, tokenString string) {
	if _, err := cli.VerifyToken(context.Background(), &pb.TokenString{TokenString: tokenString}); err != nil {
		panic(err.Error())
	}
	log.Println("Verify Token OK.")
}

func main() {
	srv := micro.NewService()

	srv.Init()

	cli := pb.NewUserMgrService("go.micro.srv.usermgr", srv.Client())

	// CreateUser(cli)
	user := GetUserByEmail(cli)
	user2 := GetUser(cli, user.Id)
	if user != user2 {
		panic("NOT EQUAL")
	}

	tokenString := NewToken(cli, user)
	log.Println(tokenString)

	VerifyToken(cli, tokenString)
}
