package services

import (
	"errors"
	"zipride/utils"
)

func HashPassword(pass string) (string, error) {
	if !utils.PasswordStrength(pass) {
		return "", errors.New("password too weak")
	}
	return utils.GenerateHash(pass)
}

func CheckPassword(hash, pass string) bool {
	return utils.CheckPass(hash, pass)
}
