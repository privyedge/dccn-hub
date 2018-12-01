package main

import (
	"fmt"
	pb "dccn-hub/protocol"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

const (
	address    = "127.0.0.1:50051"
	datacenter = "datacenter_3"
)

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func sendTaskStatus(client pb.DccncliClient) {
	ctx, cancel := context.WithCancel(context.Background())
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

				print(ctx.Err())
				log.Fatalf("Failed to receive a note : %v", err)
			}
			fmt.Printf("Got message type: %s  taskid:  %d  name: %s image %s extra: %s \n", in.Type, in.Taskid, in.Name, in.Image,  in.Extra)

			if in.Type == "NewTask" {
				// start new task

				sendTaskStatusMessage(stream, in.Taskid, in.Name, "StartSuccess")
				//or   sendTaskStatusMessage(in.Taskid, "StartFailure")
			}

			if in.Type == "HeartBeat" {
				// do nothing
			}

			if in.Type == "CancelTask" {
				// start Cancel task
				sendTaskStatusMessage(stream, in.Taskid, in.Name, "Cancelled")
			}

			// when job done send back done message  (this is hacking way )
			//	sendTaskStatusMessage(stream, in.Taskid, "Done")

		}
	}()

	go func(stream pb.Dccncli_K8TaskClient) {
		for {

			var message = pb.K8SMessage{Type: "HeartBeat", Datacenter: datacenter, Report: "xxxxreport"}
			if err := stream.Send(&message); err != nil {
				log.Fatalf("Failed to send HeartBeat: %v", err)
			}

			fmt.Printf("send HeartBeat  \n")

			time.Sleep(time.Second * 30)
		}
	}(stream)

	//stream.CloseSend()
	<-waitc
}

func sendTaskStatusMessage(stream pb.Dccncli_K8TaskClient, taskid int64, taskName string, status string) {
	var message = pb.K8SMessage{Taskid: taskid, Taskname: taskName, Status: status, Datacenter: datacenter}
	if err := stream.Send(&message); err != nil {
		log.Fatalf("Failed to send a note: %v", err)
	}
	fmt.Printf("send TaskStatus message task:  %s  status:  %s \n", message.Taskname, message.Status)
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
