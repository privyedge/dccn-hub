package handler

import (
	"context"
	"io"
	"log"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
	dbservice "github.com/Ankr-network/dccn-hub/app-dccn-dcmgr/db_service"

	micro "github.com/micro/go-micro"
)

type DcMgrHandler struct {
	db             dbservice.DBService
	taskFeedback   micro.Publisher         // sync task information with task manager
	DcStreamCaches *DataCenterStreamCaches // hold all data center as cache
}

func New(db dbservice.DBService, feedback micro.Publisher) *DcMgrHandler {
	handler := &DcMgrHandler{
		db:             db,
		taskFeedback:   feedback,
		DcStreamCaches: NewDataCenterStreamCaches(),
	}
	return handler
}

func (p *DcMgrHandler) ServerStream(ctx context.Context, stream dcmgr.DCStreamer_ServerStreamStream) error {

	log.Println("Debug into ServerStream")
	for {
		in, err := stream.Recv()
		log.Println("Recv new request")
		if err == io.EOF {
			log.Println(err.Error())
			return nil
		}
		if err != nil {
			log.Println(err.Error())
			return err
		}

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
	}
}

func (p *DcMgrHandler) All() error {
	return nil
}

func (p *DcMgrHandler) Available() error {
	return nil
}

func (p *DcMgrHandler) Cleanup() {
	if p.DcStreamCaches != nil {
		p.DcStreamCaches.Cleanup()
	}
}
