package main

import (
	"context"
	"log"
	"time"

	grpc "github.com/micro/go-grpc"
	micro "github.com/micro/go-micro"

	_ "github.com/micro/go-plugins/broker/rabbitmq"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
	testCommon "github.com/Ankr-network/dccn-hub/app-dccn-taskmgr/examples/common"
)

// send events using the publisher
func sendEv(taskId string, p micro.Publisher) {

	// create new event
	ev := common_proto.Event{
		EventType: common_proto.Operation_TASK_CANCEL,
		OpMessage: &common_proto.Event_TaskFeedback{TaskFeedback: &common_proto.TaskFeedback{
			TaskId: taskId,
			Status: common_proto.TaskStatus_CANCEL_FAILED,
		}},
	}

	log.Printf("publishing %+v\n", ev)

	// publish an event
	if err := p.Publish(context.Background(), &ev); err != nil {
		log.Fatalf("error publishing %v\n", err)
	}
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// create a service
	service := grpc.NewService()

	// parse command line
	service.Init()

	// create publisher
	pub := micro.NewPublisher(ankr_default.MQFeedbackTask, service.Client())

	cl := taskmgr.NewTaskMgrService(ankr_default.TaskMgrRegistryServerName, service.Client())
	task := testCommon.MockTasks()[0]
	if rsp, _ := cl.CreateTask(context.TODO(), &taskmgr.CreateTaskRequest{UserId: task.UserId, Task: &task}); testCommon.IsSuccess("CreateTask", rsp.Error) {
		log.Println("CreateTask Ok")
	}

	var userTasks []*common_proto.Task
	if rsp, _ := cl.TaskList(context.TODO(), &taskmgr.ID{UserId: "1"}); testCommon.IsSuccess("TaskList", rsp.Error) {
		userTasks = append(userTasks, rsp.Tasks...)
		log.Println("TaskList Ok")
	}

	if len(userTasks) == 0 {
		log.Fatalf("no tasks belongs to %d\n", 1)
	}

	pubTask := userTasks[0]

	// pub to topic 1
	sendEv(pubTask.Id, pub)

	// waits pub message arrive to mq
	time.Sleep(2 * time.Second)

	// Verify publish event
	if rsp, _ := cl.TaskDetail(context.TODO(), &taskmgr.Request{UserId: pubTask.UserId, TaskId: pubTask.Id}); testCommon.IsSuccess("UpdateTask Verify", rsp.Error) {
		if rsp.Task.Status != common_proto.TaskStatus_CANCEL_FAILED {
			log.Fatal("UpdateTaskByFeedback do not task effect")
		} else {
			log.Println("TaskDetail Ok")
		}
	}

	log.Println("Pub End")
}
