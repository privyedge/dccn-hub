package main

import (
	"context"
	"io"
	"log"

	grpc "github.com/micro/go-grpc"

	pb "github.com/Ankr-network/dccn-hub/app_dccn_stream/proto/micro"
	stream "github.com/Ankr-network/dccn-hub/app_dccn_stream/proto/micro"
)

type StreamHandler struct {
}

func (p *StreamHandler) ServerStream(ctx context.Context, stream stream.Streamer_ServerStreamStream) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println("EOF")
			return nil
		}
		if err != nil {
			log.Println(err.Error())
			return err
		}

		log.Printf("Recevie Something here: %#v\n", *in)

		in.Flag = "Server"
		if err := stream.Send(in); err != nil {
			log.Println(err.Error())
		}
	}
}

func main() {
	srv := grpc.NewService()

	srv.Init()

	if err := pb.RegisterStreamerHandler(srv.Server(), &StreamHandler{}); err != nil {
		log.Fatal(err.Error())
	}

	if err := srv.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
