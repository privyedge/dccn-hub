package main

import (
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	"github.com/Ankr-network/refactor/app_dccn_dccenter/proto"
	"github.com/Ankr-network/refactor/app_dccn_dccenter/handler"
	"github.com/Ankr-network/refactor/app_dccn_dccenter/subscriber"
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
		micro.Name("network.ankr.srv.dccenter"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	dccenter.RegisterDccenterHandler(service.Server(), new(handler.DcCenterHandler))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("network.ankr.srv.dccenter", service.Server(), new(subscriber.DcCenterSubscriber))

	// Register Function as Subscriber
	micro.RegisterSubscriber("network.ankr.srv.dccenter", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
