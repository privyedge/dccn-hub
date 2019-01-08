package subscriber

import (
	"context"
	"errors"

	db "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/db_service"
	taskmgr "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/proto/taskmgr"
)

type Subscriber struct {
	db db.DBService
}

func New(db db.DBService) *Subscriber {
	return &Subscriber{db}
}

// UpdateTaskByFeedback receives task result from data center, returns to taskmgr
// UpdateTaskStatusByFeedback updates database status by performing feedback from the data center of the task.
// sets executor's id, updates task status.
func (p *Subscriber) UpdateTaskByFeedback(ctx context.Context, event *taskmgr.Event) error {
	switch event.Operation {
	case taskmgr.Event_CREATE:
		return p.db.ValidateAndUpdateStatus(event, taskmgr.Task_CREATED)
	case taskmgr.Event_CANCEL:
		return p.db.ValidateAndUpdateStatus(event, taskmgr.Task_CANCELED)
	case taskmgr.Event_UPDATE:
		return p.db.ValidateAndUpdateStatus(event, taskmgr.Task_UPDATED)
	case taskmgr.Event_RETURN:
		return p.db.ValidateAndUpdateResult(event)
	default:
		return errors.New("Unknown Operation Code")
	}
	return nil
}
