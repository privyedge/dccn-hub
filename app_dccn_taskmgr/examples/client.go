package main

import (
	"context"
	"log"
	"reflect"

	grpc "github.com/micro/go-grpc"
	"github.com/micro/go-micro/metadata"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
)

func isEqual(origin, dst *common_proto.Task) bool {
	ok := origin.Id == dst.Id &&
		origin.UserId == dst.UserId &&
		origin.Type == dst.Type &&
		origin.Name == dst.Name &&
		origin.Image == dst.Image &&
		origin.Replica == dst.Replica &&
		origin.DataCenter == dst.DataCenter &&
		origin.DataCenterId == dst.DataCenterId &&
		origin.Status == dst.Status &&
		origin.UniqueName == dst.UniqueName &&
		origin.Url == dst.Url &&
		origin.Hidden == dst.Hidden &&
		origin.Uptime == dst.Uptime &&
		origin.CreationDate == dst.CreationDate
	if origin.Extra != nil && dst.Extra != nil {
		ok = ok && reflect.DeepEqual(origin.Extra, dst.Extra)
	}
	return ok
}

func mockTasks() []common_proto.Task {
	return []common_proto.Task{
		common_proto.Task{
			Id:           "001",
			UserId:       1,
			Name:         "task01",
			Type:         "web",
			Image:        "nginx",
			Replica:      2,
			DataCenter:   "dc01",
			DataCenterId: 1,
		},
		common_proto.Task{
			Id:           "002",
			UserId:       1,
			Name:         "task02",
			Type:         "web",
			Image:        "nginx",
			Replica:      2,
			DataCenter:   "dc02",
			DataCenterId: 1,
		},
		common_proto.Task{
			Id:           "003",
			UserId:       2,
			Name:         "task01",
			Type:         "web",
			Image:        "nginx",
			Replica:      2,
			DataCenter:   "dc01",
			DataCenterId: 1,
		},
	}
}

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

	tasks := mockTasks()
	for i := range tasks {
		if rsp, _ := cl.CreateTask(tokenContext, &taskmgr.AddTaskRequest{UserId: tasks[i].UserId, Task: &tasks[i]}); rsp.Error != nil {
			log.Fatalf(rsp.Error.Details)
		} else {
			log.Println("CreateTask Ok")
		}
	}

	if rsp, _ := cl.TaskList(tokenContext, &taskmgr.ID{UserId: 2}); rsp.Error != nil {
		log.Fatalf("TaskList Error: ", rsp.Error.Details)
	} else {
		log.Printf("TaskList Ok: %#v\n", rsp.Tasks)
	}

	cancelTask := tasks[0]
	if rsp, _ := cl.CancelTask(tokenContext, &taskmgr.Request{UserId: cancelTask.UserId, TaskId: cancelTask.Id}); rsp != nil {
		log.Fatalf("CancelTask Error: ", rsp.Details)
	} else {
		log.Println("CancelTask Ok")
	}

	if rsp, _ := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: 1, TaskId: "001"}); rsp.Error != nil {
		log.Fatalf("TaskDetail Error: ", rsp.Error.Details)
	} else {
		if rsp.Task.Status != common_proto.TaskStatus_CANCELLED {
			log.Fatal("CancelTask not changed")
		} else {
			log.Printf("TaskDetail Ok: %#v\n", rsp.Task)
		}
	}

	task := tasks[0]
	task.Name = "updateTask"
	if rsp, _ := cl.UpdateTask(tokenContext, &taskmgr.UpdateTaskRequest{UserId: task.UserId, Task: &task}); rsp != nil {
		log.Fatalf("TaskDetail Error: ", rsp.Details)
	} else {
		log.Println("TaskDetail Ok")
	}

	if rsp, _ := cl.TaskDetail(tokenContext, &taskmgr.Request{UserId: 1, TaskId: "001"}); rsp.Error != nil {
		log.Fatalf("TaskDetail Error: ", rsp.Error.Details)
	} else {
		if !isEqual(rsp.Task, &task) {
			log.Fatal("Update Not Changed")
		} else {
			log.Printf("TaskDetail Ok: %#v\n", rsp.Task)
		}
	}
}
