package main

import (
	"log"
	"net"
  "fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "dccn-hub/protocol"
	"google.golang.org/grpc/reflection"
	util "dccn-hub/util"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}


func (s *server) AddTask(ctx context.Context, in *pb.TaskRequest) (*pb.TaskResponse, error) {
        fmt.Printf("received add task request\n")
				token := in.Usertoken
				user := util.GetUser(token)

        if user.ID == 0 {
					fmt.Printf("add new task fail for user token error\n")
					return &pb.TaskResponse{Status:"Failure", Taskid: -1}, nil
				}else{
					task := util.Task{Name: in.Name, Region: in.Region, Zone: in.Zone, Userid: user.ID}
					id := util.AddTask(task)
					fmt.Printf("add new task ID : %d\n", id)

           return &pb.TaskResponse{Status:"Success", Taskid: id}, nil
				}




}


func (s *server) TaskList(ctx context.Context, in *pb.TaskListRequest) (*pb.TaskListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		fmt.Printf("task list reqeust fail for user token error\n")
		return &pb.TaskListResponse{}, nil
	}else{
		tasks := util.TaskList(int(user.ID))

    var taskList []*pb.TaskInfo
		for i := range tasks {
				task := tasks[i]
				taskInfo := &pb.TaskInfo{}
				taskInfo.Taskid = task.ID
				taskInfo.Taskname = task.Name
				taskInfo.Status = task.Status
				taskList = append(taskList, taskInfo)
				fmt.Printf("task id : %d %s status %s \n", task.ID,task.Name, task.Status)
     }

 		return &pb.TaskListResponse{Tasksinfo: taskList}, nil
	}



}




func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDccncliServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
