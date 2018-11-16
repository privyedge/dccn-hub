package main

import (
	"log"
	"io"
	"time"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "dccn-hub/protocol"
)


const (
	address  = "localhost:50051"
)

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func sendTaskStatus(client pb.DccncliClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := client.K8Task(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {

				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			fmt.Printf("Got message %d %s %s \n", in.Taskid , in.Name, in.Extra)
		}
	}()

	  var message = pb.TaskStatus{Taskid: -1, Status:"Failure"}
		if err := stream.Send(&message); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}

			fmt.Printf("send TaskStatus  message %d %s \n", message.Taskid , message.Status)

	//stream.CloseSend()
	<-waitc
}

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDccncliClient(conn)


	//_, cancel := context.WithTimeout(context.Background(), 30 * time.Second )
	//defer cancel()

  sendTaskStatus(c)



}
