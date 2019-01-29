package main

import (
	"log"

	grpc "github.com/micro/go-grpc"
	micro "github.com/micro/go-micro"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"

	mail "github.com/Ankr-network/dccn-common/protos/email/v1/micro"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"

	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
	"github.com/Ankr-network/dccn-hub/app-dccn-api/apihandler"
	"github.com/Ankr-network/dccn-hub/app-dccn-dcmgr/handler"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-usermgr/db_service"

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
	// New Service
	srv := grpc.NewService()

	// Initialise service
	srv.Init()

	// Register User Handler
	if err := usermgr.RegisterUserMgrHandler(srv.Server(), apihandler.NewApiUser(srv.Client())); err != nil {
		log.Fatal(err.Error())
	}

	// Register Task Handler
	if err := taskmgr.RegisterTaskMgrHandler(srv.Server(), apihandler.NewApiTask(srv.Client())); err != nil {
		log.Fatal(err.Error())
	}

	// Dc Manager register handler
	// New Publisher to deploy new task action.
	taskFeedback := micro.NewPublisher(ankr_default.MQFeedbackTask, srv.Client())

	dcHandler := handler.New(taskFeedback)

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	if err := micro.RegisterSubscriber(ankr_default.MQDeployTask, srv.Server(), dcHandler); err != nil {
		log.Fatal(err.Error())
	}

	// Register Dc Manager Handler
	if err := dcmgr.RegisterDCStreamerHandler(srv.Server(), dcHandler); err != nil {
		log.Fatal(err.Error())
	}

	// Register Email Handler
	if err := mail.RegisterMailHandler(srv.Server(), apihandler.NewApiEmail(srv.Client())); err != nil {
		log.Fatal(err.Error())
	}

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
