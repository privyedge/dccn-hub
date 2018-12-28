package dbuser

import (
	pb "github.com/Ankr-network/refactor/app_dccn_usermgr/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"sync/atomic"
)

const (
	userCollection = "user"
)

// 101:  User already exists"
// 102: Other errors

// mysql or mongodb. who is better.
type UserDBService interface {
	// New Creates a new user if not exits
	New(user *pb.User) *pb.Error
	// Get get user by user's id
	Get(name string) (*pb.User, *pb.Error)
}

// UserDB implements UserDBService
type UserDB struct{
	mu sync.Mutex
	count int64
	c *mgo.Collection
}

// New Creates a new user if not exits.
func New(c *mgo.Collection) (*UserDB, error) {
	count, err := c.Find(bson.M{}).Count()
	if err != nil {
		return nil, err
	}
	return &UserDB{count: int64(count), c: c}, nil
}

// Get gets user by user's id.
func (p *UserDB) Get(name string) (*pb.User, error) {
	var result *pb.User
	if err := p.c.Find(bson.M{"name": name}).One(result); err != nil {
		return nil, err
	}
	return result, nil
}

// New creates a new user if it not exists, name unique as _id
// TODO: batch operations throught bulk
func (p *UserDB) New(user pb.User) *pb.Error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// the UI ensures the name not empty
	ac, err := p.Get(user.Name)
	if err != nil {
		return &pb.Error{Code:102, Description: err.Error()}
	}

	if ac != nil {
		return &pb.Error{Code: 101, Description: "User already exists"}
	}

	err = p.c.Insert(bson.M{
		"_id": atomic.AddInt64(&p.count, int64(1)),
		"name": user.Name,
		"nickname": user.Nickname,
		"email": user.Email,
		"password": user.Password,
		"balance": user.Balance,
		"iddeleted": user.IsDeleted,
	})

	if err != nil {
		return &pb.Error{Code:102, Description:err.Error()}
	}

	return nil
}

// Count returns the number of document
// If you do not use name as _id, can use
// this to achieve self_increase _id
func (p *UserDB) Count() (n int, err error){
	n, err = p.c.Find(bson.M{}).Count()
	return
}
