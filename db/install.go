package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	Name  string
	Phone string
}

func main() {
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("test").C("counters")
	c.Insert(bson.M{"_id": "userid", "sequencevalue": 10})
	c.Insert(bson.M{"_id": "taskid", "sequencevalue": 10})
	c.Insert(bson.M{"_id": "datacenterid", "sequencevalue": 10})

	u := session.DB("test").C("user")
	u.Insert(bson.M{"_id": 1, "id": 1, "name": "John", "token": "ed1605e17374bde6c68864d072c9f5c9", "money": 1000})

}
