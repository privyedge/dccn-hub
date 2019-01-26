package handler

import (
	"context"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func (p *DcMgrHandler) updateTask(event *common_proto.Event) error {
	return p.taskFeedback.Publish(context.TODO(), event)
}
