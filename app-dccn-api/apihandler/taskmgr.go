package apihandler

import (
	"context"
	"log"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"

	"github.com/micro/go-micro/client"
)

type ApiTask struct {
	api taskmgr.TaskMgrService
}

func (p *ApiTask) CreateTask(ctx context.Context, req *taskmgr.CreateTaskRequest, rsp *taskmgr.CreateTaskResponse) error {
	out, _ := p.api.CreateTask(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) TaskList(ctx context.Context, req *taskmgr.ID, rsp *taskmgr.TaskListResponse) error {
	log.Println("Debug into TaskList")
	out, _ := p.api.TaskList(ctx, req)
	*rsp = *out
	rsp.Tasks = append(rsp.Tasks, out.Tasks...)
	return nil
}

func (p *ApiTask) CancelTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {
	log.Println("Debug into CancelTask")
	out, _ := p.api.CancelTask(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) PurgeTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {
	log.Println("Debug into PurgeTask")
	out, _ := p.api.PurgeTask(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) TaskDetail(ctx context.Context, req *taskmgr.Request, rsp *taskmgr.TaskDetailResponse) error {
	log.Println("Debug into TaskDetail")
	out, _ := p.api.TaskDetail(ctx, req)
	*rsp = *out
	return nil
}

func (p *ApiTask) UpdateTask(ctx context.Context, req *taskmgr.UpdateTaskRequest, rsp *common_proto.Error) error {
	log.Println("Debug into UpdateTask")
	out, _ := p.api.UpdateTask(ctx, req)
	*rsp = *out
	return nil
}

func NewApiTask(c client.Client) *ApiTask {
	return &ApiTask{
		api: taskmgr.NewTaskMgrService(ankr_default.TaskMgrRegistryServerName, c),
	}
}
