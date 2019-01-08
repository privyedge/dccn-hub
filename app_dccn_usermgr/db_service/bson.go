package dbservice

import (
	"errors"
	"github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
)

type bsonUser struct {
	// Id self-increase
	Id int64 `json:"id,omitempty" bson:"_id"`
	// Name should be unique
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	// Nickname show on UI
	Nickname string `json:"nickname,omitempty" bson:"nickname,omitempty"`
	// Email user's email
	Email string `json:"email,omitempty" bson:"email,omitempty"`
	// Password string
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	// Balance user's balance in account
	Balance int32 `json:"balance,omitempty" bson:"balance,omitempty"`
	// IsDeleted user's status
	IsDeleted            bool     `json:"is_deleted,omitempty" bson:"is_deleted,omitempty"`
}

func (buser *bsonUser) Encode(user *go_micro_srv_usermgr.User) error {
	if buser == nil || user == nil {
		return errors.New("Null Pointer")
	}
	buser.Id = user.Id.Id
	buser.Name = user.Name
	buser.Nickname = user.Nickname
	buser.Email = user.Email
	buser.Password = user.Password
	buser.Balance = user.Balance
	buser.IsDeleted = user.IsDeleted

	return nil
}

func (buser *bsonUser) Decode(user *go_micro_srv_usermgr.User) error {
	if buser == nil || user == nil {
		return errors.New("Null Pointer")
	}

	user.Id.Id = buser.Id
	user.Name = buser.Name
	user.Nickname = buser.Nickname
	user.Email = buser.Email
	user.Password = buser.Password
	user.Balance = buser.Balance
	user.IsDeleted = buser.IsDeleted

	return nil
}
