package util

import (
        "fmt"
	      "log"
        "gopkg.in/mgo.v2"  // package name mgo
        "gopkg.in/mgo.v2/bson"
)

type Task struct {
        Taskid int64
        Userid int64
        Name string
        Region string
        Zone string

}

type User struct {
        Userid int64
        Name string
        usertoken string
        money int64
}

type Person struct {
        Name string
        Phone string
}

func mongodbconnect() *mgo.Session{
  session, err := mgo.Dial("127.0.0.1")
  if err != nil {
          panic(err)
  }
  //defer session.Close()

  // Optional. Switch the session to a monotonic behavior.
  session.SetMode(mgo.Monotonic, true)
  return session

}

func addTask(task Task) {
        session := mongodbconnect()
        c := session.DB("test").C(Task)
        err := c.Insert(&Person{"Ale", "+55 53 8116 9639"})
        if err != nil {
                log.Fatal(err)
        }
}

func TaskList(userid string){
  session := mongodbconnect()
  c := session.DB("test").C("people")
  result := Person{}
  iter := c.Find(bson.M{"userid": userid}).Limit(100).Iter()
  for iter.Next(&result) {
    fmt.Printf("Result: %s\n", result.Name)
   }
}
