package dbservice

import (
	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
	"gopkg.in/mgo.v2/bson"
)

type UserRecord struct {
	ID               string             `bson:"id"`
	Email            string             `bson:"email"`
	HashedPassword   string             `bson:"hashed_password"`
	Name             string             `bson:"name"`
	Token            string             `bson:"token"`
	Status           usermgr.UserStatus `bson:"status"`
	LastModifiedDate uint64             `bson:"last_modified_date"`
	CreationDate     uint64             `bson:"creation_date"`
	PubKey           string             `bson:"pub_key"`
}

var feileds = map[string]string{
	"Id":               "id",
	"Email":            "email",
	"HashedPassword":   "hashed_password",
	"Name":             "name",
	"Token":            "token",
	"Status":           "status",
	"LastModifiedDate": "last_modified_date",
	"CreationDate":     "creation_date",
	"PubKey":           "pub_key",
}

func getUpdate(fields []*usermgr.UserAttribute) bson.M {
	update := bson.M{}
	for _, attr := range fields {
		switch attr.Key {
		case "Id":
			update[feileds[attr.Key]] = attr.GetStringValue()
		case "Email":
			update[feileds[attr.Key]] = attr.GetStringValue()
		case "HashedPassword":
			update[feileds[attr.Key]] = attr.GetStringValue()
		case "Name":
			update[feileds[attr.Key]] = attr.GetStringValue()
		case "Token":
			update[feileds[attr.Key]] = attr.GetStringValue()
		case "Status":
			update[feileds[attr.Key]] = usermgr.UserStatus(attr.GetIntValue())
		case "LastModifiedDate":
			update[feileds[attr.Key]] = attr.GetIntValue()
		case "CreationDate":
			update[feileds[attr.Key]] = attr.GetIntValue()
		case "PubKey":
			update[feileds[attr.Key]] = attr.GetStringValue()
		}
	}

	return bson.M{"$set": update}
}
