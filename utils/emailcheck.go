package utils

import "net/mail"

// email check
func EmailCheck(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
