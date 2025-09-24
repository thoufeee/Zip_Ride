package utils

import (
	"context"
	"os"
	"zipride/models"

	"google.golang.org/api/idtoken"
)

// verifying google token

func VerifyGoogleToken(token string) (*models.GoogleUser, error) {

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")

	payload, err := idtoken.Validate(context.Background(), token, googleClientId)

	if err != nil {
		return nil, err
	}

	return &models.GoogleUser{
		GoogleID:  payload.Claims["sub"].(string),
		Email:     payload.Claims["email"].(string),
		FirstName: payload.Claims["given_name"].(string),
		LastName:  payload.Claims["family_name"].(string),
		Avatar:    payload.Claims["picture"].(string),
	}, nil
}
