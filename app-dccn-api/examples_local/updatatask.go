package main

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"

	//"github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"

	"log"
	"time"

	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	//usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

	//common_proto "github.com/Ankr-network/dccn-common/protos/common"
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

	taskClient := taskmgr.NewTaskMgrClient(conn)
	userClient := usermgr.NewUserMgrClient(conn)

	user := &usermgr.User{
		Name:     "user_test1",
		Nickname: "test1",
		Email:    `1231@Gmail.com`,
		Password: "12345678901",
		Balance:  199,
	}
	//if _, err := userClient.Register(context.Background(), user); err != nil {
	//	log.Fatal(err.Error())
	//} else {
	//	log.Println("Register Ok")
	//}

	var token string
	//var userId string
	if rsp, err := userClient.Login(context.TODO(), &usermgr.LoginRequest{Email: user.Email, Password: user.Password}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Printf("login Success: %s\n", rsp.Token)
		token = rsp.Token
	//	userId = rsp.UserId
	}

	md := metadata.New(map[string]string{
		"token": token,
	})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	task := apiCommon.MockTasks()[0]

	task.Id = "d7a93e80-058e-4637-93a4-0b8b7dd7d211"
	task.Image = "nginx:1.14"
	task.Name = "ok123"

	if rsp, err := taskClient.UpdateTask(tokenContext, &taskmgr.UpdateTaskRequest{UserId: "1234", Task: &task}); err != nil {
		log.Fatal(err.Error())
	} else {

		log.Printf(" result %s  error:  %s", rsp.Status,  rsp.Details);

		}



}
