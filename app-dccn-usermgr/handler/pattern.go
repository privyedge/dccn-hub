package handler

import (
	"log"
	"regexp"
)

type PatternType int

const (
	// email pattern
	// emailPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	emailPattern = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	// 8-32 chars, [0-9A-Za-z_]
	// passwordPattern = `^(\w){8,32}$`
	// passwordPattern = `^(\w){8,32}$`
	// passwordPattern = `(=?.{8,32})((:?[a-z]+)(:?[0-9]+)(:?[A-Z]+)(:?\W+))`
	// passwordPattern = `^(=?.{8,32})((:?[a-z]+)(:?[0-9]+)(:?[A-Z]+))$`
	passwordPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
	// 6-18 chars, begin with chars, [0-9A-Za-z_] is available
	// userNamePattern = `^([a-zA-Z]){1}(\w){5-17}$`
	userNamePattern = `^(\w){1}(\w){5-17}$`
	testPattern     = `^(\w){3,5}$`
)

const (
	OpEmailMatch PatternType = iota
	OpPasswordMatch
	OpUserNameMatch
	OpTestMatch
)

func matchPattern(patternType PatternType, str string) bool {
	log.Println("Get String ", str)
	var reg *regexp.Regexp
	switch patternType {
	case OpEmailMatch:
		reg = regexp.MustCompile(emailPattern)
	case OpPasswordMatch:
		reg = regexp.MustCompile(passwordPattern)
	case OpUserNameMatch:
		reg = regexp.MustCompile(userNamePattern)
	case OpTestMatch:
		reg = regexp.MustCompile(testPattern)
	}
	return reg.MatchString(str)
}
