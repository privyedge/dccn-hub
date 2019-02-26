package main

import (
	"log"

	micro "github.com/micro/go-micro"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/handler"
	"github.com/Ankr-network/dccn-hub/app-dccn-usermgr/token"

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

	startHandler(db)
}

// Init starts handler to listen.
func Init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	if conf, err = config.Load(); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Load config %+v\n", conf)
}

func startHandler(db dbservice.DBService) {
	// New Service
	srv := grpc.NewService(
		micro.Name(ankr_default.UserMgrRegistryServerName),
	)

	// Initialise service
	srv.Init()

	// New Publisher to deploy new task action.
	pubEmail := micro.NewPublisher(ankr_default.MQMail, srv.Client())

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	opt := srv.Server().Options()
	if err := opt.Broker.Connect(); err != nil {
		log.Fatal(err.Error())
	}

	userHandler := handler.New(db, token.New(), pubEmail)
	defer userHandler.Destroy()

	// Register Handler
	if err := usermgr.RegisterUserMgrHandler(srv.Server(), userHandler); err != nil {
		log.Fatal(err.Error())
	}

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
