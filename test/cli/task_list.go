package main

import (
	"log"
	"time"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "dccn-hub/protocol"
)


const (
	address  = "localhost:50051"
)



func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSimpappClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second )
	defer cancel()
	r, err := c.SendClientListRequest(ctx, &pb.ClientListRequest{Usertoken:"ed1605e17374bde6c68864d072c9f5c9" })
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

 // todo when have new proto
	fmt.Printf("received Status : %s \n", r.Taskname)



}
