package handler

import (
	"log"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	dcmgr "github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
	"github.com/google/uuid"
)

func (p *DcMgrHandler) updateDataCenter(dc *common_proto.DataCenter, stream dcmgr.DCStreamer_ServerStreamStream) error {
	// first update database
	log.Printf("into updateDataCenter  : %v ", dc)

	center , err :=  p.db.GetByName(dc.Name)


	if center.Name == "" {
		// data center dose not exist, register it
		log.Printf("insert new datacenter  : %v ", dc)
		dc.Id = uuid.New().String()
		if err = p.db.Create(dc); err != nil {
			log.Println(err.Error(), ", ", *dc)
			return err
		}
	} else {
		log.Printf("update datacenter by name : %s  ", center.Name)
		if err = p.db.Update(dc); err != nil {
			log.Println(err.Error())
			return err
		}
	}

	// then update stream
	log.Printf("update new data center stream: %s ", dc.Name)
	p.DcStreamCaches.Add(dc, stream)

	return nil
}
