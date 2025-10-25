package services

import (
	"errors"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID  uint     `json:"uid"`
	Email   string   `json:"email"`
	Subject string   `json:"subj"`
	Roles   []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uint, email, subject string, roles []string, secret string, expiryMinutes int) (string, error) {
	now := time.Now().UTC()
	claims := &Claims{
		UserID:  userID,
		Email:   email,
		Subject: subject,
		Roles:   roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenStr, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
