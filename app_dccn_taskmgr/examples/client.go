package main

import (
	"context"
	"log"

	grpc "github.com/micro/go-grpc"
	"github.com/micro/go-micro/metadata"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
	testCommon "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/examples/common"
)

var token = "token"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("client service start...")
	srv := grpc.NewService()

	srv.Init()

	cl := taskmgr.NewTaskMgrService(ankr_default.TaskMgrRegistryServerName, srv.Client())

	tokenContext := metadata.NewContext(context.Background(), map[string]string{
		"Token": token,
	})

	tasks := testCommon.MockTasks()
	for i := range tasks {
		if rsp, _ := cl.CreateTask(tokenContext, &taskmgr.CreateTaskRequest{UserId: tasks[i].UserId, Task: &tasks[i]}); testCommon.IsSuccess("CreateTask", rsp.Error) {
			log.Println("CreateTask Ok")
		}
	}

	userTasks := []*common_proto.Task{}
	if rsp, _ := cl.TaskList(tokenContext, &taskmgr.ID{UserId: 1}); testCommon.IsSuccess("TaskList", rsp.Error) {
		userTasks = append(userTasks, rsp.Tasks...)
		log.Println("TaskList Ok")
	}

	if len(userTasks) == 0 {
		log.Fatalf("no tasks belongs to %d\n", 1)
	}

	// CancelTask
	cancelTask := userTasks[0]
	if rsp, _ := cl.CancelTask(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); testCommon.IsSuccess("CancelTask", rsp) {
		log.Println("CancelTask Ok")
	}

	// Verify Canceled task
	if rsp, _ := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); testCommon.IsSuccess("CancelTask Verify", rsp.Error) {
		log.Println("TaskDetail Ok")
		if rsp.Task.Status != common_proto.TaskStatus_CANCELL {
			log.Println(rsp.Task.Status)
			log.Fatalf("CancelTask %s operation does not take effect", cancelTask.Id)
		}
		log.Println("CancelTask takes effect")
	}

	// UpdateTask
	cancelTask.Name = "updateTask"
	if rsp, _ := cl.UpdateTask(tokenContext, &taskmgr.UpdateTaskRequest{UserId: cancelTask.UserId, Task: cancelTask}); testCommon.IsSuccess("UpdateTask", rsp) {
		log.Println("TaskDetail Ok")
	}

	// Verify updated task
	if rsp, _ := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); testCommon.IsSuccess("UpdateTask Verify", rsp.Error) {
		if !testCommon.IsEqual(rsp.Task, cancelTask) || rsp.Task.Status != common_proto.TaskStatus_UPDATING {
			log.Println(rsp.Task)
			log.Println(cancelTask)
			log.Fatal("UpdateTask operation does not take effect")
		}
		log.Println("UpdateTask takes effect")
	}
}
