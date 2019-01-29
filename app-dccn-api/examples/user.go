package main

import (
	"context"
	"log"

	pb "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"google.golang.org/grpc"
)

var addr = "localhost:50051"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)

	user := &pb.User{
		Name:     "user_test",
		Nickname: "test",
		Email:    `123@Gmail.com`,
		Password: "1234567890",
		Balance:  99,
	}

	cli := pb.NewUserMgrClient(conn)
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
		log.Println("OK, Token is invalid")
	}

	log.Println("END")
}
