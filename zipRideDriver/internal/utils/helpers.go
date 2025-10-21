package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(plain string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(plain), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func CheckPassword(hash, plain string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) == nil
}

// CheckPasswordHash is an alias for CheckPassword for compatibility
func CheckPasswordHash(plain, hash string) bool {
	return CheckPassword(hash, plain)
}

// GenerateJWT generates a JWT token for a driver
func GenerateJWT(driverID uint, email string, secret string) (string, error) {
	claims := jwt.MapClaims{
		"driver_id": driverID,
		"email":     email,
		"exp":       time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days expiry
		"iat":       time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
