package main

import (
	"context"
	"log"
	"strings"

	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	"golang.org/x/crypto/bcrypt"
)

func isEqual(origin, dbUser *pb.User) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(origin.Password)); err != nil {
		log.Println(err.Error())
		return false
	}
	return strings.ToLower(origin.Email) == dbUser.Email &&
		origin.Name == dbUser.Name &&
		origin.Nickname == dbUser.Nickname &&
		origin.Id == dbUser.Id &&
		origin.Balance == dbUser.Balance &&
		origin.IsDeleted == dbUser.IsDeleted
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	srv := micro.NewService()

	srv.Init()

	user := &pb.User{
		Id:       1,
		Name:     "user_test",
		Nickname: "test",
		Email:    `123@Gmail.com`,
		Password: "1234567890",
		Balance:  99,
	}

	cli := pb.NewUserMgrService("go.micro.srv.usermgr", client.DefaultClient)
	if _, err := cli.Create(context.Background(), user); err != nil {
		log.Fatal(err.Error())
	}

	u, err := cli.Get(context.Background(), &pb.Email{Email: user.Email})
	if err != nil {
		log.Fatal(err.Error())
	}

	if !isEqual(user, u) {
		log.Fatalf("want: %#v\n, but %#v\n", user, u)
	}

	if u, err := cli.Login(context.TODO(), &pb.LoginRequest{Email: user.Email, Password: user.Password}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("login feedback: %+v", u.Token)
	}

	token, err := cli.NewToken(context.TODO(), user)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("Receive Token: ", token)
	}

	// Verify same Password
	_, err = cli.VerifyToken(context.TODO(), &pb.Token{Token: token.Token})
	if err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("VerifyToken OK")
	}

	// Verify different Password
	_, err = cli.VerifyToken(context.TODO(), &pb.Token{Token: "fyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NDcyODA4NTIsImlzcyI6ImFwcF9kY2NuX3VzZXJtZ3IifQ.5k3bMjtryTPDZ_v_-_3tgUXEba6eqvN56fa2P7y3wj9"})
	if err != nil {
		log.Println("Invalid OK.")
	} else {
		log.Fatal("VerifyToken Failed.")
	}
}
