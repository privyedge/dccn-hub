package handler

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos/common"
	"log"
)

func (p *DcMgrHandler) updateTask(task *common_proto.Task) error {

	log.Printf("into updateTask from datacenter msg  : %v ", task)
	return p.taskFeedback.Publish(context.TODO(), task)
}
