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

func (p *ApiTask) CreateTask(
	ctx context.Context, req *taskmgr.CreateTaskRequest, rsp *taskmgr.CreateTaskResponse) error {

	log.Println("Debug into CreateTask")
	out, err := p.api.CreateTask(ctx, req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*rsp = *out
	return nil
}

func (p *ApiTask) TaskList(ctx context.Context, req *taskmgr.TaskListRequest, rsp *taskmgr.TaskListResponse) error {
	log.Println("Debug into TaskList")
	out, err := p.api.TaskList(ctx, req)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	*rsp = *out
	return nil
}

func (p *ApiTask) CancelTask(ctx context.Context, req *taskmgr.TaskID, rsp *common_proto.Empty) error {
	log.Println("Debug into CancelTask")
	if _, err := p.api.CancelTask(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *ApiTask) PurgeTask(ctx context.Context, req *taskmgr.TaskID, rsp *common_proto.Empty) error {
	log.Println("Debug into PurgeTask")
	if _, err := p.api.PurgeTask(ctx, req); err != nil {
		log.Print(err.Error())
		return err
	}

	return nil
}

// func (p *ApiTask) TaskDetail(ctx context.Context, req *taskmgr.Request, rsp *taskmgr.TaskDetailResponse) error {
// 	log.Println("Debug into TaskDetail")
// 	if out, err := p.api.TaskDetail(ctx, req); err != nil {
// 		log.Println(err.Error())
// 		return err
// 	} else {
// 		*rsp = *out
// 	}
// 	return nil
// }

func (p *ApiTask) UpdateTask(ctx context.Context, req *taskmgr.UpdateTaskRequest, rsp *common_proto.Empty) error {
	log.Println("Debug into UpdateTask")
	if _, err := p.api.UpdateTask(ctx, req); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func NewApiTask(c client.Client) *ApiTask {
	return &ApiTask{
		api: taskmgr.NewTaskMgrService(ankr_default.TaskMgrRegistryServerName, c),
	}
}
