package subscriber

import (
	"context"
	"github.com/Ankr-network/refactor/app_dccn_taskmgr/proto"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/metadata"
)

type TaskMgrSubscriber struct{}

func (e *TaskMgrSubscriber) Handle(ctx context.Context, event *taskmgr.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("Handler Received message: ", event, md)
	return nil
}

func Handler(ctx context.Context, event *taskmgr.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("Handler Received message: ", event, md)
	return nil
}
