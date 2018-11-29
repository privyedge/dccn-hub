package main

import (
	"fmt"
	pb "github.com/Ankr-network/dccn-hub/protocol"
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
	r, err := c.TaskList(ctx, &pb.TaskListRequest{Usertoken: "ed1605e17374bde6c68864d072c9f5c9"})
	if err != nil {
		log.Fatalf("Client: could not send: %v", err)
	}

	fmt.Printf("received task list : \n")
	tasks := r.Tasksinfo
	for i := range tasks {
		task := tasks[i]
		fmt.Printf("task id : %d     name : %s    status : %s \n", int(task.Taskid), task.Taskname, task.Status)
	}
	// todo when have new proto

}
