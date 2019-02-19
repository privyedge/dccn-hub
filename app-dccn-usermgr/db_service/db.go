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
	// Get get a user item by go_micro_srv_user mgr's email
	Get(email string) (*UserRecord, error)
	GetById(id string) (*UserRecord, error)
	GetUserByID(id string) (*UserRecord, error)
	// Create Creates a new user item if not exits
	Create(user *pb.User, password string) error
	// Update updates dc item
	Update(user *pb.User) error
	UpdateUserAttributes(id string, attr [](*usermgr.UserAttribute)) error
	// UpdateStatus update user status in db
	UpdateStatus(email string, status usermgr.UserStatus) error
	UpdateEmail(userId, newEmail string) error
	UpdatePassword(email, newPassword string) error
	UpdateRefreshToken(uid string, token string) error
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

type UserRecord struct {
	ID                 string
	Email              string
	Password           string
	Name               string
	Token              string // refresh token
	varified           bool     // email varified
	Last_modified_date uint64
	Creation_date      uint64
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
func (p *DB) Get(email string) (*UserRecord, error) {
	session := p.session.Clone()
	defer session.Close()

	var user UserRecord
	err := p.collection(session).Find(bson.M{"email": email}).One(&user)
	return &user, err
}

// GetById gets user item by email.
func (p *DB) GetById(id string) (*UserRecord, error) {
	session := p.session.Clone()
	defer session.Close()

	var user UserRecord
	err := p.collection(session).Find(bson.M{"id": id}).One(&user)
	return &user, err
}

func (p *DB) GetUserByID(id string) (*UserRecord, error) {
	session := p.session.Clone()
	defer session.Close()

	var user UserRecord
	err := p.collection(session).Find(bson.M{"id": id}).One(&user)
	return &user, err
}

// Create creates a new user item if it not exists
// TODO: batch operations through bulk
func (p *DB) Create(user *pb.User, password string) error {
	session := p.session.Clone()
	defer session.Close()
	userRecord := UserRecord{}
	userRecord.ID = user.Id
	userRecord.Email = user.Email
	userRecord.Name = user.Attributes.Name
	userRecord.Password = password
	userRecord.Last_modified_date = uint64(time.Now().Unix())
	userRecord.Creation_date = uint64(time.Now().Unix())
	return p.collection(session).Insert(userRecord)
}

// Update updates user item.
func (p *DB) Update(user *pb.User) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Update(bson.M{"email": user.Email}, user)
}

func (p *DB) UpdateStatus(email string, status usermgr.UserStatus) error {
	session := p.session.Clone()
	defer session.Close()
	// TODO: check if ok
	return p.collection(session).Update(bson.M{"email": email}, bson.M{"$set": bson.M{"status": status}})
}

func (p *DB) UpdateEmail(userId, email string) error {
	session := p.session.Clone()
	defer session.Close()
	// TODO: check if ok
	return p.collection(session).Update(bson.M{"id": userId}, bson.M{"$set": bson.M{"email": email}})
}

func (p *DB) UpdateUserAttributes(id string, attr [](*usermgr.UserAttribute)) error {
	session := p.session.Clone()
	defer session.Close()
	//todo
	return p.collection(session).Update(bson.M{"id": id}, bson.M{"$set": bson.M{"attribute": attr}})
}

func (p *DB) UpdatePassword(email, newPassword string) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Update(bson.M{"email": email}, bson.M{"$set": bson.M{"attributes.hashpassword": newPassword}})
}

// Update updates user item.
func (p *DB) UpdateRefreshToken(uid string, token string) error {
	session := p.session.Clone()
	defer session.Close()
	return p.collection(session).Update(bson.M{"id": uid}, bson.M{"$set": bson.M{"token": token}})
}

// Close closes the db connection.
func (p *DB) Close() {
	p.session.Close()
}

func (p *DB) dropCollection() {
	p.session.DB(p.dbName).C(p.collectionName).DropCollection()
}
