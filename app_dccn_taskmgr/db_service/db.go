package dbservice

import (
	"log"

	dbcommon "github.com/Ankr-network/dccn-common/db"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"
	taskmgr "github.com/Ankr-network/dccn-common/protos/taskmgr/v1"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBService interface {
	// Get gets a task item by taskmgr's id.
	Get(id string) (*taskmgr.Task, error)
	// GetAll gets all task related to user id.
	GetAll(userId int64) (*[]*taskmgr.Task, error)
	// GetByEventId gets task by event id.
	GetByEventId(eventId string) (task *[]*taskmgr.Task, err error)
	// Cancel sets task status CANCEL
	Cancel(taskId string) error
	// Create Creates a new dc item if not exits.
	Create(task *taskmgr.Task) error
	// Update updates dc item
	Update(taskId string, update bson.M) error
	// Updatetask updates dc item
	UpdateTask(taskId string, task *taskmgr.Task) error
	// Close closes db connection
	Close()
	// for test usage
	dropCollection()
}

// UserDB implements DBService
type DB struct {
	dbName              string
	collectionName      string
	eventCollectionName string
	session             *mgo.Session
}

// New returns DBService.
func New(conf dbcommon.Config) (*DB, error) {
	session, err := dbcommon.CreateDBConnection(conf)
	if err != nil {
		return nil, err
	}

	return &DB{
		dbName:         conf.DB,
		collectionName: conf.Collection,
		session:        session,
	}, nil
}

func (p *DB) collection(session *mgo.Session) *mgo.Collection {
	return session.DB(p.dbName).C(p.collectionName)
}

// Get gets task item by id.
func (p *DB) Get(taskId string) (task *taskmgr.Task, err error) {
	session := p.session.Clone()
	defer session.Close()

	err = p.collection(session).Find(bson.M{"id": taskId}).One(task)
	return
}

func (p *DB) GetAll(userId int64) (*[]*taskmgr.Task, error) {
	session := p.session.Clone()
	defer session.Close()

	var tasks []*taskmgr.Task

	if err := p.collection(session).Find(bson.M{"userid": userId}).All(&tasks); err != nil {
		return nil, err
	}
	return &tasks, nil
}

// GetByEventId gets task by event id.
func (p *DB) GetByEventId(eventId string) (*[]*taskmgr.Task, error) {
	session := p.session.Copy()
	defer session.Close()

	var tasks []*taskmgr.Task
	if err := p.collection(session).Find(bson.M{"eventid": eventId}).One(&tasks); err != nil {
		return nil, err
	}
	return &tasks, nil
}

// Create creates a new task item if it not exists
func (p *DB) Create(task *taskmgr.Task) error {
	session := p.session.Copy()
	defer session.Close()

	return p.collection(session).Insert(task)
}

// Update updates task item.
func (p *DB) Update(taskId string, update bson.M) error {
	session := p.session.Copy()
	defer session.Close()

	return p.collection(session).Update(bson.M{"taskid": taskId}, update)
}

func (p *DB) UpdateTask(taskId string, task *taskmgr.Task) error {
	session := p.session.Copy()
	defer session.Close()

	return p.collection(session).Update(bson.M{"taskid": taskId}, task)
}

// Cancel cancel task, sets task status CANCEL
func (p *DB) Cancel(taskId string) error {
	session := p.session.Copy()
	defer session.Close()

	return p.collection(session).Update(bson.M{"id": taskId}, bson.M{"$set": bson.M{"status": common_proto.TaskStatus_CANCELLED}})
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}

func (p *DB) dropCollection() {
	err := p.session.DB(p.dbName).C(p.collectionName).DropCollection()
	if err != nil {
		log.Println(err.Error())
	}
}
