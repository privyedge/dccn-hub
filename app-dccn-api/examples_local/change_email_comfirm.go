package main

import (
	"context"
	"google.golang.org/grpc/metadata"
	"time"

	//	"google.golang.org/grpc/metadata"
//	"time"

	//"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"

	//	"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
	"log"
//	"time"

//	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"google.golang.org/grpc"
//	"google.golang.org/grpc/metadata"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

//	common_proto "github.com/Ankr-network/dccn-common/protos/common"
//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)

var addr = "localhost:50051"
//var addr = "client-dev.dccn.ankr.network:50051"

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

//	taskClient := taskmgr.NewTaskMgrClient(conn)
	userClient := usermgr.NewUserMgrClient(conn)

	req := &usermgr.LoginRequest{}
	req.Email = "yousong.zhang2@gmail.com"
    req.Password = "1111112c11"



	//var token string
	//var userId string
	if rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: req.Email, Password: req.Password}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("response %+v \n", rsp)
		//log.Printf("login Success: id : %s name : %s , email %s  refresh_token : %s  access_token %s \n", rsp.User.Id, rsp.User.Attributes.Name, rsp.User.Email ,rsp.AuthenticationResult.RefreshToken, rsp.AuthenticationResult.AccessToken)
		//token = rsp.Token
		//userId = rsp.UserId
		refresh_token := rsp.AuthenticationResult.RefreshToken
		access_token := rsp.AuthenticationResult.AccessToken



		log.Printf("get access_token after login %s  refresh_token %s \n", access_token, refresh_token)

		md := metadata.New(map[string]string{
			"token": access_token,
		})

		log.Printf("get access_token after login %s \n", access_token)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()




		request := &usermgr.ConfirmEmailRequest{}
		request.NewEmail =  "yousong@ankr.com"
		request.ConfirmationCode = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTA5ODM1NDIsImp0aSI6IjQyZjQ5Y2FlLTAzYzUtNDgzOS04OWI3LTllMzdjMmZjNTk1ZSIsImlzcyI6ImFua3IubmV0d29yayJ9.IckQvHs-lse3ZKeEXP-Rf6lu3APoPUPKK8GITM9LLRo"


		if rsp, err := userClient.ConfirmEmail(tokenContext, request); err != nil {
			//	log.Fatal(err.Error())
			log.Printf("ConfirmEmail have some error : %s \n", err.Error())
		} else {
			log.Printf("ConfirmEmail result no error  %+v", rsp)
		}


	}




	log.Println("END")
}
