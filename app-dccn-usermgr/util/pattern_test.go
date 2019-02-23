package user_util

import (
	"log"
	"testing"
)

func TestCheckPassword(t *testing.T) {
	errPassword := []string{
		"abceddsf",
		"xxxxxx",
		"1234567",
	}

	for _, p := range errPassword {
		if err := CheckPassword(p); err != nil {
			log.Fatal(err.Error())
		}
	}
}
