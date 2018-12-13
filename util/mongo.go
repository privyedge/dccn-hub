package util

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var MongoDBHost = "127.0.0.1"

type Task struct {
	ID           int64
	Userid       int64
	Name         string
	Datacenter   string
	Type         string
	Replica      int
	Datacenterid int64  // mongodb name is low field
	Status       string // 1 new 2 running 3. done 4 cancelling 5.canceled 6. updating 7. updateFailed
	Uniquename   string
	URL          string
}

type User struct {
	ID    int64
	Name  string
	Token string
	Money int64
}

type Counter struct {
	ID            string
	Sequencevalue int64
}

type DataCenter struct {
	ID             int64
	Name           string
	Report         string
	LastReportTime int64
	Status         string //1. running  2. stopped  3. dropped
	DatacenterId   int
}

func GetDataCenter(name string) DataCenter {
	datacenter := DataCenter{}
	db := GetDBInstance()
	c := db.C("datacenter")
	c.Find(bson.M{"name": name}).One(&datacenter)
	return datacenter
}

func AddDataCenter(d DataCenter) int64 {
	db := GetDBInstance()
	c := db.C("datacenter")
	//p := Person{"xxxx", "123455"}
	// p._id = 19
	// fmt.Printf("Id of person: %d\n", p._id)
	id := GetID("datacenterid", db)
	msec := time.Now().UnixNano() / 1000000
	err := c.Insert(bson.M{"_id": id, "id": id, "name": d.Name, "report": d.Report, "lastReporTtime": msec, "status": "Running"})
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func UpdateDataCenter(d DataCenter, id int) {
	db := GetDBInstance()
	c := db.C("datacenter")
	fmt.Printf("UpdateDataCenter report %s \n", d.Report)

	c.Update(bson.M{"_id": id}, bson.M{"$set": bson.M{"name": d.Name, "report": d.Report}})
}

var instance *mgo.Database
var once sync.Once

func GetDBInstance() *mgo.Database {
	once.Do(func() {
		instance = mongodbconnect()
	})
	return instance
}

func mongodbconnect() *mgo.Database {
	fmt.Printf("where db name user : %s\n", MongoDBHost)
	session, err := mgo.Dial(MongoDBHost)
	if err != nil {
		panic(err)
	}
	//defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	db := session.DB("test")
	return db

}

func AddTask(task Task) int64 {
	db := GetDBInstance()
	c := db.C("task")
	id := GetID("taskid", db)
	err := c.Insert(bson.M{"_id": id, "id": id, "name": task.Name, "userid": task.Userid, "type": task.Type, "datacenter": task.Datacenter, "status": "new"})
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func GetID(name string, db *mgo.Database) int64 {
	c := db.C("counters")
	result := Counter{}
	err := c.Find(bson.M{"_id": name}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s id  %d \n", name, result.Sequencevalue)
	id := result.Sequencevalue
	id += 1
	c.Update(bson.M{"_id": name}, bson.M{"$set": bson.M{"sequencevalue": id}})

	return result.Sequencevalue
}

func TaskList(userid int) []Task {
	var tasks []Task
	db := GetDBInstance()
	c := db.C("task")
	result := Task{}
	iter := c.Find(bson.M{"userid": userid}).Limit(100).Iter()
	for iter.Next(&result) {
		tasks = append(tasks, result)
	}
	return tasks
}

func DataCeterList() []DataCenter {
	var dataCenters []DataCenter
	db := GetDBInstance()
	c := db.C("datacenter")
	result := DataCenter{}
	iter := c.Find(nil).Limit(100).Iter()
	for iter.Next(&result) {
		dataCenters = append(dataCenters, result)
	}
	return dataCenters
}

func GetTask(taskid int) Task {
	task := Task{}
	db := GetDBInstance()
	c := db.C("task")
	c.Find(bson.M{"_id": taskid}).One(&task)
	return task
}

func GetNewTask() Task {
	task := Task{}
	db := GetDBInstance()
	c := db.C("task")
	c.Find(bson.M{"status": "new"}).One(&task)
	return task
}

func UpdateTask(taskid int, status string, datacentrid int) {
	db := GetDBInstance()
	c := db.C("task")
	if datacentrid == 0 {
		c.Update(bson.M{"_id": taskid}, bson.M{"$set": bson.M{"status": status}})
	} else {
		c.Update(bson.M{"_id": taskid}, bson.M{"$set": bson.M{"status": status, "datacenterid": datacentrid}})
	}
}

func UpdateTaskReplica(taskid int, replica int) {
	db := GetDBInstance()
	c := db.C("task")
	c.Update(bson.M{"_id": taskid}, bson.M{"$set": bson.M{"replica": replica}})
}

func UpdateTaskUnqueName(taskid int, uniqueName string) {
	db := GetDBInstance()
	c := db.C("task")
	c.Update(bson.M{"_id": taskid}, bson.M{"$set": bson.M{"uniquename": uniqueName}})
}

func UpdateTaskImage(taskid int, image string) {
	db := GetDBInstance()
	c := db.C("task")
	c.Update(bson.M{"_id": taskid}, bson.M{"$set": bson.M{"name": image}})
}

func UpdateTaskURL(taskid int, url string) {
	db := GetDBInstance()
	c := db.C("task")
	c.Update(bson.M{"_id": taskid}, bson.M{"$set": bson.M{"url": url}})
}

func CancelTask(taskid int) {
	db := GetDBInstance()
	c := db.C("task")
	c.Update(bson.M{"_id": taskid}, bson.M{"$set": bson.M{"status": "cancel"}})
}

func AddUser(user User) {

	db := GetDBInstance()
	c := db.C("user")
	id := GetID("userid", db)
	fmt.Printf("Id of user: %d\n", id)
	c.Insert(bson.M{"_id": id, "id": id, "name": user.Name, "token": user.Token, "money": user.Money})

}

func SelectFreeDatacenter() int {
	dcIds := []int{1, 2}
	index := rand.Intn(len(dcIds))
	return dcIds[index]

}

func GetDatacentersMap() map[int64]string {
	var dcs map[int64]string = map[int64]string{}
	dclist := DataCeterList()
	for i := range dclist {
		dc := dclist[i]
		dcs[dc.ID] = dc.Name
	}
	return dcs
}
func GetDatacenter(name string) DataCenter {
	dc := DataCenter{}
	db := GetDBInstance()
	c := db.C("datacenter")
	c.Find(bson.M{"name": name}).One(&dc)
	return dc
}

func GetUser(token string) User {
	user := User{}
	db := GetDBInstance()
	c := db.C("user")
	c.Find(bson.M{"token": token}).One(&user)
	return user
}

func GetTaskNameAsTaskIDForK8s(t Task) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	name := reg.ReplaceAllString(t.Name, "")

	return name + "-" + strconv.Itoa(int(t.ID))
}

func GetTaskIDFromTaskNameForK8s(name string) int64 {
	s := strings.Split(name, "-")
	if len(s) == 2 {
		value, err := strconv.Atoi(s[1])
		if err != nil {
			// handle error
			fmt.Println(err)
			os.Exit(2)
		}
		return int64(value)
	} else {
		return 0
	}
}
