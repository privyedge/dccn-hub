package dbservice

import (
	"sync/atomic"

	go_micro_srv_dcmgr "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/dcmgr"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBService interface {
	// Get gets a dc item by pb's id.
	Get(id int64) (*go_micro_srv_dcmgr.DataCenter, error)
	// Create Creates a new dc item if not exits.
	Add(center go_micro_srv_dcmgr.DataCenter) error
	// Update updates dc item
	Update(center *go_micro_srv_dcmgr.DataCenter) error
	// Close closes db connection
	Close()
}

// UserDB implements DBService
type DB struct {
	dbName         string
	collectionName string
	count          int64
	session        *mgo.Session
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
func (p *DB) Get(id int64) (center *go_micro_srv_dcmgr.DataCenter, err error) {
	session := p.session.Clone()
	defer session.Close()

	var bcenter bsonDataCenter
	if err = p.collection(session).FindId(id).One(&bcenter); err != nil {
		return
	}
	err = bcenter.Decode(center)
	return
}

// Create creates a new data center item if it not exists
func (p *DB) Add(center go_micro_srv_dcmgr.DataCenter) error {
	center.Id = atomic.AddInt64(&p.count, int64(1))
	var bcenter bsonDataCenter
	if err := bcenter.Encode(&center); err != nil {
		return err
	}

	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Insert(&bcenter)
}

// Update updates datacenter item.
func (p *DB) Update(center *go_micro_srv_dcmgr.DataCenter) error {
	var bcenter bsonDataCenter
	if err := bcenter.Encode(center); err != nil {
		return err
	}

	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Update(bson.M{"_id": bcenter.Id}, &bcenter)
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}
