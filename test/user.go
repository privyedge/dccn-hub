package main

import (
        "fmt"
        mongodb "dccn-hub/util"
)


func main() {
 //user := mongodb.User{Name:"John", Token: "ed1605e17374bde6c68864d072c9f5c9", Money:1000 }


 //mongodb.AddUser(user)
//id := mongodb.GetID()

  token := "ed1605e17374bde6c68864d072c9f5c9"
  user := mongodb.GetUser(token)
//
 fmt.Printf("user id %d  name %s \n", user.ID, user.Name)
 fmt.Println("done job \n")

}
