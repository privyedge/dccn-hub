package dbservice

import (
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	"gopkg.in/mgo.v2/bson"
)

type UserRecord struct {
	ID               string
	Email            string
	HashedPassword   string
	Name             string
	Token            string             // refresh token
	Status           usermgr.UserStatus // user status
	LastModifiedDate uint64
	CreationDate     uint64
	PubKey           string
}

func getUpdate(fields []*usermgr.UserAttribute) bson.M {
	update := bson.M{}
	for _, attr := range fields {
		switch attr.Key {
		case "ID":
			update[attr.Key] = attr.GetStringValue()
		case "Email":
			update[attr.Key] = attr.GetStringValue()
		case "HashedPassword":
			update[attr.Key] = attr.GetStringValue()
		case "Name":
			update[attr.Key] = attr.GetStringValue()
		case "Token":
			update[attr.Key] = attr.GetStringValue()
		case "Status":
			update[attr.Key] = usermgr.UserStatus(attr.GetIntValue())
		case "LastModifiedDate":
			update[attr.Key] = attr.GetIntValue()
		case "CreationDate":
			update[attr.Key] = attr.GetIntValue()
		case "PubKey":
			update[attr.Key] = attr.GetStringValue()
		}
	}

	return bson.M{"$set": update}
}
