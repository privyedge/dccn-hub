package main

import (
	"context"
	"log"

	grpc "github.com/micro/go-grpc"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	pb "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	srv := grpc.NewService()

	srv.Init()

	user := &pb.User{
		Name:     "user_test",
		Nickname: "test",
		Email:    `123@Gmail.com`,
		Password: "1234567890",
		Balance:  99,
	}

	cli := pb.NewUserMgrService(ankr_default.UserMgrRegistryServerName, srv.Client())
	if _, err := cli.Register(context.Background(), user); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Register Ok")
	}

	var token string
	if rsp, err := cli.Login(context.TODO(), &pb.LoginRequest{Email: user.Email, Password: user.Password}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("login Success: %s\n", rsp.Token)
		token = rsp.Token
	}

	// Verify Login Token
	if _, err := cli.VerifyToken(context.TODO(), &pb.Token{Token: token}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Verify Login Token Success")
	}

	if rsp, err := cli.NewToken(context.TODO(), user); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("NewToken Success: ", rsp.Token)
		token = rsp.Token
	}

	// Verify NewToken
	if _, err := cli.VerifyToken(context.TODO(), &pb.Token{Token: token}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Verify Login Token Success")
	}

	// Verify different Password
	if _, err := cli.VerifyToken(context.TODO(), &pb.Token{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fyJleHAiOjE1NDgzNzQ5MTksImlzcyI6ImFua3JfbmV0d29yayJ9.crx45JXV6nXiWZtIWLfsMLjA24B2D0_8NYTpujBKilA"}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Token invalid")
	}

	log.Println("END")
}
