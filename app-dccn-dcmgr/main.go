package main

import (
	"log"

	micro "github.com/micro/go-micro"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"

	"github.com/Ankr-network/dccn-hub/app-dccn-dcmgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-dcmgr/db_service"
	"github.com/Ankr-network/dccn-hub/app-dccn-dcmgr/handler"

	grpc "github.com/micro/go-grpc"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	conf config.Config
	db   dbservice.DBService
	err  error
)

func main() {
	Init()

	if db, err = dbservice.New(conf.DB); err != nil {
		log.Fatal(err.Error())
	}
	defer db.Close()

	startHandler()
}

// Init starts handler to listen.
func Init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if conf, err = config.Load(); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Load config %+v\n", conf)
}

func startHandler() {
	srv := grpc.NewService(
		micro.Name(ankr_default.DcMgrRegistryServerName),
	)

	// Initialise service
	srv.Init()

	// New Publisher to deploy new task action.
	taskFeedback := micro.NewPublisher(ankr_default.MQFeedbackTask, srv.Client())

	dcHandler := handler.New(db, taskFeedback)
	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	opt := srv.Server().Options()
	opt.Broker.Connect()
	if err := micro.RegisterSubscriber(ankr_default.MQDeployTask, srv.Server(), dcHandler); err != nil {
		log.Fatal(err.Error())
	}

	// Register Handler
	if err := dcmgr.RegisterDCStreamerHandler(srv.Server(), dcHandler); err != nil {
		log.Fatal(err.Error())
	}

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
