package main

import (
	"context"
	"errors"
	"log"
	"os"

	grpc "github.com/micro/go-grpc"
	micro "github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/server"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	pb "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-hub/app-dccn-taskmgr/config"
	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-taskmgr/db_service"
	"github.com/Ankr-network/dccn-hub/app-dccn-taskmgr/handler"
	"github.com/Ankr-network/dccn-hub/app-dccn-taskmgr/subscriber"

	_ "github.com/micro/go-plugins/broker/rabbitmq"
)

var (
	srv  micro.Service
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
	log.Printf("load config: %#v\n", conf)
}

// StartHandler starts handler to listen.
func startHandler(db dbservice.DBService) {
	// var srv micro.Service
	// New Service
	srv = grpc.NewService(
		micro.Name(ankr_default.TaskMgrRegistryServerName),
		micro.WrapHandler(AuthWrapper),
	)

	// Initialise srv
	srv.Init()

	// New Publisher to deploy new task action.
	deployTask := micro.NewPublisher(ankr_default.MQDeployTask, srv.Client())

	// Register Function as TaskStatusFeedback to update task by data center manager's feedback.
	if err := micro.RegisterSubscriber(ankr_default.MQFeedbackTask, srv.Server(), subscriber.New(db)); err != nil {
		log.Fatal(err.Error())
	}

	// Register Handler
	if err := pb.RegisterTaskMgrHandler(srv.Server(), handler.New(db, deployTask)); err != nil {
		log.Fatal(err.Error())
	}

	// Run srv
	if err := srv.Run(); err != nil {
		log.Println(err.Error())
	}
}
func AuthWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		if os.Getenv("DISABLE_AUTH") == "true" {
			log.Println("disable auth")
			return fn(ctx, req, resp)
		}
		meta, ok := metadata.FromContext(ctx)
		if !ok {
			log.Println("no auth meta-data found in request")
			return errors.New("no auth meta-data found in request")
		}

		// Note this is now uppercase (not entirely sure why this is...)
		token := meta["token"]
		log.Println("Authenticating with token: ", token)

		// Auth here
		// Really shouldn't be using a global here, find a better way
		// of doing this, since you can't pass it into a wrapper.
		userMgrService := usermgr.NewUserMgrService(ankr_default.UserMgrRegistryServerName, srv.Client())
		rsp, _ := userMgrService.VerifyToken(context.Background(), &usermgr.Token{Token: token})
		if rsp != nil && rsp.Status == common_proto.Status_FAILURE {
			log.Println(rsp.Details)
			return errors.New(rsp.Details)
		}
		err = fn(ctx, req, resp)
		log.Println(err.Error())
		return err
	}
}
