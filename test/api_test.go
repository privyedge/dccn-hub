package test

import (
	"fmt"
	"github.com/Ankr-network/dccn-common"
	pb "github.com/Ankr-network/dccn-common/protocol/cli"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"testing"
	"time"
)

var address = ":" + ankr_const.DefaultPort

var totalTask = 0
var c pb.DccncliClient
var ctx context.Context
var cancel context.CancelFunc
var taskID int64

const replicaValue int64 = 2

func ClientCreate() {
	if len(os.Args) == 3 {
		address = os.Args[2] + ":" + ankr_const.DefaultPort
	} else if len(os.Args) == 2 {
		address = os.Args[1] + ":" + ankr_const.DefaultPort
	}

	fmt.Printf("connect Ankr Hub with address %s \n", address)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return
	}
	c = pb.NewDccncliClient(conn)
	ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)

}

func TestTaskList(t *testing.T) {
	ClientCreate()
	r, err := c.TaskList(ctx, &pb.TaskListRequest{Usertoken: ankr_const.DefaultUserToken})
	if err != nil {
		msg := fmt.Sprintf("TestTaskList : could not send: %v", err)
		t.Error(msg)
		return
	}

	totalTask = len(r.Tasksinfo)

}

func TestCreateTask(t *testing.T) {
	r, err := c.AddTask(ctx, &pb.AddTaskRequest{Name: "nginx:1.12", Usertoken: ankr_const.DefaultUserToken})
	if err != nil {
		t.Error("failed")
		return
	}

	if r.Status != ankr_const.CliReplyStatusSuccess {
		t.Error("reply status error")
	}

	taskID = r.Taskid

	time.Sleep(time.Second * 1)
	status := getTaskStatus(taskID).Status

	if status != ankr_const.TaskStatusRunning {
		msg := fmt.Sprintf("TestCreateTask status error, expect running  : %s", status)
		t.Error(msg)
	}

	if getTaskListCount() != totalTask+1 {
		t.Error("TestCreateTask  total task list error after create task")
	}

}

func TestUpdateTask(t *testing.T) {

	r, err := c.UpdateTask(ctx, &pb.UpdateTaskRequest{Taskid: taskID, Replica: replicaValue, Name: "", Usertoken: ankr_const.DefaultUserToken})
	if err != nil {
		msg := fmt.Sprintf("Client: could not send: %v", err)
		t.Error(msg)
	}

	if r.Status != ankr_const.CliReplyStatusSuccess {
		t.Error("reply status error")
	}

	time.Sleep(time.Second * 1)
	replica := getTaskStatus(taskID).Replica

	if replicaValue != replica {
		msg := fmt.Sprintf("TestUpdateTask replica error, expect %d  : actual %d", replicaValue, replica)
		t.Error(msg)
	}

}

func TestCancelTask(t *testing.T) {
	r, err := c.CancelTask(ctx, &pb.CancelTaskRequest{Taskid: taskID, Usertoken: ankr_const.DefaultUserToken})
	if err != nil {
		msg := fmt.Sprintf("Client: could not send: %v", err)
		t.Error(msg)
	}
	if r.Status != ankr_const.CliReplyStatusSuccess {
		t.Error("reply status error")
	}

	time.Sleep(time.Second * 1)
	status := getTaskStatus(taskID).Status

	if status != ankr_const.TaskStatusCancelled {
		msg := fmt.Sprintf("TestCancelTask task status error, expect %s  : actual %s", ankr_const.TaskStatusCancelled, status)
		t.Error(msg)
	}

}

func TestPurgeTask(t *testing.T) {
	r, err := c.PurgeTask(ctx, &pb.PurgeTaskRequest{Taskid: taskID, Usertoken: ankr_const.DefaultUserToken})
	if err != nil {
		msg := fmt.Sprintf("Client: could not send: %v", err)
		t.Error(msg)
	}
	if r.Status != ankr_const.CliReplyStatusSuccess {
		t.Error("reply status error")
	}

	if getTaskListCount() != totalTask {
		t.Error("TestPurgeTask  total task list error after purge task")
	}

}

func TestDataCenterList(t *testing.T) {
	r, err := c.DataCenterList(ctx, &pb.DataCenterListRequest{Usertoken: ankr_const.DefaultUserToken})
	if err != nil {
		msg := fmt.Sprintf("Client: could not send: %v", err)
		t.Error(msg)
	}

	dcList := r.DcList
	if checkAvailabeDataCenter(dcList) == false {
		t.Error("No DataCenter Available")
	}

}

func checkAvailabeDataCenter(dcList []*pb.DataCenterInfo) bool {
	for i := range dcList {
		dataCetner := dcList[i]
		if dataCetner.Status == ankr_const.DataCenteStatusOnLine {
			return true
		}
	}

	return false
}

func getTaskStatus(taskID int64) *pb.TaskInfo {
	r, err := c.TaskDetail(ctx, &pb.TaskDetailRequest{Taskid: taskID, Usertoken: "ed1605e17374bde6c68864d072c9f5c9"})
	if err != nil {
		fmt.Printf("Client: could not send: %v", err)
		return &pb.TaskInfo{}
	}

	return r.Taskinfo
}

func getTaskListCount() int {
	r, err := c.TaskList(ctx, &pb.TaskListRequest{Usertoken: ankr_const.DefaultUserToken})
	if err != nil {
		return 0
	}
	return len(r.Tasksinfo)
}
