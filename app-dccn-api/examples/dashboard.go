package main

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"

	//"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"

	//"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"

	//"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
	"log"
	"time"

	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

	//common_proto "github.com/Ankr-network/dccn-common/protos/common"
	//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)

//var addr = "localhost:50051"
var addr = "client.dccn.ankr.network:50051"

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

	dcClient := dcmgr.NewDCAPIClient(conn)
	userClient := usermgr.NewUserMgrClient(conn)
	taskClient := taskmgr.NewTaskMgrClient(conn)

	req := &usermgr.LoginRequest{}
	req.Email = `yousong.zhang@gmail.com`
	req.Password = "zddzys123"

	//var userId string
	if rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: req.Email, Password: req.Password}); err != nil {
		if err == ankr_default.ErrPasswordError {
			log.Printf("password error  %s", err.Error())
		}

		a := err.Error()

		log.Printf(">>>>>%d  %s <<<<", len(a), a)
	} else {
		log.Printf("response %+v \n", rsp)
		//log.Printf("login Success: id : %s name : %s , email %s  refresh_token : %s  access_token %s \n", rsp.User.Id, rsp.User.Attributes.Name, rsp.User.Email ,rsp.AuthenticationResult.RefreshToken, rsp.AuthenticationResult.AccessToken)
		//token = rsp.Token
		//userId = rsp.UserId
		refresh_token := rsp.AuthenticationResult.RefreshToken
		access_token := rsp.AuthenticationResult.AccessToken

		log.Printf("refresh_token  %s  access_token %s", refresh_token, access_token)

		md := metadata.New(map[string]string{
			"token": access_token,
		})

		log.Printf("get access_token after login %s \n", access_token)
		ctx := metadata.NewOutgoingContext(context.Background(), md)

		tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		// var userTasks []*common_proto.Task

		log.Printf("\n\n")
		if rsp, err := dcClient.NetworkInfo(tokenContext, &common_proto.Empty{}); err != nil {
			log.Fatal(err.Error())
		} else {
			log.Printf("network info host count %+v   traffic %d", rsp, rsp.Traffic)

		}

		//log.Printf("DataCenterLeaderBoard info >>>")
		if rsp, err := dcClient.DataCenterLeaderBoard(tokenContext, &common_proto.Empty{}); err != nil {
			log.Fatal(err.Error())
		} else {
			list := rsp.List
			log.Printf("DataCenterLeaderBoard list : %+v ", list)

		}


		//log.Printf("TaskOverview info >>>")
		if rsp, err := taskClient.TaskOverview(tokenContext, &common_proto.Empty{}); err != nil {
			log.Fatal(err.Error())
		} else {
		  // count := rsp.EnvironmentCount
			log.Printf("TaskOverview list : %d", rsp)

		}



		if rsp, err := taskClient.TaskLeaderBoard(tokenContext, &common_proto.Empty{}); err != nil {
			log.Fatal(err.Error())
		} else {
			// count := rsp.EnvironmentCount
			list := rsp.List
			log.Printf("TaskLeaderBoard list : %+v", list)

		}

		//log.Printf("tasklist info >>>")
		//if rsp, err := taskClient.TaskList(tokenContext, &taskmgr.TaskListRequest{}); err != nil {
		//	log.Fatal(err.Error())
		//} else {
		//	// count := rsp.EnvironmentCount
		//	list := rsp.Tasks
		//	log.Printf("tasklist list : %+v", list)
		//
		//}
		//
		//log.Println("END")
	}
}
