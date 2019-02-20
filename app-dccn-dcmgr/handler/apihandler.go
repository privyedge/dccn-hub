package handler

import (
	"context"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1/micro"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-hub/app-dccn-dcmgr/db_service"
	"log"
)

type DcMgrAPIHandler struct {
	db dbservice.DBService
}

func NewAPIHandler(db dbservice.DBService) *DcMgrAPIHandler {
	handler := &DcMgrAPIHandler{
		db: db,
	}
	return handler
}

func (p *DcMgrAPIHandler) DataCenterList(
	ctx context.Context, req *common_proto.Empty, rsp *dcmgr.DataCenterListResponse) error {
	//
	log.Println("api service receive DataCenterList from client")

	if list, err :=  p.db.GetAll(); err != nil {
		log.Println(err.Error())
		log.Println("DataCenterList failure")
		return err
	} else {
		log.Printf("DataCenterList successfully count: %d", len(*list))
		rsp.DcList = *list
	}
	return nil
}




