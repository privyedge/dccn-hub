package dbservice

import (
	"sync/atomic"

	pb "github.com/Ankr-network/dccn-hub/app_dccn_taskmgr/proto/taskmgr"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBService interface {
	// Get gets a task item by pb's id.
	Get(id int64) (*pb.Task, error)
	// GetByEventId gets task by event id.
	GetByEventId(eventId string) (task *pb.Task, err error)
	// ValidateAndUpdateStatus validates by event id, if ok, update task status.
	ValidateAndUpdateStatus(event *pb.Event, status pb.Task_Status) error
	// ValidateAndUpdateResult validates by event id, if ok, update task result.
	ValidateAndUpdateResult(event *pb.Event) error
	// Cancel sets task status CANCEL
	Cancel(id int64) error
	// Create Creates a new dc item if not exits.
	Create(task pb.Task) error
	// Update updates dc item
	Update(task *pb.Task) error
	// Close closes db connection
	Close()
}

// UserDB implements DBService
type DB struct {
	dbName              string
	collectionName      string
	eventCollectionName string
	count               int64
	session             *mgo.Session
}

// New returns DBService.
func New(conf dbcommon.Config) (*DB, error) {
	session, err := dbcommon.CreateDBConnection(conf)
	if err != nil {
		return nil, err
	}

	count, err := session.DB(conf.DB).C(conf.Collection).Count()
	if err != nil {
		return nil, err
	}

	return &DB{
		dbName:         conf.DB,
		collectionName: conf.Collection,
		count:          int64(count),
		session:        session,
	}, nil
}

func (p *DB) collection(session *mgo.Session) *mgo.Collection {
	return session.DB(p.dbName).C(p.collectionName)
}

// Get gets user item by id.
func (p *DB) Get(id int64) (task *pb.Task, err error) {
	session := p.session.Copy()
	defer session.Close()

	var btask bsonTask
	if err = p.collection(session).FindId(id).One(&btask); err != nil {
		return
	}
	err = btask.Decode(task)
	return
}

// GetByEventId gets task by event id.
func (p *DB) GetByEventId(eventId string) (task *pb.Task, err error) {
	session := p.session.Copy()
	defer session.Close()

	var btask bsonTask
	if err = p.collection(session).Find(bson.M{"event_id": eventId}).One(&btask); err != nil {
		return
	}
	err = btask.Decode(task)
	return
}

// Create creates a new task item if it not exists
func (p *DB) Create(task pb.Task) error {
	task.Status = pb.Task_CANCELING

	task.Id = atomic.AddInt64(&p.count, int64(1))
	var btask bsonTask
	if err := btask.Encode(&task); err != nil {
		return err
	}

	session := p.session.Copy()
	defer session.Close()
	return p.collection(session).Insert(&btask)
}

// Update updates task item.
func (p *DB) Update(task *pb.Task) error {
	task.Status = pb.Task_UPDATING

	var btask bsonTask
	if err := btask.Encode(task); err != nil {
		return err
	}

	session := p.session.Copy()
	defer session.Close()
	return p.collection(session).Update(bson.M{"_id": btask.Id}, &btask)
}

// ValidateAndUpdateStatus validates by event id, if ok, update task status.
func (p *DB) ValidateAndUpdateStatus(event *pb.Event, status pb.Task_Status) error {
	task, err := p.GetByEventId(event.Id)
	if err != nil {
		return err
	}

	session := p.session.Copy()
	defer session.Close()
	return p.collection(session).Update(bson.M{"_id": task.Id}, bson.M{"$set": bson.M{"status": uint32(status)}})
}

// ValidateAndUpdateResult validates by event id, if ok, update task result.
func (p *DB) ValidateAndUpdateResult(event *pb.Event) error {
	task, err := p.GetByEventId(event.Id)
	if err != nil {
		return err
	}

	session := p.session.Copy()
	defer session.Close()
	return p.collection(session).Update(bson.M{"_id": task.Id}, bson.M{"$set": bson.M{"result": event.Result}})
}

// Cancel cancel task, sets task status CANCEL
func (p *DB) Cancel(id int64) error {
	session := p.session.Copy()
	defer session.Close()
	return p.collection(session).Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"status": pb.Task_CANCELING}})
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}
