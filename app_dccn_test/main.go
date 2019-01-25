package main

import (
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	"github.com/Ankr-network/dccn-hub/app_dccn_test/handler"
	"github.com/Ankr-network/dccn-hub/app_dccn_test/subscriber"

	example "github.com/Ankr-network/dccn-hub/app_dccn_test/proto/example"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.app_dccn_test"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.app_dccn_test", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.app_dccn_test", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
