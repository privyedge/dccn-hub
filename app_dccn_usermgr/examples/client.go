package main

import (
	"context"
	"log"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	pb "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	dccnwrapper "github.com/Ankr-network/dccn-common/wrapper"

	grpc "github.com/micro/go-grpc"
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
	if rsp, _ := cli.Register(context.Background(), user); dccnwrapper.IsSuccess("Register", rsp) {
		log.Println("Register Ok")
	}

	var token string
	if rsp, _ := cli.Login(context.TODO(), &pb.LoginRequest{Email: user.Email, Password: user.Password}); dccnwrapper.IsSuccess("Login", rsp.Error) {
		log.Printf("login Success: %s\n", rsp.Token)
		token = rsp.Token
	}

	// Verify Login Token
	if rsp, _ := cli.VerifyToken(context.TODO(), &pb.Token{Token: token}); dccnwrapper.IsSuccess("VerifyToken", rsp) {
		log.Println("Verify Login Token Success")
	}

	if rsp, _ := cli.NewToken(context.TODO(), user); dccnwrapper.IsSuccess("NewToken", rsp.Error) {
		log.Println("NewToken Success: ", rsp.Token)
		token = rsp.Token
	}

	// Verify NewToken
	if rsp, _ := cli.VerifyToken(context.TODO(), &pb.Token{Token: token}); dccnwrapper.IsSuccess("VerifyToken", rsp) {
		log.Println("Verify Login Token Success")
	}

	// Verify different Password
	if rsp, _ := cli.VerifyToken(context.TODO(), &pb.Token{Token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.fyJleHAiOjE1NDgzNzQ5MTksImlzcyI6ImFua3JfbmV0d29yayJ9.crx45JXV6nXiWZtIWLfsMLjA24B2D0_8NYTpujBKilA"}); rsp.Status == common_proto.Status_FAILURE {
		log.Println("Token invalid")
	}

	log.Println("END")
}
