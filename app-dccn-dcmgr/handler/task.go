package handler

import (
	"context"
	"log"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
)

func (p *DcMgrHandler) updateTask(event *common_proto.Event) error {

	log.Printf("into updateTask from datacenter msg  : %v ", event)
	return p.taskFeedback.Publish(context.TODO(), event)
}
