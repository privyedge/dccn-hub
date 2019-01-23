package dbservice

import (
	"reflect"
	"testing"

	dbcommon "github.com/Ankr-network/dccn-common/db"
	dcmgr "github.com/Ankr-network/dccn-common/protos/common/proto/v1"
	mgo "gopkg.in/mgo.v2"
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

func TestNew(t *testing.T) {
	type args struct {
		conf dbcommon.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *DB
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_collection(t *testing.T) {
	type fields struct {
		dbName         string
		collectionName string
		session        *mgo.Session
	}
	type args struct {
		session *mgo.Session
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *mgo.Collection
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DB{
				dbName:         tt.fields.dbName,
				collectionName: tt.fields.collectionName,
				session:        tt.fields.session,
			}
			if got := p.collection(tt.args.session); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.collection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_Close(t *testing.T) {
	type fields struct {
		dbName         string
		collectionName string
		session        *mgo.Session
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DB{
				dbName:         tt.fields.dbName,
				collectionName: tt.fields.collectionName,
				session:        tt.fields.session,
			}
			p.Close()
		})
	}
}

func TestDB_dropCollection(t *testing.T) {
	type fields struct {
		dbName         string
		collectionName string
		session        *mgo.Session
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &DB{
				dbName:         tt.fields.dbName,
				collectionName: tt.fields.collectionName,
				session:        tt.fields.session,
			}
			p.dropCollection()
		})
	}
}
