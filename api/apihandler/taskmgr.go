package apihandler

import (
	"context"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	"github.com/Ankr-network/dccn-common/protos/taskmgr/v1"

	"github.com/micro/go-micro/client"
)

type ApiTask struct {
	api taskmgr.TaskMgrService
}

func (p *ApiTask) AddTask(ctx context.Context, req *taskmgr.AddTaskRequest, rsp *taskmgr.AddTaskResponse) error {
	out, _ := p.api.AddTask(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) TaskList(ctx context.Context, req *taskmgr.ID, rsp *taskmgr.TaskListResponse) error {
	out, _ := p.api.TaskList(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) CancelTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {
	out, _ := p.api.CancelTask(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) PurgeTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {
	out, _ := p.api.PurgeTask(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) TaskDetail(ctx context.Context, req *taskmgr.Request, rsp *taskmgr.TaskDetailResponse) error {
	out, _ := p.api.TaskDetail(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) UpdateTask(ctx context.Context, req *taskmgr.UpdateTaskRequest, rsp *common_proto.Error) error {
	out, _ := p.api.UpdateTask(ctx, req)
	*rsp = *out
	return nil
}

func NewApiTask(c client.Client) *ApiTask {
	return &ApiTask{
		api: taskmgr.NewTaskMgrService(ankr_default.DcMgrRegistryServerName, c),
	}
}
