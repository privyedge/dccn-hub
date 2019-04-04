package dbservice

import (
	"time"

	dbcommon "github.com/Ankr-network/dccn-common/db"
	pb "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"
	micro "github.com/Ankr-network/dccn-usermgr/app-dccn-usermgr/micro"
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

}

// DB implements DBService
type DB struct {
	collection        *mgo.Collection
}

// New returns DBService.
func New(conf dbcommon.Config) (*DB, error) {
	collection := micro.GetCollection("user")
	return &DB{
		collection: collection,
	}, nil
}

// CreateUser creates a new user item if it not exists
// TODO: batch operations through bulk
func (p *DB) CreateUser(user *pb.User, hashedPassword string) error {

	err := p.collection.Insert(&UserRecord{
		ID:               user.Id,
		Email:            user.Email,
		Name:             user.Attributes.Name,
		Status:           user.Status,
		HashedPassword:   hashedPassword,
		LastModifiedDate: uint64(time.Now().Unix()),
		CreationDate:     uint64(time.Now().Unix()),
	})
	return err
}

// Get gets user item by email.
func (p *DB) GetUser(id string) (*UserRecord, error) {
	var user UserRecord
	err := p.collection.Find(bson.M{"id": id}).One(&user)
	return &user, err
}

// GetUserByEmail gets user item by email.
func (p *DB) GetUserByEmail(email string) (*UserRecord, error) {
	var user UserRecord
	err := p.collection.Find(bson.M{"email": email}).One(&user)
	return &user, err
}

// UpdateUser updates user item.
func (p *DB) UpdateUser(id string, fields []*usermgr.UserAttribute) error {
	update, err := getUpdate(fields)
	if err != nil {
		return err
	}

	return p.collection.Update(bson.M{"id": id}, update)
}

// UpdateUserByEmail updates user item.
// UpdateUserByEmail updates user item.
func (p *DB) UpdateUserByEmail(email string, fields []*usermgr.UserAttribute) error {
	update, err := getUpdate(fields)
	if err != nil {
		return err
	}
	return p.collection.Update(bson.M{"email": email}, update)
}

// Close closes the db connection.
func (p *DB) Close() {
	p.Close()
}

