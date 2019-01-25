package handler

import (
	"log"
	"time"

	"github.com/google/uuid"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
)

func (p *DcMgrHandler) send(stream dcmgr.DCStreamer_ServerStreamStream, msg *common_proto.Event) bool {
	if stream != nil {
		if err := stream.Send(msg); err != nil {
			return false
		}
	}
	return true
}

func (p *DcMgrHandler) checkHealth() {
	for {
		for dc, stream := range p.dcStreams {
			if p.send(stream, &common_proto.Event{EventType: common_proto.Operation_HEARTBEAT}) == false {
				log.Printf("%s %s\n", dc, ankr_default.ErrConnection)
				delete(p.dcStreams, dc)
				log.Println(stream.Close())
				log.Println(p.db.UpdateStatus(dc, common_proto.Status_UNAVALIABLE))
			}
		}
		time.Sleep(time.Second * time.Duration(ankr_default.HeartBeatInterval))
	}
}

func (p *DcMgrHandler) updateDataCenter(dc *common_proto.DataCenter, stream dcmgr.DCStreamer_ServerStreamStream) error {
	if dc.Id == "" {
		// data center dose not exist, register it
		dc.Id = uuid.New().String()
		if err := p.db.Create(dc); err != nil {
			return err
		}
	} else {
		if err := p.db.Update(dc); err != nil {
			return err
		}
	}

	p.dcStreams[dc.Name] = stream
	return nil
}
