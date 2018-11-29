package main

import (
	"fmt"
	pb "github.com/Ankr-network/dccn-hub/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	address = "localhost:50051"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("please input taskid like go run cancel_task.go 93 \n")
		return
	}
	nums, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("please input taskid like go run cancel_task.go 93 \n", os.Args[1])
		return
	}

	// taskID := len(os.Args)
	fmt.Printf("taskID: %d \n", nums)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDccncliClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	r, err := c.CancelTask(ctx, &pb.CancelTaskRequest{Taskid: int64(nums), Usertoken: "ed1605e17374bde6c68864d072c9f5c9"})
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

	fmt.Printf("received Status : %s \n", r.Status)

}
