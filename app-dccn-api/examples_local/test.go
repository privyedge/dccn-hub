package main

import (
	"encoding/base64"
	"log"
	"strings"
)

func main() {


		s := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NTA5OTAxODAsImp0aSI6InlvdXNvbmdAYW5rci5jb20iLCJpc3MiOiJhbmtyLm5ldHdvcmsifQ.aBfqwJH8n6erG2tfVcA06LiknY1ybzeyn5XWiGE2KVo"

 b := "{\"exp\":1550990180,\"jti\":\"yousong@ankr.com\",\"iss\":\"ankr.network\"}a"

	result:= base64.StdEncoding.EncodeToString([]byte((b)))

	log.Printf("--->%s<---", result)

	parts := strings.Split(s, ".")

	log.Printf("--->%s<---", parts[1])

	decoded, err := base64.StdEncoding.DecodeString(parts[1]+"==")
	if err != nil {
		log.Printf("this is a error %s", err.Error())
	}



	log.Printf("value %s \n", decoded)
}



func getIdFromToken(refreshToken string) string {
	parts := strings.Split(refreshToken, ".")
	if len(parts) != 3 {
		log.Printf("1111111")
		return ""
	}

	log.Printf("why  ->%s<----", parts[1])

	decoded, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("22222 error %s %+v", err.Error(), decoded)
		return ""
	}
	log.Printf("4444")
	return "ok if fire"

}