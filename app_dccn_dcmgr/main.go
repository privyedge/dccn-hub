package main

import (
	"time"

	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/db_service"
	"github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/handler"
	go_micro_srv_dcmgr "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/dcmgr"
	"github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/wrapper"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"
	microconfig "github.com/Ankr-network/dccn-hub/common/micro_config"

	log "github.com/micro/go-log"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

func StartHandler(db dbservice.DBService, conf *microconfig.Config) {
	// New Service
	service := micro.NewService(
		micro.Name(conf.ServerName),
		micro.Version(conf.ServerVersion),
		micro.RegisterTTL(time.Second*time.Duration(conf.RegisterTTL)),
		micro.RegisterInterval(time.Second*time.Duration(conf.RegisterInterval)),
		micro.WrapHandler(wrapper.AuthWrapper),
	)

	// Initialise service
	service.Init()

	// Register Handler
	go_micro_srv_dcmgr.RegisterDcMgrHandler(service.Server(), handler.NewDcMgrHandler(db))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbConfig, err := dbcommon.LoadFromEnv()
	if err != nil {
		println(err.Error())
		return
	}

	microConfig, err := microconfig.LoadFromEnv()
	if err != nil {
		println(err.Error())
		return
	}

	db, err := dbservice.New(dbConfig)
	if err != nil {
		println(err.Error())
		return
	}
	defer db.Close()

	StartHandler(db, microConfig)
}
