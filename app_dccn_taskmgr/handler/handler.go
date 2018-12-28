package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/Ankr-network/dccn-rpc"
	"github.com/Ankr-network/refactor/app_dccn_taskmgr/proto"
	"github.com/Ankr-network/refactor/util"
	"github.com/micro/go-micro"
)

type TaskMgrHandler struct{
	name      string
	newPub micro.Publisher
	cancelPub micro.Publisher
}

func New(newPub, cancelPub micro.Publisher) *TaskMgrHandler {
	return &TaskMgrHandler{
		newPub:   newPub,
		cancelPub: cancelPub,
	}
}

func (p *TaskMgrHandler) AddTask(ctx context.Context, req *taskmgr.AddTaskRequest, rsp *taskmgr.AddTaskResponse) error {
		util.WriteLog("received add task request")
	token := req.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		util.WriteLog("add new task fail for user token error")
		rsp.Status = ankr_const.CliReplyStatusFailure
		rsp.Taskid = -1

		return errors.New("add new task fail for user token error")
	} else {

		// check datacenter name valid
		if len(req.Datacenter) != 0 {
			dc := util.GetDatacenter(req.Datacenter)

			if dc.ID == 0 {
				util.WriteLog("add new task fail for datacenter name does not exist")
				rsp.Status = ankr_const.CliReplyStatusFailure
				rsp.Taskid = -1

				return errors.New("add new task fail for datacenter name does not exist")
			}
		}
		//end check datacenter name

		task := util.Task{Name: req.Name, Datacenter: req.Datacenter, Userid: user.ID, Type: req.Type}
		id := util.AddTask(task)
		task.ID = id

		tastName := util.GetTaskNameAsTaskIDForK8s(task)
		util.UpdateTaskUnqueName(int(id), tastName)

		if req.Replica <= 0 || req.Replica > 100 {
			logStr := fmt.Sprintf("replica is eror %d use default 1 ", req.Replica)
			util.WriteLog(logStr)
			req.Replica = 1
		}

		util.UpdateTaskReplica(int(id), int(req.Replica))
		e := util.Event{}
		e.Type = ankr_const.NewTaskEvent
		e.TaskID = int(task.ID)

		if err := p.newPub.Publish(context.Background(), e); err != nil {
			util.WriteLog(err.Error())
			return err
		}

		rsp.Status = ankr_const.CliReplyStatusSuccess
		rsp.Taskid = id
		return nil
	}
}

// TaskList get all task list from DBConfig
func (p *TaskMgrHandler) TaskList(ctx context.Context, in *taskmgr.TaskListRequest, rsp *taskmgr.TaskListResponse) error {
		token := in.Usertoken
	user := util.GetUser(token)
	util.WriteLog("task list reqeust")

	if user.ID == 0 {
		util.WriteLog("task list reqeust fail for user token error")
		return errors.New("task list reqeust fail for user token error")
	} else {
		tasks := util.TaskList(int(user.ID))

		dcs := util.GetDatacentersMap()

		var taskList []*taskmgr.TaskInfo
		for i := range tasks {
			task := tasks[i]
			taskInfo := &taskmgr.TaskInfo{}
			taskInfo.Taskid = task.ID
			taskInfo.Taskname = task.Name
			taskInfo.Status = task.Status
			taskInfo.Replica = int64(task.Replica)
			taskInfo.Datacenter = dcs[task.Datacenterid]
			if len(taskInfo.Datacenter) == 0 {
				taskInfo.Datacenter = task.Datacenter // for user assign datacenter name but not startsuccess
			}
			taskList = append(taskList, taskInfo)
			//util.WriteLog("task id : %d %s status %s", task.ID,task.Name, task.Status)
		}

		rsp.Tasksinfo = taskList
		return nil
	}
}

