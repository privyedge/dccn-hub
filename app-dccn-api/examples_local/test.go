package main

import (
	"encoding/base64"
	"log"
    "strings"

)

func main() {


		s := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTA5ODM0MzcsImp0aSI6IjQyZjQ5Y2FlLTAzYzUtNDgzOS04OWI3LTllMzdjMmZjNTk1ZSIsImlzcyI6ImFua3IubmV0d29yayJ9.qfwSuyfOzWLywUDgKHWp1Ur29Sv50Tf8SRT40IEjkLk"


	parts := strings.Split(s, ".")

	log.Printf("pargts ->%s<-", parts[1])

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("this is a error %s", err.Error())
	}

	log.Printf("value %s \n", decoded)
}
