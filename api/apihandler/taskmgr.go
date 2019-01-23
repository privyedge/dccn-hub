package apihandler

import (
	"context"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/dcmgr/v1"
	"github.com/Ankr-network/dccn-common/protos/taskmgr/v1"

	"github.com/micro/go-micro/client"
)

type ApiTask struct {
	api dcmgr.DCStreamerService
}

func (*ApiTask) AddTask(context.Context, *taskmgr.AddTaskRequest, *taskmgr.AddTaskResponse) error {
	panic("implement me")
}

func (*ApiTask) TaskList(context.Context, *taskmgr.ID, *taskmgr.TaskListResponse) error {
	panic("implement me")
}

func (*ApiTask) CancelTask(context.Context, *taskmgr.Request, *common_proto.Error) error {
	panic("implement me")
}

func (*ApiTask) PurgeTask(context.Context, *taskmgr.Request, *common_proto.Error) error {
	panic("implement me")
}

func (*ApiTask) TaskDetail(context.Context, *taskmgr.Request, *taskmgr.TaskDetailResponse) error {
	panic("implement me")
}

func (*ApiTask) UpdateTask(context.Context, *taskmgr.UpdateTaskRequest, *common_proto.Error) error {
	panic("implement me")
}

func NewApiTask(c client.Client) *ApiTask {
	return &ApiTask{
		api: dcmgr.NewDCStreamerService(ankr_default.DcMgrRegistryServerName, c),
	}
}
