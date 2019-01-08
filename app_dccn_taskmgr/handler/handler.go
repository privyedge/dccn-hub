package handler

import (
	"context"
	"time"

	db "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/db_service"
	pb "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/proto/taskmgr"
	micro "github.com/micro/go-micro"
	"github.com/pborman/uuid"
)

type TaskMgrHandler struct {
	db        db.DBService
	newPub    micro.Publisher
	cancelPub micro.Publisher
}

func (p *TaskMgrHandler) Get(ctx context.Context, id *pb.ID, task *pb.Task) error {
	var err error
	task, err = p.db.Get(id.Id)
	return err
}

func (p *TaskMgrHandler) Create(ctx context.Context, task *pb.Task, rsp *pb.Response) error {
	event := pb.Event{
		Id:        uuid.NewUUID().String(),
		Timestamp: time.Now().Unix(),
		Operation: pb.Event_CREATE,
	}
	// TODO:
	if err := p.newPub.Publish(context.Background(), event); err != nil {
		return err
	}
	task.EventId = event.Id
	return p.db.Create(*task)
}

func (p *TaskMgrHandler) Cancel(ctx context.Context, id *pb.ID, rsp *pb.Response) error {
	task, err := p.db.Get(id.Id)
	if err != nil {
		return nil
	}
	event := pb.Event{
		Id:           task.EventId,
		DataCenterId: task.DataCenterId,
		TaskId:       id.Id,
		Timestamp:    time.Now().Unix(),
		Operation:    pb.Event_CANCEL,
	}
	// TODO: fanout module, channel type 'fanout'
	if err := p.newPub.Publish(context.Background(), event); err != nil {
		return err
	}
	return p.db.Cancel(id.Id)
}

func (p *TaskMgrHandler) Update(ctx context.Context, task *pb.Task, rsp *pb.Response) error {
	task, err := p.db.Get(task.Id)
	if err != nil {
		return nil
	}
	event := pb.Event{
		Id:           task.EventId,
		DataCenterId: task.DataCenterId,
		Timestamp:    time.Now().Unix(),
		Operation:    pb.Event_UPDATE,
	}
	if err := p.newPub.Publish(context.Background(), event); err != nil {
		return err
	}
	return p.db.Update(task)
}

func New(db db.DBService, newPub micro.Publisher, cancelPub micro.Publisher) *TaskMgrHandler {
	return &TaskMgrHandler{
		db:        db,
		newPub:    newPub,
		cancelPub: cancelPub,
	}
}
