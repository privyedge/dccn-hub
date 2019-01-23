package apihandler

import (
	"context"

	"github.com/micro/go-micro/client"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1"
)

type ApiDataCenter struct {
	api dcmgr.DCStreamerService
}

func (*ApiDataCenter) ServerStream(context.Context, dcmgr.DCStreamer_ServerStreamStream) error {
	panic("implement me")
}

func NewApiDataCenter(c client.Client) *ApiDataCenter {
	return &ApiDataCenter{
		api: dcmgr.NewDCStreamerService(ankr_default.DcMgrRegistryServerName, c),
	}
}
