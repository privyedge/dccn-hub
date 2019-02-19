package user_util

import (
	"regexp"
)

type PatternType int

const (
	// email pattern
	emailPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	// 8-32 chars,  begin with chars, two of {lower case, tower case, digit} at least
	passwordPattern = `^(\w){6,32}$`
	// 6-18 chars, begin with chars, [0-9A-Za-z_] is available
	userNamePattern = `^[A-Za-z]{1}(\w){5,17}$`
)

const (
	OpEmailMatch PatternType = iota
	OpPasswordMatch
	OpUserNameMatch
)

func MatchPattern(patternType PatternType, str string) bool {
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
