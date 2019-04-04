package dbservice

import (
	"gopkg.in/mgo.v2/bson"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/grpc"

	ankr_default "github.com/Ankr-network/dccn-common/protos"
	user_util "github.com/Ankr-network/dccn-hub/app-dccn-usermgr/util"
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
	EmailChangeConfirmCode string
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

func getUpdate(fields []*usermgr.UserAttribute) (bson.M, error) {
	update := bson.M{}
	for _, attr := range fields {
		switch attr.Key {
		case "Email":
			update[feileds[attr.Key]] = attr.GetStringValue()
		case "HashedPassword":
			update[feileds[attr.Key]] = attr.GetStringValue()
		case "Name":
			if err := user_util.CheckName(attr.GetStringValue()); err != nil {
				return nil, ankr_default.ErrUserNameFormat
			}
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

	return bson.M{"$set": update}, nil
}
