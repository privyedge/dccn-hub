package main

import (
	"log"
	"time"

	"github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/db_service"
	taskfilter "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/filter"
	"github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/handler"
	pb "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/proto/taskmgr"
	"github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/subscriber"
	"github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/wrapper"

	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/client"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	configPath = "config.json"
)

// StartHandler starts hander to listen.
func StartHandler(db dbservice.DBService, conf *config.Config) {
	// New Service
	service := micro.NewService(
		micro.Name(conf.SrvName),
		micro.Version(conf.Version),
		micro.RegisterTTL(time.Second*time.Duration(conf.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(conf.Interval)),
		micro.WrapHandler(wrapper.AuthWrapper),
	)

	// Initialise service
	service.Init()

	// New Publisher filter with policy
	newPub := micro.NewPublisher(conf.TopicPubNewTask, client.NewClient(client.Wrap(taskfilter.NewACWrapper)))
	// Cancel Publisher broadcast
	cancelPub := micro.NewPublisher(conf.TopicPubCancelTask, client.DefaultClient)

	// Register Function as Subscriber to update task by feedback.
	micro.RegisterSubscriber(conf.TopicSubResultTask, service.Server(), subscriber.New(db))

	// Register Handler
	pb.RegisterTaskMgrHandler(service.Server(), handler.New(db, newPub, cancelPub))

	// Run service
	if err := service.Run(); err != nil {
		log.Println(err.Error())
	}
}

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
