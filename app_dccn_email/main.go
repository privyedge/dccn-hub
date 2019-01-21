package main

import (
	"log"

	micro "github.com/micro/go-micro"

	mail "github.com/Ankr-network/dccn-common/protos/email/v1"
	"github.com/Ankr-network/dccn-hub/app_dccn_email/handler"
	"github.com/Ankr-network/dccn-hub/app_dccn_email/subscriber"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.srv.v1.mail"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	mail.RegisterMailHandler(service.Server(), new(handler.Mail))

	// Register Function as TaskStatusFeedback
	micro.RegisterSubscriber("go.micro.srv.v1.mail", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
