package main

import (
	"log"
	"time"

	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/handler"
	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/token"

	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	configPath = "config.json"
)

func main() {
	conf, err := config.New(configPath)
	if err != nil {
		println(err.Error())
		return
	}

	db, err := dbservice.New(conf.DBConfig)
	if err != nil {
		println(err.Error())
		return
	}
	defer db.Close()

	StartHandler(db, conf)
}

// StartHandler starts hander to listen.
func StartHandler(db dbservice.DBService, conf *config.Config) {
	// New Service
	service := micro.NewService(
		micro.Name(conf.SrvName),
		micro.Version(conf.Version),
		micro.RegisterTTL(time.Second*time.Duration(conf.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(conf.Interval)),
	)

	// Initialise service
	service.Init()

	// Register Handler
	pb.RegisterUserMgrHandler(service.Server(), handler.New(db, token.New(&conf.TokenConfig)))

	// Run service
	if err := service.Run(); err != nil {
		log.Println(err.Error())
	}
}
