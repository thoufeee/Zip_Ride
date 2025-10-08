package utils

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// password strength check

func PasswordStrength(pass string) bool {

	if len(pass) < 6 {
		return false
	}

	var hasupper, hasdigit, hasspecial bool

	for _, ch := range pass {

		switch {
		case unicode.IsUpper(ch):
			hasupper = true
		case unicode.IsDigit(ch):
			hasdigit = true
		case unicode.IsPunct(ch), unicode.IsSymbol(ch):
			hasspecial = true
		}
	}

	return hasupper && hasdigit && hasspecial
}

// hashing password

func GenerateHash(pass string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 10)
	return string(hash), err
}

// check password

func CheckPass(pass, newpass string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(pass), []byte(newpass))
	return err == nil
}
