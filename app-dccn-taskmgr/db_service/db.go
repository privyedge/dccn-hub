package dbservice

import (
	"log"

	dbcommon "github.com/Ankr-network/dccn-common/db"
	common_proto "github.com/Ankr-network/dccn-common/protos/common"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBService interface {
	// Get gets a task item by taskmgr's id.
	Get(id string) (*common_proto.Task, error)
	// GetAll gets all task related to user id.
	GetAll(userId string) (*[]*common_proto.Task, error)
	// GetByEventId gets task by event id.
	GetByEventId(eventId string) (task *[]*common_proto.Task, err error)
	// Cancel sets task status CANCEL
	Cancel(taskId string) error
	// Create Creates a new dc item if not exits.
	Create(task *common_proto.Task) error
	// Update updates dc item
	Update(taskId string, update bson.M) error
	// UpdateTask updates dc item
	UpdateTask(taskId string, task *common_proto.Task) error
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
func (p *DB) Get(taskId string) (*common_proto.Task, error) {
	session := p.session.Clone()
	defer session.Close()

	var task common_proto.Task
	err := p.collection(session).Find(bson.M{"id": taskId}).One(&task)
	return &task, err
}

func (p *DB) GetAll(userId string) (*[]*common_proto.Task, error) {
	session := p.session.Clone()
	defer session.Close()

	var tasks []*common_proto.Task

	log.Printf("find tasks with uid %s", userId)

	if err := p.collection(session).Find(bson.M{"userid": userId}).All(&tasks); err != nil {
		return nil, err
	}
	return &tasks, nil
}

// GetByEventId gets task by event id.
func (p *DB) GetByEventId(eventId string) (*[]*common_proto.Task, error) {
	session := p.session.Copy()
	defer session.Close()

	var tasks []*common_proto.Task
	if err := p.collection(session).Find(bson.M{"eventid": eventId}).One(&tasks); err != nil {
		return nil, err
	}
	return &tasks, nil
}

// Create creates a new task item if it not exists
func (p *DB) Create(task *common_proto.Task) error {
	session := p.session.Copy()
	defer session.Close()

	log.Printf("create task %+v\n", task)
	return p.collection(session).Insert(task)
}

// Update updates task item.
func (p *DB) Update(taskId string, update bson.M) error {
	session := p.session.Copy()
	defer session.Close()

	return p.collection(session).Update(bson.M{"id": taskId}, update)
}

func (p *DB) UpdateTask(taskId string, task *common_proto.Task) error {
	session := p.session.Copy()
	defer session.Close()

	fields := bson.M{}

	if len(task.Name) > 0 {
		fields["name"] = task.Name
	}

	if task.Replica > 0 {
		fields["replica"] = task.Replica
	}

	if task.Status > 0 {
		fields["status"] = task.Status
	}

	if task.Hidden {
		fields["hidden"] = task.Hidden
	}

	if len(task.Image) > 0 {
		fields["image"] = task.Image
	}

	if len(task.DataCenter) > 0 {
		fields["datacenter"] = task.DataCenter
	}


	if len(task.Url) > 0 {
		fields["url"] = task.Url
	}

	return p.collection(session).Update(bson.M{"id": taskId}, bson.M{"$set": fields})


	//return p.collection(session).Update(bson.M{"id": taskId}, task)
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
