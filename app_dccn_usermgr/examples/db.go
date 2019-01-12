package main

import (
	"log"

	go_micro_srv_usermgr "github.com/Ankr-network/dccn-hub/app_dccn_usermgr/proto/usermgr"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {
	s, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		log.Fatal(err.Error())
	}

	var u go_micro_srv_usermgr.User
	if err = s.DB("dccn").C("user").Find(bson.M{"email": "123@gmail.com"}).One(&u); err != nil {
		log.Fatal(err.Error())
	}

	log.Printf("%#v\n", u)
}
