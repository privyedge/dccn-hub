package dbaccount

import (
	pb "github.com/Ankr-network/refactor/app_dccn_account/proto"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

const (
	AccountCollection = "account"
)

// 101:  User already exists"
// 102: Other errors

// mysql or mongodb. who is better.
type AccountDBService interface {
	// New Creates a new account if not exits
	New(account *pb.Account) *pb.Error
	// Get get account by account's id
	Get(name string) (*pb.Account, *pb.Error)
}

// AccountDB implements AccountDBService
type AccountDB struct{
	mu sync.Mutex
	c *mgo.Collection
}

// New Creates a new account if not exits.
func New(c *mgo.Collection) *AccountDB {
	return &AccountDB{c: c}
}

// Get gets account by account's id.
func (p *AccountDB) Get(name string) (*pb.Account, error) {
	var result *pb.Account
	if err := p.c.FindId(name).One(result); err != nil {
		return nil, err
	}
	return result, nil
}

// New creates a new account if it not exists, name unique as _id
// TODO: batch operations throught bulk
func (p *AccountDB) New(account pb.Account) *pb.Error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// the UI ensures the name not empty
	ac, err := p.Get(account.Name)
	if err != nil {
		return &pb.Error{Code:102, Description: err.Error()}
	}

	if ac != nil {
		return &pb.Error{Code: 101, Description: "User already exists"}
	}

	err = p.c.Insert(bson.M{
		"_id": account.Name,
		"nickname": account.Nickname,
		"email": account.Email,
		"password": account.Password,
		"balance": account.Balance,
	})

	if err != nil {
		return &pb.Error{Code:102, Description:err.Error()}
	}

	return nil
}

// Count returns the number of document
// If you do not use name as _id, can use
// this to achieve self_increase _id
func (p *AccountDB) Count() (n int, err error){
	n, err = p.c.Find(bson.M{}).Count()
	return
}
