package main

import (
	"github.com/Ankr-network/refactor/app_dccn_account/handler"
	"github.com/Ankr-network/refactor/app_dccn_account/proto"
	"github.com/Ankr-network/refactor/app_dccn_account/config"
	"github.com/Ankr-network/refactor/util"
	"github.com/Ankr-network/refactor/util/db"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"gopkg.in/mgo.v2"
	"time"
)

var (
	configPath = "config/account.json"
	conf *config.Config
	db *mgo.Session
	err error
)

func main() {

	StartConfigWatcher(configPath)
	defer conf.Finalize()

	if db, err = utildb.CreateDBConnection(conf.DB); err != nil {
		util.WriteLog(err.Error())
		return
	}

	StartHandler(db.DB(conf.DBName).C(conf.Collection))
	defer db.Close()
}

// StartConfigWatcher loads config and watch the update
func StartConfigWatcher(path string) {
	conf, err = config.New(path)
	if err != nil {
		util.WriteLog(err.Error())
		return
	}

	go conf.Watch(path, func() {
		tmp, err := utildb.CreateDBConnection(conf.DB)
		if err != nil {
			util.WriteLog("new db failed: " +  err.Error())
		}

		// TODO: ensure the collection setting right
		// tmp.Ensure()
		db, tmp = tmp, db
		tmp.Close()
	})
}

// StartHandler startes hander to listen.
func StartHandler(table *mgo.Collection) {
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
	accountmgr.RegisterAccountMgrHandler(service.Server(), handler.New(table))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
