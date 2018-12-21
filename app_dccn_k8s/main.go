package main

import (
	"fmt"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	"github.com/Ankr-network/refactor/app_dccn_k8s/proto"
	"github.com/Ankr-network/refactor/app_dccn_k8s/handler"
	"github.com/Ankr-network/refactor/app_dccn_k8s/subscriber"
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
		micro.Name("network.ankr.srv.k8s"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	k8s.RegisterK8SHandler(service.Server(), new(handler.K8sHandler))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("network.ankr.srv.k8s", service.Server(), new(subscriber.K8sSubscriber))

	// Register Function as Subscriber
	micro.RegisterSubscriber("network.ankr.srv.k8s", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
