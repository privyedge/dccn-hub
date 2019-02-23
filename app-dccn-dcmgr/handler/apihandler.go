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


func (p *DcMgrAPIHandler) DataCenterLeaderBoard(ctx context.Context, req *common_proto.Empty,
	                              rsp *dcmgr.DataCenterLeaderBoardResponse ) error {
	//rsp = & dcmgr.DataCenterLeaderBoardResponse{}
	list := make([]*dcmgr.DataCenterLeaderBoardDetail, 0)
	{
		detail := dcmgr.DataCenterLeaderBoardDetail{}
		detail.Name = "us_cloud"
		detail.Number = 99.81
		list = append(list, &detail)
	}

	{
		detail := dcmgr.DataCenterLeaderBoardDetail{}
		detail.Name = "asia_cloud"
		detail.Number = 97.71
		list = append(list, &detail)
	}

	{
		detail := dcmgr.DataCenterLeaderBoardDetail{}
		detail.Name = "europe_cloud"
		detail.Number = 96.89
		list = append(list, &detail)
	}

	rsp.List = list
	log.Printf("DataCenterLeaderBoard %+v", rsp.List)
	return nil
}


func (p *DcMgrAPIHandler) NetworkInfo(ctx context.Context, req *common_proto.Empty,
	rsp *dcmgr.NetworkInfoResponse) error{

	rsp.UserCount = 299
	rsp.ContainerCount = 1342
	rsp.EnvironmentCount = 450
	rsp.HostCount = 137
    return nil
}
