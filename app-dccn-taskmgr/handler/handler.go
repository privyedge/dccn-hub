package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/Ankr-network/dccn-common/protos"
	"github.com/Ankr-network/dccn-common/protos/common"
	"github.com/google/uuid"
	"github.com/gorhill/cronexpr"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/metadata"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"

	"github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
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



type Token struct {
	Exp int64
	Jti string
	Iss string
}


func getUserID(ctx context.Context) string{
	meta, ok := metadata.FromContext(ctx)
	// Note this is now uppercase (not entirely sure why this is...)
	var token string
	if ok {
		token = meta["token"]
	}

	parts := strings.Split(token, ".")

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		fmt.Println("decode error:", err)

	}
	fmt.Println(string(decoded))

	var dat Token

	if err := json.Unmarshal(decoded, &dat); err != nil {
		panic(err)
	}

	return string(dat.Jti)
}

func (p *TaskMgrHandler) CreateTask(ctx context.Context, req *taskmgr.CreateTaskRequest, rsp *taskmgr.CreateTaskResponse) error {
	uid := getUserID(ctx)
	log.Println("task manager service CreateTask")

	if req.Task.Attributes.Replica < 0 || req.Task.Attributes.Replica >= 100 {
		log.Println(ankr_default.ErrReplicaTooMany)
		return ankr_default.ErrReplicaTooMany
	}

	log.Printf("CreateTask task %+v", req)

	if req.Task.Attributes.Replica == 0 {
		req.Task.Attributes.Replica = 1
	}

	if req.Task.Type == common_proto.TaskType_CRONJOB { // check schudule filed
		_ , err := cronexpr.Parse(req.Task.GetTypeCronJob().Schedule)
		if err != nil {
			log.Printf("check crobjob scheducle fomat error %s \n", err.Error())
			return ankr_default.ErrCronJobScheduleFormat
		}

	}

	req.Task.Status = common_proto.TaskStatus_STARTING
	req.Task.Id = uuid.New().String()
	rsp.TaskId = req.Task.Id


	event := common_proto.DCResponse{
		OpType: common_proto.DCOperation_TASK_CREATE,
		OpPayload: &common_proto.DCResponse_Task{Task: req.Task},
	}


	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		log.Println(ankr_default.ErrPublish)
		return ankr_default.ErrPublish
	}else{
		log.Println("task manager service send CreateTask MQ message to dc manager service (api)")
	}

	if err := p.db.Create(req.Task, uid); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

// Must return nil for gRPC handler
func (p *TaskMgrHandler) CancelTask(ctx context.Context, req *taskmgr.TaskID, rsp *common_proto.Empty) error {
	userId :=getUserID(ctx)
	log.Println("Debug into CancelTask")
	if err := checkId(userId, req.TaskId); err != nil {
		log.Println(err.Error())
		return err
	}
	task, err := p.checkOwner(userId, req.TaskId)
	if err != nil {
		log.Println(err.Error())
		return err
	}


	event := common_proto.DCResponse{
		OpType: common_proto.DCOperation_TASK_CANCEL,
		OpPayload: &common_proto.DCResponse_Task{Task: task},
	}

	if err := p.deployTask.Publish(context.Background(), &event); err != nil {
		log.Println(err.Error())
		return err
	}

	if err := p.db.Update(task.Id, bson.M{"$set": bson.M{"status": common_proto.TaskStatus_CANCELLED}}); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func convertToTaskMessage(task db.TaskRecord) common_proto.Task {
	message := common_proto.Task{}
	message.Id = task.ID
	message.Name = task.Name
	message.Type = task.Type
	message.Status = task.Status
	//deployMessage := common_proto.TaskTypeDeployment{Image : task.Image}
	if task.Type == common_proto.TaskType_DEPLOYMENT {
		t := common_proto.Task_TypeDeployment{TypeDeployment: &common_proto.TaskTypeDeployment{Image:task.Image}}
		message.TypeData = &t
	}

	if task.Type == common_proto.TaskType_JOB {
		t := common_proto.Task_TypeJob{TypeJob: &common_proto.TaskTypeJob{Image:task.Image}}
		message.TypeData = &t
	}

	if task.Type == common_proto.TaskType_CRONJOB {
		t := common_proto.Task_TypeCronJob{TypeCronJob: &common_proto.TaskTypeCronJob{Image:task.Image, Schedule: task.Schedule}}
		message.TypeData = &t
	}

	return message

}

func (p *TaskMgrHandler) TaskList(ctx context.Context, req *taskmgr.TaskListRequest, rsp *taskmgr.TaskListResponse) error {
	userId := getUserID(ctx)
	log.Println("task service into TaskList")

	tasks, err := p.db.GetAll(userId)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	tasksWithoutHidden := make([]*common_proto.Task, 0)

	for i := 0; i < len(tasks); i++ {
		if tasks[i].Hidden != true {
			taskMessage := convertToTaskMessage(tasks[i])
			tasksWithoutHidden = append(tasksWithoutHidden, &taskMessage)
		}
	}

	rsp.Tasks = tasksWithoutHidden

	return nil
}

func (p *TaskMgrHandler) UpdateTask(ctx context.Context, req *taskmgr.UpdateTaskRequest, rsp *common_proto.Empty) error {
	userId := getUserID(ctx)

	if err := checkId(userId, req.Task.Id); err != nil {
		log.Println(err.Error())
		return err
	}


	task, err := p.checkOwner(userId, req.Task.Id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	if req.Task.Attributes.Replica == 0 {
		req.Task.Attributes.Replica = task.Attributes.Replica
	}

	if req.Task.Attributes.Replica < 0 || req.Task.Attributes.Replica >= 100 {
		log.Println(ankr_default.ErrReplicaTooMany.Error())
		return ankr_default.ErrReplicaTooMany
	}

	if task.Status == common_proto.TaskStatus_CANCELLED ||
		task.Status == common_proto.TaskStatus_DONE {
		log.Println(ankr_default.ErrStatusNotSupportOperation.Error())
		return ankr_default.ErrStatusNotSupportOperation
	}

	event := common_proto.DCResponse{
		OpType: common_proto.DCOperation_TASK_UPDATE,
		OpPayload: &common_proto.DCResponse_Task{Task: task},
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

func (p *TaskMgrHandler) PurgeTask(ctx context.Context, req *taskmgr.TaskID, rsp *common_proto.Empty) error {
	error := p.CancelTask(ctx, req, rsp)
	if error == nil {
		log.Printf(" PurgeTask  %s \n", req.TaskId)
		p.db.Update(req.TaskId, bson.M{"$set": bson.M{"hidden": true}})
	}
	return error
}

func (p *TaskMgrHandler) checkOwner(userId, taskId string) (*common_proto.Task, error) {
	task, err := p.db.Get(taskId)

	if err != nil {
		return nil, err
	}

	log.Printf("taskid : %s user id -%s-   user_token_id -%s-  ", taskId, task.Userid, userId)

	if task.Userid != userId {
		return nil, ankr_default.ErrUserNotOwn
	}

	taskMessage := convertToTaskMessage(task)

	return &taskMessage, nil
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
