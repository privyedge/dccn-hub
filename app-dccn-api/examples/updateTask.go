package main

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"

	//"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"

	"log"
	"time"

	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

	//common_proto "github.com/Ankr-network/dccn-common/protos/common"
//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)

//var addr = "localhost:50051"
var addr = "client-dev.dccn.ankr.network:50051"

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

	taskClient := taskmgr.NewTaskMgrClient(conn)
	userClient := usermgr.NewUserMgrClient(conn)

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





		task :=  common_proto.Task{}
		task.Attributes = &common_proto.TaskAttributes{}
		task.Attributes.Replica = 3

		task.Id = "ea6be29e-f1cd-42ba-948f-a4414af3b076"
		task.Type =  common_proto.TaskType_CRONJOB
		task.TypeData = &common_proto.Task_TypeCronJob{TypeCronJob: &common_proto.TaskTypeCronJob{Image:"nginx:1.13", Schedule:"* * * * *"}}




		if _, err := taskClient.UpdateTask(tokenContext, &taskmgr.UpdateTaskRequest{Task: &task}); err != nil {
			log.Fatal(err.Error())
		} else {

			log.Printf(" UpdateTask success\n")

		}

	}
}
