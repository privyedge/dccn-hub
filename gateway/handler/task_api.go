package handler

import (
	"context"

	taskmgr "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/proto/taskmgr"
	pb "github.com/Ankr-network/dccn-hub/gateway/proto/taskmgr"
)

type TaskApi struct {
	client taskmgr.TaskMgrService
}

func NewTaskApi(client taskmgr.TaskMgrService) *TaskApi {
	return &TaskApi{client: client}
}

func taskSrvToApi(task *taskmgr.Task) *pb.Task {
	var t pb.Task
	t.Id = task.Id
	t.Result = task.Result
	t.EventId = task.EventId
	t.DataCenterId = task.DataCenterId
	t.ReturnTime = task.ReturnTime
	t.UpdateTime = task.UpdateTime
	t.CancelTime = task.CancelTime
	t.CreateTime = task.CreateTime
	t.Status = pb.Task_Status(task.Status)
	t.UserId = task.UserId
	t.Extra = task.Extra
	t.Type = task.Type
	t.StartupTime = task.StartupTime
	t.Name = task.Name

	return &t
}

func taskApiToSrv(t *pb.Task) *taskmgr.Task {
	var task taskmgr.Task
	task.Id = t.Id
	task.Result = t.Result
	task.EventId = t.EventId
	task.DataCenterId = t.DataCenterId
	task.ReturnTime = t.ReturnTime
	task.UpdateTime = t.UpdateTime
	task.CancelTime = t.CancelTime
	task.CreateTime = t.CreateTime
	task.Status = taskmgr.Task_Status(t.Status)
	task.UserId = t.UserId
	task.Extra = t.Extra
	task.Type = t.Type
	task.StartupTime = t.StartupTime
	task.Name = t.Name

	return &task
}

func (p *TaskApi) Get(ctx context.Context, id *pb.ID, task *pb.Task) error {
	t, err := p.client.Get(ctx, &taskmgr.ID{Id: id.Id})
	if err != nil {
		return err
	}
	task = taskSrvToApi(t)
	return nil
}

func (p *TaskApi) Create(ctx context.Context, task *pb.Task, rsp *pb.Response) error {
	if _, err := p.client.Create(ctx, taskApiToSrv(task)); err != nil {
		return err
	}
	return nil
}

func (p *TaskApi) Cancel(ctx context.Context, id *pb.ID, rsp *pb.Response) error {
	if _, err := p.client.Cancel(ctx, &taskmgr.ID{Id: id.Id}); err != nil {
		return err
	}
	return nil
}

func (p *TaskApi) Update(ctx context.Context, task *pb.Task, rsp *pb.Response) error {
	if _, err := p.client.Update(ctx, taskApiToSrv(task)); err != nil {
		return err
	}
	return nil
}
