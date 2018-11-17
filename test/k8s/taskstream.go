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
	datacenter = "datacenter_2"
)

// runRouteChat receives a sequence of route notes, while sending notes for various locations.
func sendTaskStatus(client pb.DccncliClient) {
	 ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
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
			fmt.Printf("Got message type: %s  taskid:  %d  name: %s extra: %s \n", in.Type, in.Taskid , in.Name, in.Extra)

      if in.Type == "NewTask" {
				// start new task

          sendTaskStatusMessage(stream, in.Taskid, "StartSuccess")
					//or   sendTaskStatusMessage(in.Taskid, "StartFailure")
			}

			if in.Type  == "HeartBeat" {
				// do nothing
			}

			if in.Type  == "CancelTask" {
				// start Cancel task
				  sendTaskStatusMessage(stream, in.Taskid, "Cancelled")
			}


      // when job done send back done message  (this is hacking way )
				//	sendTaskStatusMessage(stream, in.Taskid, "Done")




		}
	}()


	go func(stream pb.Dccncli_K8TaskClient) {
		 for {

			 var message = pb.K8SMessage{Type: "HeartBeat",  Datacenter: datacenter, Report:"xxxxreport"}
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

func sendTaskStatusMessage(stream pb.Dccncli_K8TaskClient, taskid int64, status string) {
	var message = pb.K8SMessage{Taskid: taskid, Status:status, Datacenter: datacenter}
	 if err := stream.Send(&message); err != nil {
		 log.Fatalf("Failed to send a note: %v", err)
	 }
	 fmt.Printf("send TaskStatus  message %d %s \n", message.Taskid , message.Status)
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
