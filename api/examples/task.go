package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/micro/go-micro/metadata"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	apiCommon "github.com/Ankr-network/dccn-hub/api/examples/common"
)

var addr string
var token = "token"

func init() {
	addr = os.Getenv("API_ADDRESS")
	log.Println("API ADDRESS: ", addr)
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err.Error())
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			log.Println(err.Error())
		}
	}(conn)

	cl := taskmgr.NewTaskMgrClient(conn)

	tokenContext := metadata.NewContext(context.Background(), map[string]string{
		"Token": token,
	})

	tasks := apiCommon.MockTasks()
	log.Println("Test CreateTask")
	for i := range tasks {
		log.Println("Debug CreateTask ", i)
		if rsp, _ := cl.CreateTask(tokenContext, &taskmgr.CreateTaskRequest{UserId: tasks[i].UserId, Task: &tasks[i]}); apiCommon.IsSuccess("CreateTask", rsp.Error) {
			log.Println("CreateTask Ok")
		}
		return
	}

	log.Println("Test TaskList")
	var userTasks []*common_proto.Task
	if rsp, _ := cl.TaskList(tokenContext, &taskmgr.ID{UserId: "1"}); apiCommon.IsSuccess("TaskList", rsp.Error) {
		userTasks = append(userTasks, rsp.Tasks...)
		log.Println("TaskList Ok")
	}

	if len(userTasks) == 0 {
		log.Fatalf("no tasks belongs to %d\n", 1)
	}

	// CancelTask
	cancelTask := userTasks[0]
	log.Println("Test CancelTask")
	if rsp, _ := cl.CancelTask(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); apiCommon.IsSuccess("CancelTask", rsp) {
		log.Println("CancelTask Ok")
	}

	// Verify Canceled task
	log.Println("Test TaskDetail")
	if rsp, _ := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); apiCommon.IsSuccess("CancelTask Verify", rsp.Error) {
		log.Println("TaskDetail Ok")
		if rsp.Task.Status != common_proto.TaskStatus_CANCEL {
			log.Println(rsp.Task.Status)
			log.Fatalf("CancelTask %s operation does not take effect", cancelTask.Id)
		}
		log.Println("CancelTask takes effect")
	}

	// UpdateTask
	cancelTask.Name = "updateTask"
	if rsp, _ := cl.UpdateTask(tokenContext, &taskmgr.UpdateTaskRequest{UserId: cancelTask.UserId, Task: cancelTask}); apiCommon.IsSuccess("UpdateTask", rsp) {
		log.Println("TaskDetail Ok")
	}

	// Verify updated task
	if rsp, _ := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); apiCommon.IsSuccess("UpdateTask Verify", rsp.Error) {
		if !apiCommon.IsEqual(rsp.Task, cancelTask) || rsp.Task.Status != common_proto.TaskStatus_UPDATING {
			log.Println(rsp.Task)
			log.Println(cancelTask)
			log.Fatal("UpdateTask operation does not take effect")
		}
		log.Println("UpdateTask takes effect")
	}
}
