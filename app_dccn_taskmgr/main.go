package main

import (
	"fmt"
	"time"

	"github.com/Ankr-network/refactor/app_dccn_taskmgr/handler"
	"github.com/Ankr-network/refactor/app_dccn_taskmgr/proto"
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
	// New Service
	service := micro.NewService(
		micro.Name("network.ankr.srv.taskmgr"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*10),
	)

	// Initialise service
	service.Init()

	// Register Handler
	taskmgr.RegisterTaskMgrHandler(service.Server(), new(handler.TaskMgrHandler))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("network.ankr.srv.taskmgr", service.Server(), new(subscriber.TaskMgrSubscriber))

	// Register Function as Subscriber
	micro.RegisterSubscriber("network.ankr.srv.taskmgr", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
