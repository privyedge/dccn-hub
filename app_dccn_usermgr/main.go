package main

import (
	"github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"gopkg.in/mgo.v2"
	"time"

	"github.com/Ankr-network/refactor/util"
	"github.com/Ankr-network/refactor/util/db"
	"github.com/Ankr-network/refactor/proto/usermgr"
	"github.com/Ankr-network/refactor/app_dccn_usermgr/token"
	"github.com/Ankr-network/refactor/app_dccn_usermgr/config"
	"github.com/Ankr-network/refactor/app_dccn_usermgr/db_user"
	"github.com/Ankr-network/refactor/app_dccn_usermgr/handler"
)

var (
	configPath = "refactor/app_dccn_usermgr/config/config.json"
	conf *config.Config
	db *mgo.Session
	err error
)

func main() {

	LoadConfig(configPath)
	defer conf.Finalize()

	if db, err = utildb.CreateDBConnection(conf.DBConfig); err != nil {
		util.WriteLog(err.Error())
		return
	}

	StartHandler(db.DB(conf.DB).C(conf.Collection), conf)
	defer db.Close()
}

// LoadConfig loads config and watch the update
func LoadConfig(path string) {
	conf, err = config.New(path)
	if err != nil {
		util.WriteLog(err.Error())
		panic(err.Error())
	}
}

// StartHandler startes hander to listen.
func StartHandler(c *mgo.Collection, conf *config.Config) {
	// New Service
	service := k8s.NewService(
		micro.Name("network.ankr.srv.user"),
		micro.Version("user daemon"),
		micro.RegisterTTL(time.Second*time.Duration(conf.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(conf.Interval)),
	)

	// Initialise service
	service.Init()

	// Register Handler
	usermgr.RegisterUserMgrHandler(service.Server(), handler.New(dbuser.New(c), token.New(&conf.TokenConfig)))

	// Maybe send message service here

	// Run service
	if err := service.Run(); err != nil {
		util.WriteLog(err.Error())
	}
}
