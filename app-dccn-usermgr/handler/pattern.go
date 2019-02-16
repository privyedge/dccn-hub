package handler

import (
	"regexp"
)

type PatternType int

const (
	// email pattern
	emailPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	// 8-32 chars,  begin with chars, two of {lower case, tower case, digit} at least
	passwordPattern = `^(\w){8,32}$`
	// passwordPattern = `^(?![a-zA-Z]+$)(?![A-Z0-9]+$)(?![A-Z\\W_!@#$%^&*`~()-+=]+$)(?![a-z0-9]+$)(?![a-z\\W_!@#$%^&*`~()-+=]+$)(?![0-9\\W_!@#$%^&*`~()-+=]+$)[a-zA-Z0-9\\W_!@#$%^&*`~()-+=]{8,30}$`
	// passwordPattern = `^(![a-z])+$)(?![A-Z]+$)(?![0-9]+$)(?![a-zA-Z]+$)(?![a-z0-9]+$)(?![A-Z0-9]+$)[0-9A-Za-z]{8,31}$`
	// 6-18 chars, begin with chars, [0-9A-Za-z_] is available
	// userNamePattern = `^(\w){5-17}$`
	userNamePattern = `^[A-Za-z]{1}(\w){5,17}$`
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
