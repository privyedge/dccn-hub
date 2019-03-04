package main

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"

	"log"
	"time"

	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

	//	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	//	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)
//var addr = "localhost:50051"
var addr = "client-dev.dccn.ankr.network:50051"

//var addr = "afcac29ea274711e99cb106bbae7419f-1982485008.us-west-1.elb.amazonaws.com:50051"

//func parseError(err string) string{
//
//}

func main() {

	log.SetFlags(log.LstdFlags | log.Llongfile)
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

		task := common_proto.Task{}
		//task.DataCenterName = "datacenter_tokyo"
		task.Name = "task"
		task.Type = common_proto.TaskType_DEPLOYMENT
		task.Attributes = &common_proto.TaskAttributes{}
		task.Attributes.Replica = 1
		t := common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image:"nginx:1.12"}}
		task.TypeData = &t

		if rsp, err := taskClient.CreateTask(tokenContext, &taskmgr.CreateTaskRequest{Task: &task}); err != nil {

			//log.Println("detail create %+v " + rsp)
			log.Fatal(err)
		} else {
			log.Println("create task successfully : taskid   " + rsp.TaskId)
		}

	}

	//// var userTasks []*common_proto.Task
	//userTasks := make([]*common_proto.Task, 0)
	//if rsp, err := taskClient.TaskList(tokenContext, &taskmgr.ID{UserId: userId}); err != nil {
	//	log.Fatal(err.Error())
	//} else {
	//	userTasks = append(userTasks, rsp.Tasks...)
	//	if len(userTasks) == 0 {
	//		log.Fatalf("no tasks belongs to %s", userId)
	//	} else {
	//		log.Println(len(userTasks), "tasks belongs to ", user.Email)
	//		for i := 0; i < len(userTasks); i++ {
	//			log.Println(userTasks[i])
	//		}
	//
	//	}
	//}
	//
	//log.Println("END")
}
