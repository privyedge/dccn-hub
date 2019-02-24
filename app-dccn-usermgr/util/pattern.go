package user_util

import (
	"log"
	"regexp"
	"unicode"
	"errors"

	"github.com/Ankr-network/dccn-common/protos"
)

const (
	// email pattern
	emailPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
)

func CheckEmail(email string) error {
	if len(email) < 4 {
		return errors.New("email too short")
	}

	_, err := regexp.MatchString(emailPattern, email)
	return err
}

// [0-9A-Za-z-]; at least 2, must contain en letter and digit
func CheckName(name string) error {
	//var (
	//	letter bool
	//	digit  bool
	//)

	if len(name) < 2 {
		return ankr_default.ErrUserNameFormat
	}

	//for _, c := range name {
	//	switch {
	//	case unicode.IsNumber(c):
	//		digit = true
	//	case unicode.IsLetter(c):
	//		letter = true
	//	default:
	//		return ankr_default.ErrUnexpectedChar
	//	}
	//
	//	if digit && letter {
	//		return nil
	//	}
	//}
	return nil
}

func CheckPassword(password string) error {
	log.Printf("CheckPassword %s \n", password)
	var (
		letter bool
		digit  bool
	)

	if len(password) < 6 {
		return ankr_default.ErrPasswordLength
	}

	for _, c := range password {
		switch {
		case unicode.IsNumber(c):
			digit = true
		case unicode.IsLetter(c):
			letter = true
		default:
			return ankr_default.ErrUnexpectedChar
		}

		if digit && letter {
			return nil
		}

	}

	return ankr_default.ErrPasswordFormat
}

func CheckRegister(name, email, password string) error {
	if err := CheckName(name); err != nil {
		return err
	}

	if err := CheckPassword(password); err != nil {
		return err
	}

	if err := CheckEmail(email); err != nil {
		return err
	}
	return nil
}
