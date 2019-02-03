package handler

import (
	"fmt"
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
	mu *sync.RWMutex
	// TODO: Redis here
	streams map[string]dcmgr.DCStreamer_ServerStreamStream
}

func NewDataCenterStreamCaches() *DataCenterStreamCaches {

	cache := &DataCenterStreamCaches{
		mu:      new(sync.RWMutex),
		streams: make(map[string]dcmgr.DCStreamer_ServerStreamStream),
	}
	go cache.checkHealthy()

	return cache
}

func (p *DataCenterStreamCaches) Add(dc *common_proto.DataCenter, stream dcmgr.DCStreamer_ServerStreamStream) {

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, ok := p.streams[dc.Name]; ok {
		return
	}

	log.Println("Debug into DataCenterStreamCaches'Add")
	p.streams[dc.Name] = stream
}

func (p *DataCenterStreamCaches) Remove(dc string) {

	p.mu.Lock()
	if _, ok := p.streams[dc]; ok {
		log.Println("Debug into DataCenterStreamCaches'Remove: ", dc)
		delete(p.streams, dc)
	}
	p.mu.Unlock()
}

func (p *DataCenterStreamCaches) Len(dc string) int {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return len(p.streams)
}

func (p *DataCenterStreamCaches) Has(dc string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	_, ok := p.streams[dc]
	return ok
}

func (p *DataCenterStreamCaches) One(dc string) (dcmgr.DCStreamer_ServerStreamStream, error) {

	log.Println("Debug into DataCenterStreamCaches'SelectFreeDataCenter")

	p.mu.RLock()
	defer p.mu.RUnlock()

	switch dc {
	case "":
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

	default:
		return p.Get(dc)
	}
}

func (p *DataCenterStreamCaches) Get(dc string) (dcmgr.DCStreamer_ServerStreamStream, error) {

	p.mu.RLock()
	defer p.mu.RUnlock()
	if ss, ok := p.streams[dc]; ok {
		return ss, nil
	}
	return nil, fmt.Errorf("%s not avaiable", dc)
}

func (p *DataCenterStreamCaches) All() map[string]dcmgr.DCStreamer_ServerStreamStream {
	p.mu.RLock()
	defer p.mu.RUnlock()

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

func (p *DataCenterStreamCaches) Cleanup() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, stream := range p.streams {
		if err := stream.Close(); err != nil {
			log.Fatal(err.Error())
		}
	}
}
