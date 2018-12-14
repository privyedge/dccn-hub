package api_listner

import (
	"fmt"
	"github.com/Ankr-network/dccn-hub/util"
	pb "github.com/Ankr-network/dccn-rpc/protocol_new/cli"
	"github.com/Ankr-network/dccn-rpc/server_rpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/reflection"
	"log"

	//"net"
	"os"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	name string
}

func (s *server) TaskDetail(ctx context.Context, in *pb.TaskDetailRequest) (*pb.TaskDetailResponse, error) {
	fmt.Printf("received task detail request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		fmt.Printf("add new task fail for user token error \n")

		return &pb.TaskDetailResponse{Body: ""}, nil
	}
	task := util.GetTask(int(in.Taskid))
	if task.Userid != user.ID { // can not get other user task
		return &pb.TaskDetailResponse{Body: ""}, nil
	}
	return &pb.TaskDetailResponse{Body: task.URL}, nil

}

func (s *server) AddTask(ctx context.Context, in *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	fmt.Printf("received add task request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	if user.ID == 0 {
		fmt.Printf("add new task fail for user token error \n")
		return &pb.AddTaskResponse{Status: "Failure", Taskid: -1}, nil
	} else {

		// check datacenter name valid
		if len(in.Datacenter) != 0 {
			dc := util.GetDatacenter(in.Datacenter)

			if dc.ID == 0 {
				fmt.Printf("add new task fail for datacenter name does not exist \n")
				return &pb.AddTaskResponse{Status: "Failure", Taskid: -1}, nil

			}
		}
		//end check datacenter name

		task := util.Task{Name: in.Name, Datacenter: in.Datacenter, Userid: user.ID, Type: in.Type}
		id := util.AddTask(task)
		task.ID = id

		tastName := util.GetTaskNameAsTaskIDForK8s(task)
		util.UpdateTaskUnqueName(int(id), tastName)

		if in.Replica <= 0 || in.Replica > 100 {
			fmt.Printf("replica is eror %d use default 1  \n", in.Replica)
			in.Replica = 1
		}

		util.UpdateTaskReplica(int(id), int(in.Replica))
		e := util.Event{}
		e.Type = util.NewTaskEvent
		e.TaskID = int(task.ID)

		util.Send(util.TaskManagerQueueName, e)

		return &pb.AddTaskResponse{Status: "Success", Taskid: id}, nil
	}

}

func (s *server) TaskList(ctx context.Context, in *pb.TaskListRequest) (*pb.TaskListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)
	fmt.Printf("task list reqeust \n")

	if user.ID == 0 {
		fmt.Printf("task list reqeust fail for user token error\n")
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
			//fmt.Printf("task id : %d %s status %s \n", task.ID,task.Name, task.Status)
		}

		return &pb.TaskListResponse{Tasksinfo: taskList}, nil
	}

}

func (s *server) DataCenterList(ctx context.Context, in *pb.DataCenterListRequest) (*pb.DataCenterListResponse, error) {
	token := in.Usertoken
	user := util.GetUser(token)
	fmt.Printf("task list reqeust \n")

	if user.ID == 0 {
		fmt.Printf("task list reqeust fail for user token error\n")
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
			//fmt.Printf("task id : %d %s status %s \n", task.ID,task.Name, task.Status)
		}

		return &pb.DataCenterListResponse{DcList: dcList}, nil
	}

}

func (s *server) CancelTask(ctx context.Context, in *pb.CancelTaskRequest) (*pb.CancelTaskResponse, error) {
	fmt.Printf("received cancel task request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		fmt.Printf("can not find task\n")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	if user.ID == 0 {
		fmt.Printf("cancel task fail for user token error\n")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	if task.Userid != user.ID {
		fmt.Printf("task uid != user id \n")
		return &pb.CancelTaskResponse{Status: "Failure"}, nil
	}

	fmt.Printf("task %d in DataCenter %d \n", task.ID, int(task.Datacenterid))

	//sendMessageToK8(taskType string, taskid int64, name string, extra string)
	// todo test this function
	//datacenter := s.dcstreams[int(task.Datacenterid)]
	//if datacenter == nil {
	//	fmt.Printf("can not find datacenter \n")
	//	util.UpdateTask(int(in.Taskid), "cancelfailed", 0)
	//	return &pb.CancelTaskResponse{Status: "Failure"}, nil
	//
	//} else {
	//	fmt.Printf("send cancel message to datacenter id  %d \n", int(task.Datacenterid))
	//	util.UpdateTask(int(in.Taskid), "cancelling", 0)

	//if sendMessageToK8(datacenter, "CancelTask", in.Taskid, task.Uniquename, task.Name, task.Replica, "") == false {
	//	delete(s.dcstreams, int(task.Datacenterid))
	//}
	e := util.Event{}
	e.Type = util.CancelTaskEvent
	e.TaskID = int(task.ID)

	util.Send(util.TaskManagerQueueName, e)

	return &pb.CancelTaskResponse{Status: "Success"}, nil
	//}

}

func (s *server) UpdateTask(ctx context.Context, in *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	fmt.Printf("received update task request\n")
	token := in.Usertoken
	user := util.GetUser(token)

	task := util.GetTask(int(in.Taskid))

	if task.ID == 0 {
		fmt.Printf("can not find task\n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if user.ID == 0 {
		fmt.Printf("cancel task fail for user token error\n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if task.Userid != user.ID {
		fmt.Printf("task uid != user id \n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	if len(task.Uniquename) == 0 {
		fmt.Printf("task does not have Uniquename in mongodb \n")
		return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	}

	fmt.Printf("task %d in DataCenter %d \n", task.ID, int(task.Datacenterid))

	//datacenter := s.dcstreams[int(task.Datacenterid)]
	//if datacenter == nil {
	//	fmt.Printf("can not find datacenter \n")
	//	util.UpdateTask(int(in.Taskid), "updateFailed", 0)
	//	return &pb.UpdateTaskResponse{Status: "Failure"}, nil
	//} else {

	//check replica is valid
	if task.Replica == 0 { // support previous
		task.Replica = 1
	}

	if in.Replica <= 0 || in.Replica > 100 {
		fmt.Printf("replica is eror %d \n", in.Replica)
		in.Replica = int64(task.Replica)
	}

	if len(in.Name) == 0 {
		in.Name = task.Name
	}

	fmt.Printf("send update message to datacenter id  %d  replica  %d  image : %s\n", int(task.Datacenterid), int(in.Replica), task.Name)

	util.UpdateTask(int(in.Taskid), "updating", 0)

	// if they are same use 0 as default value
	if int(in.Replica) == task.Replica {
		in.Replica = 0
	}

	if in.Name == task.Name {
		in.Name = ""
	}

	e := util.Event{}
	e.Type = util.UpdateTaskEvent
	e.TaskID = int(task.ID)
	e.Replica = int(in.Replica)
	e.Name = in.Name

	util.Send(util.TaskManagerQueueName, e)

	//if sendMessageToK8(datacenter, "UpdateTask", in.Taskid, task.Uniquename, in.Name, int(in.Replica), "") == false {
	//	delete(s.dcstreams, int(task.Datacenterid))
	//}
	return &pb.UpdateTaskResponse{Status: "Success"}, nil
	//}

}

func StartService() {
	if len(os.Args) == 2 {
		util.MongoDBHost = os.Args[1]

	}

	util.WriteLog("this is a test")

	lis, s := server_rpc.Connect(port)
	ss := server{}

	pb.RegisterDccncliServer(s, &ss)
	// Register reflection service on gRPC server.

	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
