package main

import (
	"fmt"
	mongodb "github.com/Ankr-network/dccn-hub/util"
)

func main() {
	fmt.Printf("test user  \n")
	user := mongodb.User{Name: "John", Token: "ed1605e17374bde6c68864d072c9f5c9", Money: 1000}

	mongodb.AddUser(user)

	//  token := "ed1605e17374bde6c68864d072c9f5c9"
	//
	//  user := mongodb.GetUser(token)
	//
	// if user.ID == 0 {
	//   fmt.Printf("user not exit \n")
	// } else{
	//   fmt.Printf("user id %d  name %s \n", user.ID, user.Name)
	// }
	//

}
