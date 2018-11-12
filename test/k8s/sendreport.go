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
	c := pb.NewDccncliClient(conn)


	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second )
	defer cancel()
	r, err := c.K8ReportStatus(ctx, &pb.ReportRequest{Name:"datacenter_2",Report:"job 1 job2 job3 host 100", Host:"127.0.0.1", Port:5005 })
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

	fmt.Printf("received Status : %s \n", r.Status)



}
