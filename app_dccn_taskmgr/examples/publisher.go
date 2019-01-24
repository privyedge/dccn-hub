package main

import (
	"context"
	"log"

	grpc "github.com/micro/go-grpc"
	micro "github.com/micro/go-micro"

	_ "github.com/micro/go-plugins/broker/rabbitmq"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1/micro"
	testCommon "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/examples/common"
)

// send events using the publisher
func sendEv(taskId string, p micro.Publisher) {

	// create new event
	ev := common_proto.Event{
		EventType: common_proto.Operation_TASK_CANCEL,
		OpMessage: &common_proto.Event_TaskFeedback{TaskFeedback: &common_proto.TaskFeedback{
			TaskId: taskId,
			Status: common_proto.TaskStatus_CANCELL_FAILED,
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

	// pub to topic 1
	sendEv(task.Id, pub)

	// Verify publish event
	if rsp, _ := cl.TaskDetail(context.TODO(), &taskmgr.Request{UserId: task.UserId, TaskId: task.Id}); testCommon.IsSuccess("UpdateTask Verify", rsp.Error) {
		if rsp.Task.Status != common_proto.TaskStatus_CANCELL_FAILED {
			log.Fatal("Publish Failed")
		} else {
			log.Println("TaskDetail Ok")
		}
	}

	log.Println("Pub End")

}
