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
	PublicKeys       string
}

func getUpdate(fields []*usermgr.UserAttribute) bson.M {
	update := bson.M{}
	for _, attr := range fields {
		switch attr.Key {
		case "id":
			update[attr.Key] = attr.GetStringValue()
		case "email":
			update[attr.Key] = attr.GetStringValue()
		case "hashedpassword":
			update[attr.Key] = attr.GetStringValue()
		case "name":
			update[attr.Key] = attr.GetStringValue()
		case "token":
			update[attr.Key] = attr.GetStringValue()
		case "status":
			update[attr.Key] = usermgr.UserStatus(attr.GetIntValue())
		case "lastmodifieddate":
			update[attr.Key] = attr.GetIntValue()
		case "creationdate":
			update[attr.Key] = attr.GetIntValue()
		case "publickeys":
			update[attr.Key] = attr.GetStringValue()
		}
	}

	return bson.M{"$set": update}
}
