package dbservice

import (
	"log"
	"testing"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	mgo "gopkg.in/mgo.v2"

	usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func mockDbService() DBService {
	s, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatal(err.Error())
	}

	return &DB{
		dbName:         "sang",
		collectionName: "test",
		session:        s,
	}
}

func mockUserRecord() *UserRecord {
	return &UserRecord{
		ID:               uuid.New().String(),
		Name:             "test_local",
		Email:            "xuexiacm@163.com",
		LastModifiedDate: uint64(time.Now().Unix()),
		CreationDate:     uint64(time.Now().Unix()),
	}
}

func mockPbUser() *usermgr.User {
	return &usermgr.User{
		Id:     "fa010697-b2a2-4ce0-be4b-544f097a6822 ",
		Email:  "xuexiacm@163.com ",
		Status: usermgr.UserStatus_CONFIRMING,
		Attributes: &usermgr.UserAttributes{
			Name:             "ankrtest_sang",
			LastModifiedDate: uint64(time.Now().Unix()),
			CreationDate:     uint64(time.Now().Unix()),
		},
	}
}

func TestCreateUser(t *testing.T) {
	db := mockDbService()
	defer db.dropCollection()

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err.Error())
	}
	if err := db.CreateUser(mockPbUser(), string(hashedPwd)); err != nil {
		log.Fatal(err.Error())
	}
}

func TestUpdateUser(t *testing.T) {
	db := mockDbService()
	// defer db.dropCollection()

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte("12345"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err.Error())
	}

	u := mockPbUser()

	if err := db.CreateUser(u, string(hashedPwd)); err != nil {
		log.Fatal(err.Error())
	}

	attr := []*usermgr.UserAttribute{
		{
			Key:   "HashedPassword",
			Value: &usermgr.UserAttribute_StringValue{StringValue: "13456890"},
		},
	}
	if err := db.UpdateUser(u.Id, attr); err != nil {
		log.Fatal(err.Error())
	}

	emailAttr := []*usermgr.UserAttribute{
		{
			Key:   "Email",
			Value: &usermgr.UserAttribute_StringValue{StringValue: "994336359@qq.com"},
		},
		{
			Key: "Name",
		},
		{
			Key: "publickeys",
			// : "xiaowu",
			Value: &usermgr.UserAttribute_StringValue{StringValue: "xiaowu"},
		},
	}

	if err := db.UpdateUserByEmail(u.Email, emailAttr); err != nil {
		log.Fatal(err.Error())
	}

}
