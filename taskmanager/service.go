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


func (s *server) SendClientStartRequest(ctx context.Context, in *pb.ClientStartRequest) (*pb.ServerStartResponse, error) {
        fmt.Printf("received add task request\n")
				token := in.Usertoken
				user := util.GetUser(token)

        if user.ID == 0 {
					fmt.Printf("add new task fail for user token error\n")
					return &pb.ServerStartResponse{Status:"Failure", Taskid: -1}, nil
				}else{
					task := util.Task{Name: in.Name, Region: in.Region, Zone: in.Zone, Userid: user.ID}
					id := util.AddTask(task)
					fmt.Printf("add new task ID : %d\n", id)
           return &pb.ServerStartResponse{Status:"Success", Taskid: id}, nil
				}




}


func (s *server) SendClientListRequest(ctx context.Context, in *pb.ClientListRequest) (*pb.ServerListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		fmt.Printf("task list reqeust fail for user token error\n")
		return &pb.ServerListResponse{}, nil
	}else{
		tasks := util.TaskList(int(user.ID))

		for i := range tasks {
				task := tasks[i]
				fmt.Printf("task id : %d \n", task.ID)
     }

	}
	// todo when have new list proto
	return &pb.ServerListResponse{Taskname: "this is a list"}, nil
}




func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterSimpappServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
