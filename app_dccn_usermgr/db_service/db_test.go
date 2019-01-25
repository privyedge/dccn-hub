package dbservice

import (
	"testing"

	dbcommon "github.com/Ankr-network/dccn-common/db"
	go_micro_srv_usermgr "github.com/Ankr-network/dccn-common/protos/usermgr/v1/micro"
)

func mockDB() (DBService, error) {
	conf := dbcommon.Config{
		DB:         "114feb0961f8edfa8f514b67c6ef8af3",
		Collection: "user",
		Host:       "127.0.0.1:27017",
		Timeout:    15,
		PoolLimit:  15,
	}

	return New(conf)
}

func TestDB_New(t *testing.T) {
	db, err := mockDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	db.Close()
}

func mockUser() *go_micro_srv_usermgr.User {
	return &go_micro_srv_usermgr.User{
		Name:     "user_test",
		Nickname: "test",
		Email:    "123@Gmail.com",
		Password: "123456",
		Balance:  99,
	}
}

func TestDB_Create(t *testing.T) {
	db, err := mockDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()
	defer db.dropCollection()

	if err = db.Create(mockUser()); err != nil {
		t.Fatal(err.Error())
	}
}

func isEqual(origin, dbUser *go_micro_srv_usermgr.User) bool {
	return origin.Email == dbUser.Email &&
		origin.Password == dbUser.Password &&
		origin.Name == dbUser.Name &&
		origin.Nickname == dbUser.Nickname &&
		origin.Id == dbUser.Id &&
		origin.Balance == dbUser.Balance &&
		origin.IsDeleted == dbUser.IsDeleted
}

func TestDB_Get(t *testing.T) {
	db, err := mockDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()
	defer db.dropCollection()

	user := mockUser()
	if err = db.Create(user); err != nil {
		t.Fatal(err.Error())
	}

	var u *go_micro_srv_usermgr.User
	u, err = db.Get(user.Email)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !isEqual(user, u) {
		t.Fatalf("want %+v, but %+v\n", *user, *u)
	}
	t.Logf("Get Ok: %#v\n", *u)
}

func TestDB_Update(t *testing.T) {
	db, err := mockDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()
	defer db.dropCollection()

	user := mockUser()
	if err = db.Create(user); err != nil {
		t.Fatal(err.Error())
	}

	user.Nickname = "12345"
	if err = db.Update(user); err != nil {
		t.Fatal(err.Error())
	}

	u, err := db.Get(user.Email)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !isEqual(user, u) {
		t.Fatalf("UPDATE DB ERROR")
	}
}
