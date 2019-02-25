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

var addr = "client-stage.dccn.ankr.network:50051"

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

	userClient := usermgr.NewUserMgrClient(conn)

	//&email=
	req := usermgr.ConfirmRegistrationRequest{
		Email:            "starky@udm.ru",
		ConfirmationCode: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTExMTczMjIsImp0aSI6IjE0NDkxYTIzLWIyZWEtNGNkNi04YzY5LTE5MzViNzIxNWM5OSIsImlzcyI6ImFua3IubmV0d29yayJ9.lVFpGIi6gAiDaq6JRiJMq3xppLQ-2GaqBpLNSQN43wA",
	}

	if _, err := userClient.ConfirmRegistration(context.Background(), &req); err != nil {
		//	log.Fatal(err.Error())
		log.Fatal("receive have some error : %s \n", err.Error())
	} else {
		log.Printf("Register result no error  ")
	}

}
