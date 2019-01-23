package handler

import (
	"context"
	"io"
	"log"
	"sync"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/db_service"

	micro "github.com/micro/go-micro"
)

type DcMgrHandler struct {
	mu           sync.Mutex // protect dc streams
	db           dbservice.DBService
	taskFeedback micro.Publisher                                // sync task information with task manager
	dcStreams    map[string]dcmgr.DCStreamer_ServerStreamStream // hold all data center as cache
}

func New(feedback micro.Publisher) *DcMgrHandler {
	handler := &DcMgrHandler{
		taskFeedback: feedback,
		dcStreams:    make(map[string]dcmgr.DCStreamer_ServerStreamStream),
	}
	go handler.checkHealth()
	return handler
}

func (p *DcMgrHandler) ServerStream(ctx context.Context, stream dcmgr.DCStreamer_ServerStreamStream) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			log.Println(err.Error())
			return nil
		}
		if err != nil {
			log.Println(err.Error())
			return err
		}

		p.mu.Lock()
		switch in.EventType {
		case common_proto.Operation_HEARTBEAT: // update data center in cache
			if err := p.updateDataCenter(in.GetDataCenter(), stream); err != nil {
				log.Println(err.Error())
			}
		case common_proto.Operation_TASK_CREATE, common_proto.Operation_TASK_UPDATE, common_proto.Operation_TASK_CANCEL: // update task status
			if err := p.updateTask(in); err != nil {
				log.Println(err.Error())
			}
		default:
			log.Println(ankr_default.ErrUnknown.Error())
		}
		p.mu.Unlock()
	}
}
