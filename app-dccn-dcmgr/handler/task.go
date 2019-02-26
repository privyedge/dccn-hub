package handler

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos/common"
	"log"
)

func (p *DcMgrHandler) updateTask(stream *common_proto.DCStream) error {

	log.Printf("into updateTask from datacenter msg  : %v ", stream)
	return p.taskFeedback.Publish(context.TODO(), stream)
}
