package handler

import (
	"context"
	"errors"

	"github.com/Ankr-network/refactor/util"
	pb "github.com/Ankr-network/refactor/app_dccn_dccenter/proto"
)

type DcCenterHandler struct{}

func (p *DcCenterHandler) DataCenterList(ctx context.Context, req *pb.DataCenterListRequest, rsp *pb.DataCenterListResponse) error {
	token := req.Usertoken
	// TODO: Abstract DBConfig interface
	user := util.GetUser(token)
	util.WriteLog("task list reqeust")

	if user.ID == 0 {
		util.WriteLog("task list reqeust fail for taskmgr token error")
		return errors.New("task list reqeust fail for taskmgr token error")
	} else {
		dataCenters := util.DataCeterList()

		var dcList []*pb.DataCenterInfo
		for i := range dataCenters {
			dataCenter := dataCenters[i]
			dcInfo := &pb.DataCenterInfo{}
			dcInfo.Id = dataCenter.ID
			dcInfo.Name = dataCenter.Name
			dcList = append(dcList, dcInfo)
			//util.WriteLog("task id : %d %s status %s", task.ID,task.Name, task.Status)
		}

		rsp.DcList = dcList
		return nil
	}
}
