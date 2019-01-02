package subscriber

import (
	"context"
	"github.com/Ankr-network/refactor/proto"
	"github.com/micro/go-log"
	"github.com/micro/go-micro/metadata"
)

// GetResult receives task result from data center, returns to taskmgr
func GetResult(ctx context.Context, event *taskmgr.TaskResult) error {
	md, _ := metadata.FromContext(ctx)
	log.Logf("Received Result: ", event, md)
	// TODO: Deposit the result into DBConfig; mysql or mongodb
	return nil
}
