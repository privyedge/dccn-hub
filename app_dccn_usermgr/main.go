package main

import (
	"log"

	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/db_service"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/handler"
	pb "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/token"

	"github.com/micro/cli"
	micro "github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	conf config.Config
)

func main() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)

	srv := Init()

	db, err := dbservice.New(conf.DB)
	if err != nil {
		println(err.Error())
		return
	}
	defer db.Close()

	startHandler(srv, db)
}

// startHandler starts handler to listen.
func Init() micro.Service {
	// New Service
	srv := micro.NewService(
		micro.Flags(
			cli.BoolFlag{
				Name:  "TEST_FLAG",
				Usage: "for test usage",
			},
			cli.StringFlag{
				Name:  "DB_HOST",
				Usage: "DB URL",
			},
			cli.StringFlag{
				Name:  "DB_NAME",
				Usage: "Database name",
			},
			cli.StringFlag{
				Name:  "DB_COLLECTION",
				Usage: "Collection name",
			},
			cli.IntFlag{
				Name:  "DB_TIMEOUT",
				Usage: "Connect DB timeout value",
			},
			cli.IntFlag{
				Name:  "DB_POOL_LIMIT",
				Usage: "Max DB connections",
			},
			// JWT
			cli.IntFlag{
				Name:  "TOKEN_ACTIVE_TIME",
				Usage: "JWT ActiveTie",
			},
			cli.IntFlag{
				Name:  "TOKEN_NOT_BEFORE",
				Usage: "JWT NotBefore",
			},
		),
	)
	conf.Token = token.DefaultTokenConfig()
	var testFlag bool
	// Initialise service
	srv.Init(
		micro.Action(func(ctx *cli.Context) {
			testFlag = ctx.Bool("TEST_FLAG")
			// DB
			conf.DB.Host = ctx.String("DB_HOST")
			conf.DB.DB = ctx.String("DB_NAME")
			conf.DB.Collection = ctx.String("DB_COLLECTION")
			conf.DB.Timeout = ctx.Int("DB_TIMEOUT")
			conf.DB.PoolLimit = ctx.Int("DB_POOL_LIMIT")
			// TOKEN
			conf.Token.ActiveTime = ctx.Int("TOKEN_ACTIVE_TIME")
			conf.Token.NotBefore = int64(ctx.Int("TOKEN_NOT_BEFORE"))
		}),
	)

	if testFlag {
		// md5 -s "ankr_network_test_usermgr_db"
		conf.DB.DB = "114feb0961f8edfa8f514b67c6ef8af3"
	}
	return srv
}

func startHandler(srv micro.Service, db dbservice.DBService) {
	// Register Handler
	pb.RegisterUserMgrHandler(srv.Server(), handler.New(db, token.New(&conf.Token)))

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
