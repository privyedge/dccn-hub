package main

import (
        "fmt"
        mongodb "dccn-hub/util"
)


func main() {
    fmt.Println("test add test:")

  //  mongodb.TaskList()

  task := mongodb.Task{Name:"docker_image_name", Region:"us_west", Zone:"ca", Userid:1} }


  mongodb.addTask(task)

}
