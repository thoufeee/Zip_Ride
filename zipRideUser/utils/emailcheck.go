package utils

import (
	"net/mail"
	"regexp"
	"strings"
)

// email check
func EmailCheck(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// phone number chcek
func PhoneNumberCheck(phone string) (string, bool) {
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.TrimPrefix(phone, "+91")
	phone = strings.TrimPrefix(phone, "0")

	//    checking all charaters are digits and exactly ten
	matched, _ := regexp.MatchString(`^[0-9]{10}$`, phone)

	if !matched {
		return "", false
	}

	return phone, true
}
