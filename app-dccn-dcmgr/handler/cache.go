package handler

import (
	"log"
	"math/rand"
	"sync"
	"time"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
)

var heartbeat = &common_proto.Event{EventType: common_proto.Operation_HEARTBEAT}

type DataCenterStreamCaches struct {
	mu *sync.Mutex
	// TODO: Redis here
	streams map[string]dcmgr.DCStreamer_ServerStreamStream
}

func NewDataCenterStreamCaches() *DataCenterStreamCaches {

	cache := &DataCenterStreamCaches{
		mu:      new(sync.Mutex),
		streams: make(map[string]dcmgr.DCStreamer_ServerStreamStream),
	}
	go cache.checkHealthy()

	return cache
}

func (p *DataCenterStreamCaches) Add(dc *common_proto.DataCenter, stream dcmgr.DCStreamer_ServerStreamStream) {

	log.Println("Debug into DataCenterStreamCaches'Add")
	p.mu.Lock()
	p.streams[dc.Name] = stream
	p.mu.Unlock()
}

func (p *DataCenterStreamCaches) Remove(dc string) {

	log.Println("Debug into DataCenterStreamCaches'Remove")
	p.mu.Lock()
	if _, ok := p.streams[dc]; ok {
		delete(p.streams, dc)
	}
	p.mu.Unlock()
}

func (p *DataCenterStreamCaches) Len(dc string) int {
	p.mu.Lock()
	defer p.mu.Unlock()

	return len(p.streams)
}

func (p *DataCenterStreamCaches) One() (dcmgr.DCStreamer_ServerStreamStream, error) {

	log.Println("Debug into DataCenterStreamCaches'SelectFreeDataCenter")

	p.mu.Lock()
	defer p.mu.Unlock()

	if len(p.streams) <= 0 {
		return nil, ankr_default.ErrNoAvailableDataCenter
	}

	dcs := make([]string, len(p.streams))
	var i = 0
	for dc := range p.streams {
		dcs[i] = dc
		i++
	}

	randIndex := rand.Intn(len(dcs))
	return p.streams[dcs[randIndex]], nil
}

func (p *DataCenterStreamCaches) All() map[string]dcmgr.DCStreamer_ServerStreamStream {
	return p.streams
}

func (p *DataCenterStreamCaches) checkHealthy() {
	for {
		for dc, stream := range p.streams {
			if err := stream.Send(heartbeat); err != nil {
				p.Remove(dc)
				log.Println(stream.Close())
				// log.Println(p.db.UpdateStatus(dc, common_proto.Status_UNAVAILABLE))
			}
		}
		time.Sleep(time.Second * time.Duration(ankr_default.HeartBeatInterval))
	}
}
