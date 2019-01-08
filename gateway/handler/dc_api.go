package handler

import (
	"context"

	dcmgr "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/dcmgr"
	pb "github.com/Ankr-network/dccn-hub/gateway/proto/dcmgr"
)

type DataCenterApi struct {
	client dcmgr.DcMgrService
}

func NewDataCenterApi(client dcmgr.DcMgrService) *DataCenterApi {
	return &DataCenterApi{client: client}
}

func dataCenterSrvToApi(dataCenter *dcmgr.DataCenter) *pb.DataCenter {
	var dc pb.DataCenter
	dc.Id = dataCenter.Id
	dc.Name = dataCenter.Name
	dc.Status = dataCenter.Status

	return &dc
}

func dataCenterApiToSrv(dc *pb.DataCenter) *dcmgr.DataCenter {
	var dataCenter dcmgr.DataCenter
	dataCenter.Id = dc.Id
	dataCenter.Name = dc.Name
	dataCenter.Status = dc.Status

	return &dataCenter
}

func (p *DataCenterApi) Get(ctx context.Context, id *pb.ID, dc *pb.DataCenter) error {
	dataCenter, err := p.client.Get(ctx, &dcmgr.ID{Id: id.Id})
	if err != nil {
		return err
	}

	dc = dataCenterSrvToApi(dataCenter)
	return nil
}

func (p *DataCenterApi) Add(ctx context.Context, dc *pb.DataCenter, rsp *pb.Response) error {
	if _, err := p.client.Add(ctx, dataCenterApiToSrv(dc)); err != nil {
		return err
	}
	return nil
}

func (p *DataCenterApi) Update(ctx context.Context, dc *pb.DataCenter, rsp *pb.Response) error {
	if _, err := p.client.Update(ctx, dataCenterApiToSrv(dc)); err != nil {
		return err
	}
	return nil
}
