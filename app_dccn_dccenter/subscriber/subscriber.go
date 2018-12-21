package subscriber

import (
	"context"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/metadata"

	"github.com/Ankr-network/refactor/app_dccn_dccenter/proto"
)

type DcCenterSubscriber struct{}

func (e *DcCenterSubscriber) Handle(ctx context.Context, event *dccenter.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("Handler Received message: ", event, md)
	return nil
}

func Handler(ctx context.Context, event *dccenter.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("Handler Received message: ", event, md)
	return nil
}
