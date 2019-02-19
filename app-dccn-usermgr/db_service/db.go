package dbservice

import (
	"time"

	dbcommon "github.com/Ankr-network/dccn-common/db"
	pb "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mysql or mongodb. who is better.
type DBService interface {
	// Get gets a user item by user id
	GetUser(id string) (*UserRecord, error)
	// Get get a user item by user email, email unique
	GetUserByEmail(email string) (*UserRecord, error)
	// Create Creates a new user item if not exits
	CreateUser(user *pb.User, hashedPassword string) error
	// Update updates dc item by id
	UpdateUser(id string, fields []*usermgr.UserAttribute) error
	// UpdateByEmail updates dc item by email
	UpdateUserByEmail(email string, fields []*usermgr.UserAttribute) error
	// Close closes db connection
	Close()
	// dropCollection for testing usage
	dropCollection()
}

// DB implements DBService
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

// Get gets user item by email.
func (p *DB) GetUser(id string) (*UserRecord, error) {
	session := p.session.Clone()
	defer session.Close()

	var user UserRecord
	err := p.collection(session).Find(bson.M{"id": id}).One(&user)
	return &user, err
}

// GetUserByEmail gets user item by email.
func (p *DB) GetUserByEmail(email string) (*UserRecord, error) {
	session := p.session.Clone()
	defer session.Close()

	var user UserRecord
	err := p.collection(session).Find(bson.M{"email": email}).One(&user)
	return &user, err
}

// CreateUser creates a new user item if it not exists
// TODO: batch operations through bulk
func (p *DB) CreateUser(user *pb.User, hashedPassword string) error {
	session := p.session.Clone()
	defer session.Close()
	userRecord := UserRecord{}
	userRecord.ID = user.Id
	userRecord.Email = user.Email
	userRecord.Name = user.Attributes.Name
	userRecord.HashedPassword = hashedPassword
	userRecord.LastModifiedDate = uint64(time.Now().Unix())
	userRecord.CreationDate = uint64(time.Now().Unix())
	return p.collection(session).Insert(userRecord)
}

// UpdateUser updates user item.
func (p *DB) UpdateUser(id string, fields []*usermgr.UserAttribute) error {
	session := p.session.Clone()
	defer session.Close()

	return p.collection(session).Update(bson.M{"id": id}, getUpdate(fields))
}

// UpdateUserByEmail updates user item.
func (p *DB) UpdateUserByEmail(email string, fields []*usermgr.UserAttribute) error {
	session := p.session.Clone()
	defer session.Close()

	return p.collection(session).Update(bson.M{"email": email}, getUpdate(fields))
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}

func (p *DB) dropCollection() {
	p.session.DB(p.dbName).C(p.collectionName).DropCollection()
}
