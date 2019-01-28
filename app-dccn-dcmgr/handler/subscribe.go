package handler

import (
	"context"
	"log"
	"math/rand"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
)

func (p *DcMgrHandler) selectFreeDataCenter() (dcmgr.DCStreamer_ServerStreamStream, error) {
	dcs := make([]string, 0, len(p.dcStreams))
	var i = 0
	for dc := range p.dcStreams {
		dcs[i] = dc
		i++
	}

	if len(dcs) <= 0 {
		return nil, ankr_default.ErrNoAvailableDataCenter
	}
	randIndex := rand.Intn(len(dcs))
	return p.dcStreams[dcs[randIndex]], nil
}

// UpdateTaskByFeedback receives task result from data center, returns to v1
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *DcMgrHandler) UpdateTaskByFeedback(ctx context.Context, event *common_proto.Event) error {
	switch event.EventType {
	case common_proto.Operation_TASK_CREATE, common_proto.Operation_TASK_CANCEL, common_proto.Operation_TASK_UPDATE:
		stream, err := p.selectFreeDataCenter()
		if err != nil {
			return err
		}
		if !p.send(stream, event) {
			log.Printf("%s: %v", ankr_default.ErrSyncTaskInfo.Error(), *event)
			return ankr_default.ErrSyncTaskInfo
		}
	default:
		return ankr_default.ErrUnknown
	}
	return nil
}
