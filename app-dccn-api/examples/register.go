package main

import (
	"context"
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

	//	taskClient := taskmgr.NewTaskMgrClient(conn)
	userClient := usermgr.NewUserMgrClient(conn)


	req := usermgr.RegisterRequest{}
	req.User = &usermgr.User{}
	req.User.Email = `1223@Gmail.com`
	req.User.Attributes = &usermgr.UserAttributes{}
	req.User.Attributes.Name = "ankrtest"
	req.Password = "111111c"

	if _, err := userClient.Register(context.Background(), &req); err != nil {
		//	log.Fatal(err.Error())
		log.Printf("receive have some error : %s \n", err.Error())
	} else {
		log.Printf("Register result no error  ")
	}

}
