package handler

import (
	"context"
	"log"

	"gopkg.in/mgo.v2/bson"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	"github.com/google/uuid"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"

	micro "github.com/micro/go-micro"

	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
	db "github.com/Ankr-network/dccn-hub/app-dccn-taskmgr/db_service"
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

	if err := checkId(req.UserId, req.TaskId); err != nil {
		log.Println(err.Error())
		return err
	}
	task, err := p.checkOwner(req.UserId, req.TaskId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	rsp.Task = task

	return nil
}

func (p *TaskMgrHandler) CreateTask(ctx context.Context, req *taskmgr.CreateTaskRequest, rsp *taskmgr.CreateTaskResponse) error {

	log.Println("Debug into CreateTask")
	if req.UserId == "" {
		log.Println(ankr_default.ErrUserNotExist.Error())
		return ankr_default.ErrUserNotExist
	}

	if req.Task.Replica < 0 || req.Task.Replica >= 100 {
		log.Println(ankr_default.ErrReplicaTooMany)
		return ankr_default.ErrReplicaTooMany
	}

	req.Task.Status = common_proto.TaskStatus_START
	req.Task.UserId = req.UserId
	req.Task.Id = uuid.New().String()
	rsp.TaskId = req.Task.Id

	event := common_proto.Event{
		EventType: common_proto.Operation_TASK_CREATE,
		OpMessage: &common_proto.Event_Task{Task: req.Task},
	}

	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		log.Println(ankr_default.ErrPublish)
		return ankr_default.ErrPublish
	}

	if err := p.db.Create(req.Task); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// Must return nil for gRPC handler
func (p *TaskMgrHandler) CancelTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {

	log.Println("Debug into CancelTask")
	if err := checkId(req.UserId, req.TaskId); err != nil {
		log.Println(err.Error())
		return err
	}
	task, err := p.checkOwner(req.UserId, req.TaskId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if task.Status != common_proto.TaskStatus_RUNNING &&
		task.Status != common_proto.TaskStatus_START &&
		task.Status != common_proto.TaskStatus_UPDATING {
		log.Println(ankr_default.ErrStatusNotSupportOperation)
		return ankr_default.ErrStatusNotSupportOperation
	}

	event := common_proto.Event{
		EventType: common_proto.Operation_TASK_CANCEL,
		OpMessage: &common_proto.Event_Task{Task: task},
	}

	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		log.Println(err.Error())
		return err
	}

	if err := p.db.Update(task.Id, bson.M{"$set": bson.M{"status": common_proto.TaskStatus_CANCEL}}); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (p *TaskMgrHandler) TaskList(ctx context.Context, req *taskmgr.ID, rsp *taskmgr.TaskListResponse) error {

	log.Println("Debug into TaskList")

	if req.UserId == "" {
		log.Println(ankr_default.ErrUserNotExist)
		return ankr_default.ErrUserNotExist
	}

	tasks, err := p.db.GetAll(req.UserId)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	rsp.Tasks = append(rsp.Tasks, *tasks...)

	return nil
}

func (p *TaskMgrHandler) UpdateTask(ctx context.Context, req *taskmgr.UpdateTaskRequest, rsp *common_proto.Error) error {

	log.Println("Debug into UpdateTask")

	if err := checkId(req.UserId, req.Task.Id); err != nil {
		log.Println(err.Error())
		return err
	}
	task, err := p.checkOwner(req.UserId, req.Task.Id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if req.Task.Replica <= 0 || req.Task.Replica >= 100 {
		log.Println(ankr_default.ErrReplicaTooMany.Error())
		return ankr_default.ErrReplicaTooMany
	}

	if task.Status == common_proto.TaskStatus_CANCELLED ||
		task.Status == common_proto.TaskStatus_DONE {
		log.Println(ankr_default.ErrStatusNotSupportOperation.Error())
		return ankr_default.ErrStatusNotSupportOperation
	}

	event := common_proto.Event{
		EventType: common_proto.Operation_TASK_UPDATE,
		OpMessage: &common_proto.Event_Task{Task: task},
	}

	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		log.Println(err.Error())
		return err
	}
	// TODO: wait deamon notify
	req.Task.Status = common_proto.TaskStatus_UPDATING
	if err := p.db.UpdateTask(req.Task.Id, req.Task); err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}

func (p *TaskMgrHandler) PurgeTask(ctx context.Context, req *taskmgr.Request, rsp *common_proto.Error) error {

	log.Println("Debug into PurgeTask")
	return nil
}

func (p *TaskMgrHandler) checkOwner(userId, taskId string) (*common_proto.Task, error) {
	task, err := p.db.Get(taskId)
	if err != nil {
		return nil, err
	}

	if task.UserId != userId {
		return nil, ankr_default.ErrUserNotOwn
	}
	return task, nil
}

func checkId(userId, taskId string) error {
	if userId == "" {
		return ankr_default.ErrUserNotExist
	}

	if taskId == "" {
		return ankr_default.ErrUserNotOwn
	}

	return nil
}
