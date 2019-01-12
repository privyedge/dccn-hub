package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

var passwd = "123456"

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	hashpw, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err.Error())
	}
	encodePw := string(hashpw)
	log.Println(encodePw)
	err = bcrypt.CompareHashAndPassword([]byte(encodePw), []byte(hashpw))
	if err != nil {
		log.Fatal(err.Error())
	}
}
