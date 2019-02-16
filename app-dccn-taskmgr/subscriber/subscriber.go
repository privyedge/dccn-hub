package subscriber

import (
	"context"
	"log"

	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	db "github.com/Ankr-network/dccn-hub/app-dccn-taskmgr/db_service"
	"gopkg.in/mgo.v2/bson"
)

type TaskStatusFeedback struct {
	db db.DBService
}

func New(db db.DBService) *TaskStatusFeedback {
	return &TaskStatusFeedback{db}
}

// UpdateTaskByFeedback receives task result from data center, returns to v1
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *TaskStatusFeedback) HandlerFeedbackEventFromDataCenter(ctx context.Context, event *common_proto.Event) error {

	feedback := event.GetTaskFeedback()
	log.Printf("HandlerFeedbackEventFromDataCenter: Receive New Event: %+v", *feedback)
	var update bson.M
	switch event.EventType {
	case common_proto.Operation_TASK_CREATE:  // feedback  TaskStatus_START_FAILED  TaskStatus_START_SUCCESS => TaskStatus_RUNNING
		status := common_proto.TaskStatus_RUNNING
		if feedback.Status == common_proto.TaskStatus_START_FAILED {
			status = common_proto.TaskStatus_START_FAILED
		}
		update = bson.M{"$set": bson.M{"status": status, "datacenter": feedback.DataCenter}}
		if event.GetTaskFeedback().Url != "" {
			update = bson.M{"$set": bson.M{"status": status, "datacenter": feedback.DataCenter, "url": feedback.Url, "report": feedback.Report}}
		}
	case common_proto.Operation_TASK_UPDATE:
		status := common_proto.TaskStatus_RUNNING
		update = bson.M{"$set": bson.M{"status": status}}
	case common_proto.Operation_TASK_CANCEL:
		status := common_proto.TaskStatus_CANCELLED
		if feedback.Status == common_proto.TaskStatus_CANCEL_FAILED {
			status = common_proto.TaskStatus_CANCEL_FAILED
		}


		update = bson.M{"$set": bson.M{"status": status, "report": feedback.Report}}
	}

	return p.db.Update(feedback.TaskId, update)
}