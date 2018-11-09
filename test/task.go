package main

import (
        "fmt"
        mongodb "dccn-hub/util"
)


func main() {
    fmt.Println("test add test:")

  // //  mongodb.TaskList()
  //
  // task := mongodb.Task{Name:"docker_image_name", Region:"us_west", Zone:"ca", Userid:1}
  //
  //
  // id := mongodb.AddTask(task)
  // fmt.Printf("new task id : %d \n", id)

   userid := 1
   tasks := mongodb.TaskList(userid)

   for i := range tasks {
       task := tasks[i]
       fmt.Printf("task id : %d \n", task.ID)
}




}
