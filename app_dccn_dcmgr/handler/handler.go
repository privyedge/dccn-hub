package handler

import (
	"context"
	"github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/db_service"
	"github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/dcmgr"
)

type DcMgrHandler struct {
	db dbservice.DBService
}

func NewDcMgrHandler(db dbservice.DBService) *DcMgrHandler {
	return &DcMgrHandler{db}
}

func (p *DcMgrHandler) Get(ctx context.Context, id *go_micro_srv_dcmgr.ID, center *go_micro_srv_dcmgr.DataCenter) error {
	var err error
	center, err = p.db.Get(id.Id)
	return err
}

func (p *DcMgrHandler) Add(ctx context.Context, center *go_micro_srv_dcmgr.DataCenter, rsp *go_micro_srv_dcmgr.Response) error {
	return p.db.Add(*center)
}

func (p *DcMgrHandler) Update(ctx context.Context, center *go_micro_srv_dcmgr.DataCenter, rsp *go_micro_srv_dcmgr.Response) error {
	return p.db.Update(center)
}

