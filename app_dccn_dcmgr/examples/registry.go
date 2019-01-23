package main

import (
	"context"
	"log"

	dcmgr "github.com/Ankr-network/dccn-common/protos/common/proto/v1"
	"github.com/micro/go-micro/client"
)

func isEqual(origin, dbdc *dcmgr.DataCenter) bool {
	return origin.Id == dbdc.Id && origin.Status == dbdc.Status && origin.Name == dbdc.Name
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// log.Println("app-dccn-v1-dcmgr-client service start...")
	// srv := grpc.NewService()
	// // srv := micro.NewService()
	// log.Printf("%#v\n", srv.Client())

	// srv.Init()

	datacenter := &dcmgr.DataCenter{
		Id:     1,
		Name:   "dc000",
		Status: 1,
	}

	// cli := dcmgr.NewDcMgrService("go.micro.srv.v1.dcmgr", srv.Client())
	cli := dcmgr.NewDcMgrService("go.micro.srv.v1.dcmgr", client.NewClient())
	// cli := dcmgr.NewDcMgrService("", srv.Client())
	// cli := dcmgr.NewDcMgrService("", client.NewClient())
	if _, err := cli.Create(context.Background(), datacenter); err != nil {
		log.Fatal(err.Error())
	}
	log.Println("Create Ok")

	c, err := cli.Get(context.Background(), &dcmgr.ID{Id: datacenter.Id})
	if err != nil {
		log.Fatal(err.Error())
	}

	if !isEqual(datacenter, c) {
		log.Fatalf("want: %#v\n, but %#v\n", datacenter, c)
	}
	log.Println("Get Ok.")
}
