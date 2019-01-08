package dbservice

import (
	"errors"

	taskmgr "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/proto/taskmgr"
)

type bsonTask struct {
	// id task id, unique
	Id int64 `json:"id,omitempty" bson:"_id"`
	// event_id the task belongs.
	EventId string `json:"event_id,omitempty" bson:"event_id,omitempty"`
	// user_id the task belongs.
	UserId int64 `json:"user_id,omitempty" bson:"user_id,omitempty"`
	// name task name
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// task type
	Type string `json:"type,omitempty" bson:"type,omitempty"`
	// status [running,pending,cancel,done]
	Status uint32 `json:"status,omitempty" bson:"status,omitempty"`
	// startup time of the this task
	StartupTime uint32 `json:"startup_time,omitempty" bson:"startup_time,omitempty"`
	// mission data center id
	DataCenterId int64 `json:"data_center_id,omitempty" bson:"data_center_id,omitempty"`
	// task creation time.
	CreateTime uint64 `json:"creation_date,omitempty" bson:"create_time,omitempty"`
	// task cancel time
	CancelTime uint64 `json:"cancel_date,omitempty" bson:"cancel_time,omitempty"`
	// task updates time.
	UpdateTime uint64 `json:"cancel_date,omitempty" bson:"update_time,omitempty"`
	// task returns time.
	ReturnTime uint64 `json:"cancel_date,omitempty" bson:"return_time,omitempty"`
	// extra for other arguments
	Extra []byte `json:"extra,omitempty" bson:"extra,omitempty"`
	// result of the task
	Result []byte `json:"result,omitempty" bson:"result,omitempty"`
}

func (btask *bsonTask) Encode(task *taskmgr.Task) error {
	if btask == nil || task == nil {
		return errors.New("Null Pointer")
	}
	btask.Id = task.Id
	btask.EventId = task.EventId
	btask.UserId = task.UserId
	btask.Name = task.Name
	btask.Type = task.Type
	btask.Status = uint32(task.Status)
	btask.StartupTime = task.StartupTime
	btask.DataCenterId = task.DataCenterId
	btask.CreateTime = task.CreateTime
	btask.CancelTime = task.CancelTime
	btask.UpdateTime = task.UpdateTime
	btask.ReturnTime = task.ReturnTime
	btask.Extra = task.Extra
	btask.Result = task.Result

	return nil
}

func (btask *bsonTask) Decode(task *taskmgr.Task) error {
	if task == nil || btask == nil {
		return errors.New("Null Pointer")
	}
	task.Id = btask.Id
	task.EventId = btask.EventId
	task.UserId = btask.UserId
	task.Name = btask.Name
	task.Type = btask.Type
	task.Status = taskmgr.Task_Status(btask.Status)
	task.StartupTime = btask.StartupTime
	task.DataCenterId = btask.DataCenterId
	task.CreateTime = btask.CreateTime
	task.CancelTime = btask.CancelTime
	task.UpdateTime = btask.UpdateTime
	task.ReturnTime = btask.ReturnTime
	task.Extra = btask.Extra
	task.Result = btask.Result

	return nil
}
