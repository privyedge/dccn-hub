package handler

import (
	"context"
	"fmt"
	"log"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	"github.com/pborman/uuid"
	"gopkg.in/mgo.v2/bson"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"

	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
	db "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/db_service"
	micro "github.com/micro/go-micro"
)

type TaskMgrHandler struct {
	db         db.DBService
	deployTask micro.Publisher
}

func New(db db.DBService, deployTask micro.Publisher) *TaskMgrHandler {

	return &TaskMgrHandler{
		db:         db,
		deployTask: deployTask,
	}
}

func (p *TaskMgrHandler) TaskDetail(ctx context.Context, req *taskmgr.Request, rsp *taskmgr.TaskDetailResponse) error {

	log.Println("Debug into TaskDetail")

	if !checkId(&rsp.Error, req.UserId, req.TaskId) {
		log.Println(rsp.Error.Details)
		return nil
	}
	task, ok := p.checkOwner(&rsp.Error, req.UserId, req.TaskId)
	if !ok {
		log.Println(rsp.Error.Details)
		return nil
	}

	rsp.Task = task

	return nil
}

func (p *TaskMgrHandler) CreateTask(ctx context.Context, req *taskmgr.CreateTaskRequest, rsp *taskmgr.CreateTaskResponse) error {

	if req.UserId == 0 {
		pbError(&rsp.Error, ankr_default.ErrUserNotExist)
		log.Println(rsp.Error.Details)
		return nil
	}

	if req.Task.Replica <= 0 || req.Task.Replica >= 100 {
		pbError(&rsp.Error, ankr_default.ErrReplicaTooMany)
		log.Println(rsp.Error.Details)
		return nil
	}

	event := common_proto.Event{
		EventType: common_proto.Operation_TASK_CREATE,
		OpMessage: &common_proto.Event_Task{Task: req.Task},
	}

	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		pbError(&rsp.Error, fmt.Errorf("%s: %s", ankr_default.ErrPublish, err.Error()))
		log.Println(rsp.Error.Details)
		return nil
	}
	req.Task.Status = common_proto.TaskStatus_START
	req.Task.Id = uuid.New()
	rsp.TaskId = req.Task.Id
	if err := p.db.Create(req.Task); err != nil {
		pbError(&rsp.Error, err)
		log.Println(rsp.Error.Details)
		return nil
	}

	return nil
}

// Must return nil for gRPC handler
func (p *TaskMgrHandler) CancelTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {
	log.Println("Debug into CancelTask")
	if !checkId(&rsp, req.UserId, req.TaskId) {
		log.Println(rsp.Details)
		return nil
	}
	task, ok := p.checkOwner(&rsp, req.UserId, req.TaskId)
	if !ok {
		log.Println(rsp.Details)
		return nil
	}

	if task.Status != common_proto.TaskStatus_RUNNING &&
		task.Status != common_proto.TaskStatus_START &&
		task.Status != common_proto.TaskStatus_UPDATING {
		pbError(&rsp, fmt.Errorf("%s: %s", ankr_default.ErrStatusNotSupportOperation, task.Status))
		log.Println(rsp.Details)
		return nil
	}

	event := common_proto.Event{
		EventType: common_proto.Operation_TASK_CANCEL,
		OpMessage: &common_proto.Event_Task{Task: task},
	}

	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		pbError(&rsp, err)
		log.Println(rsp.Details)
		return nil
	}

	if err := p.db.Update(task.Id, bson.M{"$set": bson.M{"status": common_proto.TaskStatus_CANCELL}}); err != nil {
		pbError(&rsp, err)
		log.Println(rsp.Details)
		return nil
	}

	return nil
}

func (p *TaskMgrHandler) TaskList(ctx context.Context, req *taskmgr.ID, rsp *taskmgr.TaskListResponse) error {
	log.Println("Debug into TaskList")

	if req.UserId == 0 {
		pbError(&rsp.Error, ankr_default.ErrUserNotExist)
		log.Println(rsp.Error.Details)
		return nil
	}

	tasks, err := p.db.GetAll(req.UserId)
	if err != nil {
		pbError(&rsp.Error, err)
		log.Println(rsp.Error.Details)
		return nil
	}
	rsp.Tasks = append(rsp.Tasks, *tasks...)

	return nil
}

func (p *TaskMgrHandler) UpdateTask(ctx context.Context, req *taskmgr.UpdateTaskRequest, rsp *common_proto.Error) error {
	log.Println("Debug into UpdateTask")

	if !checkId(&rsp, req.UserId, req.Task.Id) {
		log.Println(rsp.Details)
		return nil
	}
	task, ok := p.checkOwner(&rsp, req.UserId, req.Task.Id)
	if !ok {
		log.Println(rsp.Details)
		return nil
	}

	if req.Task.Replica <= 0 || req.Task.Replica >= 100 {
		pbError(&rsp, ankr_default.ErrReplicaTooMany)
		log.Println(rsp.Details)
		return nil
	}

	if task.Status == common_proto.TaskStatus_CANCELLED ||
		task.Status == common_proto.TaskStatus_DONE {
		pbError(&rsp, ankr_default.ErrStatusNotSupportOperation)
		log.Println(rsp.Details)
		return nil
	}

	event := common_proto.Event{
		EventType: common_proto.Operation_TASK_UPDATE,
		OpMessage: &common_proto.Event_Task{Task: task},
	}

	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		pbError(&rsp, err)
		log.Println(rsp.Details)
		return nil
	}
	// TODO: wait deamon notify
	req.Task.Status = common_proto.TaskStatus_UPDATING
	if err := p.db.UpdateTask(req.Task.Id, req.Task); err != nil {
		pbError(&rsp, err)
		log.Println(rsp.Details)
		return nil
	}
	return nil
}

func (p *TaskMgrHandler) PurgeTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {

	return nil
}

func pbError(rsp **common_proto.Error, err error) {

	if *rsp == nil {
		*rsp = &common_proto.Error{}
	}
	(*rsp).Status = common_proto.Status_FAILURE
	(*rsp).Details = err.Error()
}

func (p *TaskMgrHandler) checkOwner(rsp **common_proto.Error, userId int64, taskId string) (*common_proto.Task, bool) {
	task, err := p.db.Get(taskId)
	if err != nil {
		pbError(rsp, err)
		return nil, false
	}

	if task.UserId != userId {
		pbError(rsp, ankr_default.ErrUserNotOwn)
		return nil, false
	}
	return task, true
}

func checkId(rsp **common_proto.Error, userId int64, taskId string) bool {
	if userId == 0 {
		pbError(rsp, ankr_default.ErrUserNotExist)
		return false
	}

	if taskId == "" {
		pbError(rsp, ankr_default.ErrUserNotOwn)
		return false
	}

	return true
}
