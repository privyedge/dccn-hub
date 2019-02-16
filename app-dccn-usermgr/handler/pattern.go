package handler

import (
	"regexp"
)

type PatternType int

const (
	// email pattern
	emailPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	// 8-32 chars, [0-9A-Za-z_]
	passwordPattern = `^(\w){8-32}$`
	// 6-18 chars, begin with chars, [0-9A-Za-z_] is available
	userNamePattern = `^([a-zA-Z]){1}(\w){5-17}$`
)

const (
	OpEmailMatch PatternType = iota
	OpPasswordMatch
	OpUserNameMatch
)

func matchPattern(patternType PatternType, str string) bool {
	var reg *regexp.Regexp
	switch patternType {
	case OpEmailMatch:
		reg = regexp.MustCompile(emailPattern)
	case OpPasswordMatch:
		reg = regexp.MustCompile(passwordPattern)
	case OpUserNameMatch:
		reg = regexp.MustCompile(userNamePattern)
	}
	return reg.MatchString(str)
}
