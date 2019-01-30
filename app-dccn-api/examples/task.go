package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/grpc"
	apiCommon "github.com/Ankr-network/dccn-hub/app-dccn-api/examples/common"
)

var addr = "client-dev.dccn.ankr.network:50051"

// var addr = "localhost:50051"
var token = "test token here"

// func init() {
// 	addr = os.Getenv("API_ADDRESS")
// 	log.Println("API ADDRESS: ", addr)
// }

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

	// request with token verify information
	md := metadata.New(map[string]string{
		"token": token,
	})

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	tokenContext, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tasks := apiCommon.MockTasks()
	log.Println("Test CreateTask")
	for i := range tasks {
		if rsp, err := cl.CreateTask(tokenContext, &taskmgr.CreateTaskRequest{UserId: tasks[i].UserId, Task: &tasks[i]}); err != nil {
			log.Fatal(err.Error())
		} else {
			log.Println(*rsp)
		}
	}

	log.Println("Test TaskList")
	var userTasks []*common_proto.Task
	if rsp, err := cl.TaskList(tokenContext, &taskmgr.ID{UserId: "1"}); err != nil {
		log.Fatal(err.Error())
	} else {
		userTasks = append(userTasks, rsp.Tasks...)
		log.Println("TaskList Ok")
		if len(userTasks) == 0 {
			log.Fatalf("no tasks belongs to %d\n", 1)
		}
	}

	// CancelTask
	cancelTask := userTasks[0]
	log.Println("Test CancelTask")
	if _, err := cl.CancelTask(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("CancelTask Ok")

		// Verify Canceled task
		log.Println("Test TaskDetail")
		if _, err := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); err != nil {
			log.Fatal(err.Error())
		} else {
			log.Println("TaskDetail Ok")
		}
	}

	// UpdateTask
	cancelTask.Name = "updateTask"
	if _, err := cl.UpdateTask(tokenContext, &taskmgr.UpdateTaskRequest{UserId: cancelTask.UserId, Task: cancelTask}); err != nil {
		log.Fatal(err.Error())
	} else {
		log.Println("TaskDetail Ok")
		// Verify updated task
		if rsp, err := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); err != nil {
			log.Fatal(err.Error())
		} else {
			if !apiCommon.IsEqual(rsp.Task, cancelTask) || rsp.Task.Status != common_proto.TaskStatus_UPDATING {
				log.Fatal("UpdateTask operation does not take effect")
			}
			log.Println("UpdateTask takes effect")
		}
	}
}
