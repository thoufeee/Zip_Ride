package utils

import (
	"os"
	"time"
	"zipride/database"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// claims

type Claims struct {
	UserId  uint   `json:"userid"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	TokenId string `json:"tokenid"`
	jwt.RegisteredClaims
}

// generate access token

func GenerateAccess(userid uint, email string, role string) (string, error) {

	jwtkey := []byte(os.Getenv("JWT_SECRET"))

	expir := time.Now().Add(60 * time.Minute)

	claims := &Claims{
		UserId: userid,
		Email:  email,
		Role:   role,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expir),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtkey)

}

// generate refresh token

func GenerateRefresh(userid uint, email, role string) (string, error) {

	refreshkey := []byte(os.Getenv("REFRESH_KEY"))

	exp := time.Now().Add(7 * 24 * time.Hour)
	refreshid := uuid.New().String()

	claims := &Claims{
		UserId:  userid,
		Email:   email,
		Role:    role,
		TokenId: refreshid,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
			ID:        refreshid,
		},
	}

	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := refresh.SignedString(refreshkey)

	if err != nil {
		return "", err
	}

	err = database.RDB.Set(database.Ctx, refreshid, token, time.Until(exp)).Err()

	if err != nil {
		return "", err
	}

	return token, nil
}
