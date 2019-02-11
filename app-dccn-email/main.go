package main

import (
	"log"

	ankr_default "github.com/Ankr-network/dccn-common/protos"

	grpc "github.com/micro/go-grpc"
	micro "github.com/micro/go-micro"

	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
	"github.com/Ankr-network/dccn-hub/app-dccn-email/handler"
	"github.com/Ankr-network/dccn-hub/app-dccn-email/subscriber"

	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name(ankr_default.EmailRegistryServerName),
	)

	// Initialise service
	service.Init()

	// Register Handler
	if err := mail.RegisterMailHandler(service.Server(), new(handler.Mail)); err != nil {
		log.Fatal(err.Error())
	}

	// Register Function as TaskStatusFeedback
	opt := service.Server().Options()
	opt.Broker.Connect()
	if err := micro.RegisterSubscriber(ankr_default.MQMail, service.Server(), subscriber.Handler); err != nil {
		log.Fatal(err.Error())
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
