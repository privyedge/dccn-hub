package dbuser

import (
	"errors"
	"sync"
	"sync/atomic"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	pb "github.com/Ankr-network/refactor/proto/usermgr"
)

// mysql or mongodb. who is better.
type DBService interface {
	// New Creates a new usermgr if not exits
	New(user *pb.User) error
	// Get get usermgr by usermgr's id
	Get(name string) (*pb.User, error)
}

// UserDB implements DBService
type UserDB struct{
	mu sync.Mutex
	count int64
	c *mgo.Collection
}

// New Creates a new usermgr if not exits.
func New(c *mgo.Collection) *UserDB {
	count, _ := c.Find(bson.M{}).Count()
	return &UserDB{count: int64(count), c: c}
}

// Get gets usermgr by usermgr's id.
func (p *UserDB) Get(name string) (*pb.User, error) {
	var result *pb.User
	if err := p.c.Find(bson.M{"name": name}).One(result); err != nil {
		return nil, err
	}
	return result, nil
}

// New creates a new usermgr if it not exists, name unique as _id
// TODO: batch operations throught bulk
func (p *UserDB) New(user *pb.User) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// the UI ensures the name not empty
	u, err := p.Get(user.Name)
	if err != nil {
		return err
	}

	if u != nil {
		return errors.New("User already exists")
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
		return err
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
