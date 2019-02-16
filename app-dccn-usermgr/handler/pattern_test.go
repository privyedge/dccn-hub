package handler

import (
	"testing"
)

func TestVerifyEmailPattern(t *testing.T) {
	okEmails := []string{
		"162@163.com",
		"abc@gmail.com",
	}
	for _, email := range okEmails {
		if !matchPattern(OpEmailMatch, email) {
			t.Fatal("verify email failed ", email)
		}
	}

	errEmails := []string{
		"123456",
		"163.com",
		"@",
		"123@.com",
		"@163.com",
	}
	for _, email := range errEmails {
		if matchPattern(OpEmailMatch, email) {
			t.Fatal("verify email failed ", email)
		}
	}
}

func TestVerifyUserPattern(t *testing.T) {
	okNames := []string{
		"a23_12kjk",
		"aaaaaakjl",
		"a787098070",
		"a__________",
		"a9___987",
		"aa___uio",
	}
	for _, name := range okNames {
		if !matchPattern(OpUserNameMatch, name) {
			t.Fatal("verify user pattern failed ", name)
		}
	}

	errNames := []string{
		"a2_98sf09f",
		"2222222a_",
		"_akj9879879",
		"a2_",
	}
	for _, name := range errNames {
		if matchPattern(OpUserNameMatch, name) {
			t.Fatal("verify user failed ", name)
		}
	}
}
func TestVerifyPasswordPattern(t *testing.T) {
	okPasswords := []string{
		"a89787asfA_KjKJ",
		"A89787asfAKjKJ",
		"9A89787asfAKjKJ",
	}
	for _, pw := range okPasswords {
		if !matchPattern(OpPasswordMatch, pw) {
			t.Fatal("verify password pattern failed ", pw)
		}
	}
	errPasswords := []string{
		"123aA",
		"123aAksAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
		"12367890765",
		"aaaaaaaaaaaa",
		"AAAAAAAAAAAAAAAA",
		"1AAAAAAAAAAAAAA",
		"1aaaaaaaaaaaaaaa",
	}
	for _, pw := range errPasswords {
		if matchPattern(OpPasswordMatch, pw) {
			t.Fatal("verify password pattern failed ", pw)
		}
	}
}

func TestPattern(t *testing.T) {
	// re := regexp.MustCompilePOSIX(`^[0-9A-Za-z_]{3,5}$`)
	// re := regexp.MustCompile(`^[0-9A-Za-z_]{3,5}$`)
	// var re *regexp.Regexp
	// re = regexp.MustCompile(`^(\w){3,5}$`)
	okStr := []string{
		"abc",
		"_____",
		"0000",
		"0abc",
		"0abc_",
		"_abc9",
		"ab_c9",
	}
	for _, str := range okStr {
		// if !re.MatchString(str) {
		if !matchPattern(OpTestMatch, str) {
			t.Fatal("error")
		}
	}

	errStrs := []string{
		"",
		"0",
		"a",
		"_",
		"000000",
		"______",
		"aaaaaaa",
	}
	for _, str := range errStrs {
		// if re.MatchString(str) {
		if matchPattern(OpTestMatch, str) {
			t.Fatal("error ", str)
		}
	}
}
