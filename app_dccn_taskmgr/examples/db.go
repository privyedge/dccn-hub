package main

import (
	"log"

	pb "github.com/Ankr-network/dccn-common/protos/common/proto/v1"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func GetAll(c *mgo.Collection) (*[]*pb.DataCenter, error) {
	dc := pb.DataCenter{
		Id:     0,
		Name:   "sanghai",
		Status: 0,
	}
	for i := 0; i < 5; i++ {
		if err := c.Insert(&dc); err != nil {
			dc.Status++
			log.Fatal(err.Error())
		}
	}

	var dcs []*pb.DataCenter
	if err := c.Find(bson.M{"id": 0}).All(&dcs); err != nil {
		log.Fatal(err.Error())
	}
	return &dcs, nil
}

func main() {
	s, err := mgo.Dial("localhost:27017")
	if err != nil {
		log.Fatal(err.Error())
	}

	c := s.DB("test").C("test")
	GetAll(c)
}
