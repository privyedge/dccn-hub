package subscriber

import (
	"context"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/metadata"

	"github.com/Ankr-network/refactor/app_dccn_k8s/proto"
)

type K8sSubscriber struct{}

func (e *K8sSubscriber) Handle(ctx context.Context, event *k8s.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("Handler Received message: ", event, md)
	return nil
}

func Handler(ctx context.Context, event *k8s.Event) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("Handler Received message: ", event, md)
	return nil
}
