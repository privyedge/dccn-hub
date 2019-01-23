package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/Ankr-network/dccn-hub/app_dccn_stream/proto/grpc"
)

var (
	addr = "localhost:50051"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}

	defer conn.Close()
	cli := pb.NewStreamerClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	task := pb.Task{
		Name:  "XiaoMing",
		Count: 0,
		Age:   15,
		Flag:  "Client",
	}
	stream, err := cli.ServerStream(ctx)
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				log.Println("EOF")
				close(waitc)
				return
			}

			if err != nil {
				log.Fatal(err.Error())
			}

			task = *in

			log.Printf("Receive New Task: %#v\n", *in)

			time.Sleep(1 * time.Second)
		}
	}()

	for {
		task.Count++
		task.Flag = "CLient"
		if err = stream.Send(&task); err != nil {
			log.Fatal(err.Error())
		}
		time.Sleep(1. * time.Second)
	}

	// if err := stream.CloseSend(); err != nil {
	// 	log.Fatal(err.Error())
	// }
	// <-waitc
}
