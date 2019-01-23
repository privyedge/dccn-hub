package main

import (
	"context"
	"log"
	"time"

	dcmgr "github.com/Ankr-network/dccn-common/protos/common/proto/v1"
	"github.com/micro/go-micro/client"
)

var serverAddress string = "localhost:50052"

func isEqual(origin, dbdc *dcmgr.DataCenter) bool {
	return origin.Id == dbdc.Id &&
		origin.Name == dbdc.Name &&
		origin.Report == dbdc.Report &&
		origin.LastReportTime == dbdc.LastReportTime
}

func main() {
	opts := client.WithAddress(serverAddress)

	// c := client.DefaultClient
	// c.Init(...client.Option)
	// log.Printf("%#v, %#v\n", c, s2)

	// log.Printf("%#v, %#v\n", c, s2)

	// s2 := grpc.NewService()
	// c2 := s2.Client()
	// log.Printf("%#v, %#v\n", c, c2, s2)
	// This is recommended to call rpc api
	cli := dcmgr.NewDcMgrService("go.micro.srv.v1.dcmgr", client.DefaultClient)
	datacenter := &dcmgr.DataCenter{
		Id:             1,
		Name:           "dc000",
		Status:         1,
		Report:         "test mock case",
		LastReportTime: int32(time.Now().Unix()),
	}
	rsp, err := cli.Create(context.TODO(), datacenter, opts)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(rsp)

	dc, err := cli.Get(context.TODO(), &dcmgr.ID{Id: 1}, opts)
	if err != nil {
		log.Fatal(err.Error())
	}

	if !isEqual(datacenter, dc) {
		log.Fatalf("want %#v, but %#v\n", datacenter, dc)
	}

	log.Println("OK")
}
