package services

import (
	"errors"
	"zipride/database"
	"zipride/internal/domain/driver/models"
	"zipride/utils"
)

func GoogleLogin(token string) (string, error) {
	googleUser, err := utils.VerifyGoogleToken(token)
	if err != nil {
		return "", errors.New("invalid google token")
	}

	var d models.Driver
	if err := database.DB.Where("email = ?", googleUser.Email).First(&d).Error; err != nil {
		// create new driver
		d = models.Driver{
			Name:            googleUser.FirstName,
			Email:           &googleUser.Email,
			IsPhoneVerified: true,
		}
		if err := database.DB.Create(&d).Error; err != nil {
			return "", errors.New("failed to create driver")
		}
	}

	var email string
	if d.Email != nil {
		email = *d.Email
	} else {
		email = ""
	}
	return utils.GenerateAccess(d.ID, email, "driver", nil)
}
