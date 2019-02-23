package main

import (
	"context"
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
	req.Email = "12112@Gmail.com"
    req.Password = "11111111"



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


	}




	log.Println("END")
}
