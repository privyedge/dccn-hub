package main

import (
	"fmt"
	"github.com/Ankr-network/refactor/app_dccn_taskmgr/filter"
	"github.com/micro/go-micro/client"
	"time"

	"github.com/Ankr-network/refactor/app_dccn_taskmgr/handler"
	"github.com/Ankr-network/refactor/proto"
	"github.com/Ankr-network/refactor/app_dccn_taskmgr/subscriber"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

func init() {
	// load the config from a file source
	if err := config.Load(file.NewSource(
		file.WithPath("./config.json"),
	)); err != nil {
		fmt.Println(err)
		return
	}

	// define our own host type
	type Host struct {
		Address string `json:"address"`
		Port    int    `json:"port"`
	}

	var host Host

	// read a database host
	if err := config.Get("hosts", "database").Scan(&host); err != nil {
		fmt.Println(err)
		return
	}

}

func main() {
	// New Service.
	// TODO: add autheracation module
	service := micro.NewService(
		micro.Name("network.ankr.srv.taskmgr"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// Initialise service.
	service.Init()

	// New Publisher filter with policy
	pubNew := micro.NewPublisher("topc.task.new", client.NewClient(client.Wrap(taskfilter.NewACWrapper)))
	// Cancel Publisher broadcast
	pubCancel := micro.NewPublisher("topic.task.cancel",client.DefaultClient )

	// Register Handler.
	taskmgr.RegisterTaskMgrHandler(service.Server(), handler.New(pubNew, pubCancel))

	// Register Function as Subscriber.
	micro.RegisterSubscriber("topic.task.result", service.Server(), subscriber.GetResult)

	// Run service.
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
