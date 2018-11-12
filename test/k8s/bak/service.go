package main

import (
	"log"
	"net"
  "fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "dccn-hub/protocol/dccnk8"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50052"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}


func (s *server) K8AddTask(ctx context.Context, in *pb.K8TaskRequest) (*pb.K8TaskResponse, error) {
    fmt.Printf("receive add task request\n")

     return &pb.K8TaskResponse{Status:"Success", Taskid: 1234}, nil

}


func (s *server) K8KeepAlive(ctx context.Context, in *pb.K8KeepaliveRequest) (*pb.K8KeepaliveResponse, error) {

   return &pb.K8KeepaliveResponse{}, nil

}




func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDccnk8Server(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
