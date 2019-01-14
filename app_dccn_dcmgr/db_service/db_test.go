package dbservice

import (
	"testing"

	dcmgr "github.com/Ankr-network/dccn-hub/app_dccn_dcmgr/proto/v1"
	dbcommon "github.com/Ankr-network/dccn-hub/common/db"
)

func mockDB() (DBService, error) {
	conf := dbcommon.Config{
		DB:         "114feb0961f8edfa8f514b67c6ef8af3",
		Collection: "dc",
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

func mockDc() *dcmgr.DataCenter {
	return &dcmgr.DataCenter{
		Id:     1,
		Name:   "dc01",
		Status: 1,
	}
}

func TestDB_Create(t *testing.T) {
	db, err := mockDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()
	defer db.dropCollection()

	if err = db.Create(mockDc()); err != nil {
		t.Fatal(err.Error())
	}
}

func isEqual(origin, dbdc *dcmgr.DataCenter) bool {
	return origin.Id == dbdc.Id && origin.Status == dbdc.Status && origin.Name == dbdc.Name
}

func TestDB_Get(t *testing.T) {
	db, err := mockDB()
	if err != nil {
		t.Fatal(err.Error())
	}
	defer db.Close()
	defer db.dropCollection()

	datacenter := mockDc()
	if err = db.Create(datacenter); err != nil {
		t.Fatal(err.Error())
	}

	var u *dcmgr.DataCenter
	u, err = db.Get(datacenter.Id)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !isEqual(datacenter, u) {
		t.Fatalf("want %+v, but %+v\n", *datacenter, *u)
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

	dataCenter := mockDc()
	if err = db.Create(dataCenter); err != nil {
		t.Fatal(err.Error())
	}

	dataCenter.Name = "000000"
	if err = db.Update(dataCenter); err != nil {
		t.Fatal(err.Error())
	}

	u, err := db.Get(dataCenter.Id)
	if err != nil {
		t.Fatal(err.Error())
	}

	if !isEqual(dataCenter, u) {
		t.Fatalf("UPDATE DB ERROR")
	}
}
