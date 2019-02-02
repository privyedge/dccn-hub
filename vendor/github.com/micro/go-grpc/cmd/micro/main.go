package main

import (
	"os"

	_ "github.com/micro/go-plugins/client/grpc"
	_ "github.com/micro/go-plugins/server/grpc"
	"github.com/micro/micro/cmd"
)

func init() {
	os.Setenv("MICRO_CLIENT", "grpc")
	os.Setenv("MICRO_SERVER", "grpc")
}

func main() {
	cmd.Init()
}
