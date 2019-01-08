package dbservice

import (
	"sync/atomic"

	go_micro_srv_usermgr "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mysql or mongodb. who is better.
type DBService interface {
	// Get get a user item by go_micro_srv_usermgr's id
	Get(id int64) (*go_micro_srv_usermgr.User, error)
	// Create Creates a new user item if not exits
	Add(user go_micro_srv_usermgr.User) error
	// Update updates dc item
	Update(user *go_micro_srv_usermgr.User) error
	// Close closes db connection
	Close()
}

// DB implements DBService
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
func (p *DB) Get(id int64) (user *go_micro_srv_usermgr.User, err error) {
	session := p.session.Copy()
	defer session.Close()

	var buser bsonUser
	if err = p.collection(session).FindId(id).One(&buser); err != nil {
		return
	}
	err = buser.Decode(user)
	return
}

// Create creates a new user item if it not exists
// TODO: batch operations throught bulk
func (p *DB) Add(user go_micro_srv_usermgr.User) error {
	user.Id.Id = atomic.AddInt64(&p.count, int64(1))
	var buser bsonUser
	if err := buser.Encode(&user); err != nil {
		return err
	}

	session := p.session.Copy()
	defer session.Close()
	return p.collection(session).Insert(&buser)
}

// Update updates user item.
func (p *DB) Update(user *go_micro_srv_usermgr.User) error {
	var buser bsonUser
	if err := buser.Encode(user); err != nil {
		return err
	}

	session := p.session.Copy()
	defer session.Close()
	return p.collection(session).Update(bson.M{"_id": buser.Id}, &buser)
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}
