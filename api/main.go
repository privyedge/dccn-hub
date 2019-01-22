package main

import (
	"log"

	pb "github.com/Ankr-network/dccn-common/protos/usermgr/v1"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/handler"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/token"

	grpc "github.com/micro/go-grpc"
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
	srv := grpc.NewService()

	// Initialise service
	srv.Init()

	// Register Handler
	if err := pb.RegisterUserMgrHandler(srv.Server(), handler.New(srv.Client())); err != nil {
		log.Fatal(err.Error())
	}

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
