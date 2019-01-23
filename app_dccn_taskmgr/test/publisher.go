package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-log/log"
	micro "github.com/micro/go-micro"
	"github.com/pborman/uuid"
)

// send events using the publisher
func sendEv(topic string, p micro.Publisher) {
	t := time.NewTicker(time.Second)

	for _ = range t.C {
		// create new event
		ev := &taskmgr.Event{
			Id:        uuid.NewUUID().String(),
			Timestamp: time.Now().Unix(),
			Message:   fmt.Sprintf("Messaging you all day on %s", topic),
			Op:        taskmgr.OpCode_ADD,
		}

		log.Logf("publishing %+v\n", ev)

		// publish an event
		if err := p.Publish(context.Background(), ev); err != nil {
			log.Logf("error publishing %v", err)
		}
	}
}

func main() {
	// create a service
	service := micro.NewService(
		micro.Name("go.micro.cli.v1"),
		micro.Version("latest"),
	)
	// parse command line
	service.Init()

	// create publisher
	pub := micro.NewPublisher("go.micro.srv.v1", service.Client())

	// pub to topic 1
	go sendEv("go.micro.srv.v1", pub)

	// block forever
	select {}
}
