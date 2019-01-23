package socket

import (
	"github.com/Ankr-network/dccn-rpc/server_rpc"
	"google.golang.org/grpc"
)

type DCHandler struct {
	server *grpc.ClientConn
}

func New(port string) {
	lis, err := server_rpc.Connect(port)
	if err != nil {
		return err
	}
}
