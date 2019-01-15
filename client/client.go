package main

import (
	"context"
	"log"

	hello "github.com/micro/examples/greeter/srv/proto/hello"
	"github.com/micro/go-micro/client"
)

var serverAddress string = "localhost:50051"

func main() {
	// This is recommended to call rpc api
	cli := hello.NewSayService("go.micro.srv.greeter", client.DefaultClient)
	rsp, err := cli.Hello(context.TODO(), &hello.Request{Name: "sanghaifa"}, client.WithAddress(serverAddress))
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Println(rsp)
}

// One of the usage
func CallMicroService() {
	req := client.NewRequest("go.micro.srv.greeter", "Say.Hello", &hello.Request{Name: "sanghai"})
	rsp := &hello.Response{}
	err := client.Call(context.TODO(), req, rsp, client.WithAddress(serverAddress))
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println(rsp)
}
