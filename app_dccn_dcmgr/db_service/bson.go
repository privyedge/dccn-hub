package dbservice

import (
	"errors"
	"github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/dcmgr"
)

type bsonDataCenter struct {
	Id                   int64    `json:"id,omitempty" bson:"_id"`
	Name                 string   `json:"name,omitempty" bson:"name,omitempty"`
	Status               int32    `json:"status,omitempty" bson:"status,omitempty"`
}

func (bcenter *bsonDataCenter) Encode(center *go_micro_srv_dcmgr.DataCenter) error {
	if bcenter == nil || center == nil {
		return errors.New("Null Pointer")
	}
	bcenter.Id = center.Id
	bcenter.Name = center.Name
	bcenter.Status = center.Status

	return nil
}

func (bcenter *bsonDataCenter) Decode(center *go_micro_srv_dcmgr.DataCenter) error {
	if bcenter == nil || center == nil {
		return errors.New("Null Pointer")
	}
	center.Id = bcenter.Id
	center.Name = bcenter.Name
	center.Status = bcenter.Status
	return nil
}

