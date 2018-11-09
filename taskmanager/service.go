package main

import (
	"log"
	"net"
  "fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "dccn-hub/protocol"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}


func (s *server) SendClientStartRequest(ctx context.Context, in *pb.ClientStartRequest) (*pb.ServerStartResponse, error) {
        fmt.Printf("received %s\n", in.Name)
	return &pb.ServerStartResponse{Status:"Success", Taskid:1}, nil
}


func (s *server) SendClientListRequest(ctx context.Context, in *pb.ClientListRequest) (*pb.ServerListResponse, error) {
        fmt.Printf("received\n")
	return &pb.ServerListResponse{}, nil
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
