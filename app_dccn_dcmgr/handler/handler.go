package handler

import (
	"context"

	dbservice "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/db_service"
	dcmgr "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/v1"
)

type DcMgrHandler struct {
	db dbservice.DBService
}

func NewDcMgrHandler(db dbservice.DBService) *DcMgrHandler {
	return &DcMgrHandler{db}
}

func (p *DcMgrHandler) Get(ctx context.Context, id *dcmgr.ID, center *dcmgr.DataCenter) error {
	c, err := p.db.Get(id.Id)
	if err != nil {
		return err
	}

	*center = *c
	return err
}

func (p *DcMgrHandler) Create(ctx context.Context, center *dcmgr.DataCenter, rsp *dcmgr.Response) error {
	return p.db.Create(center)
}
