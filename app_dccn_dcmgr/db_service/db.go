package dbservice

import (
	dcmgr "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/v1"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type DBService interface {
	// Get gets a dc item by pb's id.
	Get(id int64) (*dcmgr.DataCenter, error)
	// Create Creates a new dc item if not exits.
	Create(center *dcmgr.DataCenter) error
	// Update updates dc item
	Update(center *dcmgr.DataCenter) error
	// Close closes db connection
	Close()
	// dropCollection for testing usage
	dropCollection()
}

// UserDB implements DBService
type DB struct {
	dbName         string
	collectionName string
	session        *mgo.Session
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

// Get gets user item by id.
func (p *DB) Get(id int64) (*dcmgr.DataCenter, error) {
	session := p.session.Clone()
	defer session.Close()

	var center dcmgr.DataCenter
	err := p.collection(session).Find(bson.M{"id": id}).One(&center)
	return &center, err
}

// Create creates a new data center item if it not exists
func (p *DB) Create(center *dcmgr.DataCenter) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Insert(center)
}

// Update updates user item.
func (p *DB) Update(datacenter *dcmgr.DataCenter) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Update(bson.M{"id": datacenter.Id}, datacenter)
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}

func (p *DB) dropCollection() {
	p.session.DB(p.dbName).C(p.collectionName).DropCollection()
}
