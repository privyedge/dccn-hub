package main

import (
	go_micro_srv_dcmgr "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/dcmgr"
	go_micro_srv_taskmgr "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/proto/taskmgr"
	go_micro_srv_usermgr "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	"github.com/Ankr-network/dccn-hub/gateway/handler"
	go_micro_api_dcmgr "github.com/Ankr-network/dccn-hub/gateway/proto/dcmgr"
	go_micro_api_taskmgr "github.com/Ankr-network/dccn-hub/gateway/proto/taskmgr"
	go_micro_api_usermgr "github.com/Ankr-network/dccn-hub/gateway/proto/usermgr"
	"github.com/micro/examples/template/api/client"

	log "github.com/micro/go-log"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.gateway"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init(
		// create wrap for the Example srv client
		micro.WrapHandler(client.ExampleWrapper(service)),
	)

	// Register Handlers
	go_micro_api_usermgr.RegisterUserMgrHandler(service.Server(), handler.NewUserApi(
		go_micro_srv_usermgr.NewUserMgrService("go.micro.srv.usermgr", service.Client()),
	))
	go_micro_api_taskmgr.RegisterTaskMgrHandler(service.Server(), handler.NewTaskApi(
		go_micro_srv_taskmgr.NewTaskMgrService("go.micro.srv.taskmgr", service.Client()),
	))
	go_micro_api_dcmgr.RegisterDcMgrHandler(service.Server(), handler.NewDataCenterApi(
		go_micro_srv_dcmgr.NewDcMgrService("go.micro.srv.dcmgr", service.Client()),
	))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