// CancelTask broadcast
func (p *TaskMgrHandler) CancelTask(ctx context.Context, in *taskmgr.CancelTaskRequest, rsp *taskmgr.CancelTaskResponse) error {
		util.WriteLog("received cancel task request")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		util.WriteLog("can not find task")
		rsp.Status = "Failure"
		return errors.New("can not find task")
	}

	if user.ID == 0 {
		util.WriteLog("cancel task fail for user token error")
		rsp.Status = "Failure"
		return errors.New("cancel task fail for user token error")
	}

	if task.Userid != user.ID {
		util.WriteLog("task uid != user id")
		rsp.Status = "Failure"
		return errors.New("task uid != user id")
	}

	logStr := fmt.Sprintf("task %d in DataCenter %d", task.ID, int(task.Datacenterid))
	util.WriteLog(logStr)

	if task.Status != ankr_const.TaskStatusCancelled {
		e := util.Event{}
		e.Type = ankr_const.CancelTaskEvent
		e.TaskID = int(task.ID)
		// util.Send(ankr_const.TaskManagerQueueName, e)
		if err := p.cancelPub.Publish(context.Background(), e); err != nil {
			util.WriteLog(err.Error())
			return err
		}

	}

	rsp.Status = "Success"
	return nil
}

// UpdateTask first cancel task, then add new task
func (p *TaskMgrHandler) UpdateTask(ctx context.Context, in *taskmgr.UpdateTaskRequest, rsp *taskmgr.UpdateTaskResponse) error {
		util.WriteLog("received update task request")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		util.WriteLog("can not find task")
		rsp.Status = "Failure"
		return errors.New("can not find task")
	}

	if user.ID == 0 {
		util.WriteLog("cancel task fail for user token error")
		rsp.Status = "Failure"
		return errors.New("cancel task fail for user token error")
	}

	if task.Userid != user.ID {
		util.WriteLog("task uid != user id")
		rsp.Status = "Failure"
		return errors.New("task uid != user id")
	}

	if len(task.Uniquename) == 0 {
		util.WriteLog("task does not have Uniquename in mongodb")
		rsp.Status = ankr_const.CliReplyStatusFailure
		return errors.New("task does not have Uniquename in mongodb")
	}

	logStr := fmt.Sprintf("task %d in DataCenter %d", task.ID, int(task.Datacenterid))
	util.WriteLog(logStr)

	//check replica is valid
	if task.Replica == 0 { // support previous
		task.Replica = 1
	}

	if in.Replica <= 0 || in.Replica > 100 {
		logStr := fmt.Sprintf("replica is eror %d", in.Replica)
		util.WriteLog(logStr)
		in.Replica = int64(task.Replica)
	}

	if len(in.Name) == 0 {
		in.Name = task.Name
	}

	logStr2 := fmt.Sprintf("send update message to datacenter id  %d  replica  %d  image : %s", int(task.Datacenterid), int(in.Replica), task.Name)
	util.WriteLog(logStr2)

	util.UpdateTask(int(in.Taskid), ankr_const.TaskStatusRunning, 0)

	// if they are same use 0 as default value
	if int(in.Replica) == task.Replica {
		in.Replica = 0
	}

	if in.Name == task.Name {
		in.Name = ""
	}

	e := util.Event{}
	e.Type = ankr_const.UpdateTaskEvent
	e.TaskID = int(task.ID)
	e.Replica = int(in.Replica)
	e.Name = in.Name

	util.Send(ankr_const.TaskManagerQueueName, e)

	rsp.Status = ankr_const.CliReplyStatusSuccess
	return nil
}

// TaskDetail get the deatil info
func (p *TaskMgrHandler) TaskDetail(ctx context.Context, req *taskmgr.TaskDetailRequest, rsp *taskmgr.TaskDetailResponse) error {
	token := req.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		util.WriteLog("add new task fail for user token error")
		return errors.New("add new task fail for user token error")
	}
	task := util.GetTask(int(req.Taskid))
	if task.Userid != user.ID { // can not get other user task
		return errors.New("can not get other user task")
	}

	rsp.Body = task.URL
	return nil
}
