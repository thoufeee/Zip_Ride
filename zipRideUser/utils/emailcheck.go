package utils

import "net/mail"

// email check
func EmailCheck(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// phone number chcek
func PhoneNumberCheck(phone string) bool {
	return len(phone) == 10
}
