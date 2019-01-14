package main

import (
	"log"

	"github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/db_service"
	"github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/handler"
	pb "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/v1"

	grpc "github.com/micro/go-grpc"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	conf       config.Config
	configPath string = "config.json"
	db         dbservice.DBService
	err        error
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
	log.Println("app_dccn_usermgr service start...")

	if conf, err = config.Load(configPath); err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("Load config %+v\n", conf)
}

func startHandler(db dbservice.DBService) {
	srv := grpc.NewService()
	srv.Init()

	// Register Handler
	pb.RegisterDcMgrHandler(srv.Server(), handler.NewDcMgrHandler(db))

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
