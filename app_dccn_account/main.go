package main

import (
	"fmt"
	"github.com/Ankr-network/refactor/app_dccn_account/proto"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"time"
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
	service := k8s.NewService(
		micro.Name("network.ankr.srv.account"),
		micro.Version("latest"),
		micro.RegisterTTL(time.Second*30),
		micro.RegisterInterval(time.Second*15),
	)

	// Initialise service
	service.Init()

	// Register Handler
	accountmgr.RegisterAccountMgrHandler(service.Server(), )

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
