package main

import (
	"fmt"
		pb "github.com/Ankr-network/dccn-rpc/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDccncliClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := c.K8QueryTask(ctx, &pb.QueryTaskRequest{Name: "datacenter_2"})
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

	fmt.Printf("received new task  : %d %s %s \n", r.Taskid, r.Name, r.Extra)

}
