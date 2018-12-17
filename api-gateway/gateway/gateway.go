package api_listener

import (
	"fmt"
	"github.com/Ankr-network/dccn-hub/util"
	ankr_const "github.com/Ankr-network/dccn-rpc"
	pb "github.com/Ankr-network/dccn-rpc/protocol_new/cli"
	"github.com/Ankr-network/dccn-rpc/server_rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"log"
	"os"
)

const (
	port = ":50051"
)

type server struct {
	name string
}

func (s *server) TaskDetail(ctx context.Context, in *pb.TaskDetailRequest) (*pb.TaskDetailResponse, error) {
	util.WriteLog("received task detail request")
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		util.WriteLog("add new task fail for user token error")

		return &pb.TaskDetailResponse{Body: ""}, nil
	}
	task := util.GetTask(int(in.Taskid))
	if task.Userid != user.ID { // can not get other user task
		return &pb.TaskDetailResponse{Body: ""}, nil
	}
	return &pb.TaskDetailResponse{Body: task.URL}, nil

}

func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	util.WriteLog("received add task request")
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		util.WriteLog("add new task fail for user token error")
		return &pb.AddTaskResponse{Status: ankr_const.CliReplyStatusFailure, Taskid: -1}, nil
	} else {

		// check datacenter name valid
		if len(in.Datacenter) != 0 {
			dc := util.GetDatacenter(in.Datacenter)

			if dc.ID == 0 {
				util.WriteLog("add new task fail for datacenter name does not exist")
				return &pb.AddTaskResponse{Status: ankr_const.CliReplyStatusFailure, Taskid: -1}, nil

			}
		}
		//end check datacenter name

		task := util.Task{Name: in.Name, Datacenter: in.Datacenter, Userid: user.ID, Type: in.Type}
		id := util.AddTask(task)
		task.ID = id

		tastName := util.GetTaskNameAsTaskIDForK8s(task)
		util.UpdateTaskUnqueName(int(id), tastName)

		if in.Replica <= 0 || in.Replica > 100 {
			logStr := fmt.Sprintf("replica is eror %d use default 1 ", in.Replica)
			util.WriteLog(logStr)
			in.Replica = 1
		}

		util.UpdateTaskReplica(int(id), int(in.Replica))
		e := util.Event{}
		e.Type = ankr_const.NewTaskEvent
		e.TaskID = int(task.ID)

		util.Send(ankr_const.TaskManagerQueueName, e)

		return &pb.AddTaskResponse{Status: ankr_const.CliReplyStatusSuccess, Taskid: id}, nil
	}

}

func (s *server) TaskList(ctx context.Context, in *pb.TaskListRequest) (*pb.TaskListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)
	util.WriteLog("task list reqeust")

	if user.ID == 0 {
		util.WriteLog("task list reqeust fail for user token error")
		return &pb.TaskListResponse{}, nil
	} else {
		tasks := util.TaskList(int(user.ID))

		dcs := util.GetDatacentersMap()

		var taskList []*pb.TaskInfo
		for i := range tasks {
			task := tasks[i]
			taskInfo := &pb.TaskInfo{}
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

		return &pb.TaskListResponse{Tasksinfo: taskList}, nil
	}

}

func (s *server) DataCenterList(ctx context.Context, in *pb.DataCenterListRequest) (*pb.DataCenterListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)
	util.WriteLog("task list reqeust")

	if user.ID == 0 {
		util.WriteLog("task list reqeust fail for user token error")
		return &pb.DataCenterListResponse{}, nil
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

		return &pb.DataCenterListResponse{DcList: dcList}, nil
	}

}

func (s *server) CancelTask(ctx context.Context, in *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
	util.WriteLog("received cancel task request")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		util.WriteLog("can not find task")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	if user.ID == 0 {
		util.WriteLog("cancel task fail for user token error")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	if task.Userid != user.ID {
		util.WriteLog("task uid != user id")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	logStr := fmt.Sprintf("task %d in DataCenter %d", task.ID, int(task.Datacenterid))
	util.WriteLog(logStr)

	if task.Status != ankr_const.TaskStatusCancelled {
		e := util.Event{}
		e.Type = ankr_const.CancelTaskEvent
		e.TaskID = int(task.ID)
		util.Send(ankr_const.TaskManagerQueueName, e)

	}

	return &pb.CancelTaskResponse{Status: "Success"}, nil

}

func (s *server) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	util.WriteLog("received update task request")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		util.WriteLog("can not find task")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if user.ID == 0 {
		util.WriteLog("cancel task fail for user token error")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if task.Userid != user.ID {
		util.WriteLog("task uid != user id")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if len(task.Uniquename) == 0 {
		util.WriteLog("task does not have Uniquename in mongodb")
		return &pb.UpdateTaskResponse{Status: ankr_const.CliReplyStatusSuccess}, nil
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

	return &pb.UpdateTaskResponse{Status: ankr_const.CliReplyStatusSuccess}, nil

}

func StartService() {
	if len(os.Args) == 3 {
		util.MongoDBHost = os.Args[1]

	}

	if len(os.Args) == 3 {
		util.RabbitMQHost = os.Args[2]
	}

	util.WriteLog("Start API Listner ")

	lis, s := server_rpc.Connect(port)
	ss := server{}

	pb.RegisterDccncliServer(s, &ss)
	// Register reflection service on gRPC server.

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
