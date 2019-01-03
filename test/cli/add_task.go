package main

import (
	"fmt"
	"log"
	"time"

	pb "github.com/Ankr-network/dccn-hub/protocol"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
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
	r, err := c.AddTask(ctx, &pb.AddTaskRequest{Name: "nginx:1.12", Region: "us_west", Zone: "ca", Usertoken: "ed1605e17374bde6c68864d072c9f5c9"})
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

	fmt.Printf("received Status : %s \n", r.Status)

}
