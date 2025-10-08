package utils

import (
	"net/mail"
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

    if len(phone) != 10 {
        return "", false
    }
    return phone, true
}
