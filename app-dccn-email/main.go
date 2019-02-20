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

// Init starts handler to listen.
func Init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	Init()

	// New Service
	service := grpc.NewService(
		micro.Name(ankr_default.EmailRegistryServerName),
	)

	// Initialise service
	service.Init()

	// Register Handler
	if err := mail.RegisterMailHandler(service.Server(), new(handler.MailHandler)); err != nil {
		log.Fatal(err.Error())
	}

	// Register Function as TaskStatusFeedback
	opt := service.Server().Options()
	if err := opt.Broker.Connect(); err != nil {
		log.Fatal(err.Error())
	}

	if err := micro.RegisterSubscriber(ankr_default.MQMail, service.Server(), new(subscriber.Subscriber)); err != nil {
		log.Fatal(err.Error())
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
