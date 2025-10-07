package services

import (
	"errors"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"gorm.io/gorm"
)

// GoogleLogin verifies the Google id_token, upserts a driver, and returns accessToken, refreshToken, status, and phoneVerified
func GoogleLogin(idToken string) (string, string, string, bool, error) {
	googleUser, err := utils.VerifyGoogleToken(idToken)
	if err != nil {
		return "", "", "", false, errors.New("invalid google token")
	}

	var d models.Driver
	// Prefer lookup by GoogleID if present
	if googleUser.GoogleID != "" {
		if err := database.DB.Where("google_id = ?", googleUser.GoogleID).First(&d).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", "", false, err
		}
	}

	// Fallback lookup by email if not found and email present
	if d.ID == 0 && googleUser.Email != "" {
		if err := database.DB.Where("email = ?", googleUser.Email).First(&d).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", "", false, err
		}
	}

	if d.ID == 0 {
		// Create new driver
		d = models.Driver{
			FirstName:     googleUser.FirstName,
			LastName:      googleUser.LastName,
			Email:         googleUser.Email,
			Role:          "driver",
			IsVerified:    true,
			Status:        "pending_docs",
			PhoneVerified: false,
			Avatar:        googleUser.Avatar,
		}
		if googleUser.GoogleID != "" {
			gid := googleUser.GoogleID
			d.GoogleID = &gid
		}
		if err := database.DB.Create(&d).Error; err != nil {
			return "", "", "", false, errors.New("failed to create driver")
		}
	} else {
		// Update missing linkage/fields
		updated := false
		if d.GoogleID == nil && googleUser.GoogleID != "" {
			gid := googleUser.GoogleID
			d.GoogleID = &gid
			updated = true
		}
		if d.Avatar == "" && googleUser.Avatar != "" {
			d.Avatar = googleUser.Avatar
			updated = true
		}
		if updated {
			if err := database.DB.Save(&d).Error; err != nil {
				return "", "", "", false, err
			}
		}
	}

	accessToken, err := utils.GenerateAccess(d.ID, d.Email, "driver")
	if err != nil {
		return "", "", "", false, errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefresh(d.ID, d.Email, "driver")
	if err != nil {
		return "", "", "", false, errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, d.Status, d.PhoneVerified, nil
}
