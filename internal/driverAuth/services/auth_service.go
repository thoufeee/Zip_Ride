package services

import (
	"errors"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"gorm.io/gorm"
)

func RegisterDriver(firstName, lastName, email, phone, password string) (string, string, error) {
	if !utils.EmailCheck(email) {
		return "", "", errors.New("invalid email")
	}

	// Check if driver already exists
	var existing models.Driver
	err := database.DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return "", "", errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", "", err
	}

	// Hash password
	hashed, err := utils.GenerateHash(password)
	if err != nil {
		return "", "", errors.New("failed to hash password")
	}

	// Create new driver
	newDriver := models.Driver{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      email,
		Phone:      phone,
		Password:   hashed,
		Role:       "driver",
		IsVerified: true,
	}

	if err := database.DB.Create(&newDriver).Error; err != nil {
		return "", "", errors.New("failed to create driver")
	}

	// Generate JWT access token
	accessToken, err := utils.GenerateAccess(newDriver.ID, email, "driver")
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateRefresh(newDriver.ID, email, "driver")
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

func LoginDriver(phone, password string) (string, string, error) {
	var d models.Driver
	err := database.DB.Where("phone = ?", phone).First(&d).Error
	if err != nil {
		return "", "", errors.New("driver not found")
	}

	if !utils.CheckPass(d.Password, password) {
		return "", "", errors.New("incorrect password")
	}

	accessToken, err := utils.GenerateAccess(d.ID, d.Email, "driver")
	if err != nil {
		return "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefresh(d.ID, d.Email, "driver")
	if err != nil {
		return "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, nil
}

// EnsureDriverByPhoneAndIssueToken ensures a driver exists for the phone and issues tokens
func EnsureDriverByPhoneAndIssueToken(phone string) (string, string, string, error) {
	var d models.Driver
	err := database.DB.Where("phone = ?", phone).First(&d).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// create minimal driver record
		d = models.Driver{
			FirstName:     "",
			LastName:      "",
			Email:         "",
			Phone:         phone,
			Password:      "",
			Role:          "driver",
			IsVerified:    true,
			Status:        "pending_docs",
			PhoneVerified: true,
		}
		if err := database.DB.Create(&d).Error; err != nil {
			return "", "", "", errors.New("failed to create driver")
		}
	} else if err != nil {
		return "", "", "", err
	} else {
		// mark phone verified
		d.PhoneVerified = true
		if err := database.DB.Save(&d).Error; err != nil {
			return "", "", "", err
		}
	}

	accessToken, err := utils.GenerateAccess(d.ID, d.Email, "driver")
	if err != nil {
		return "", "", "", errors.New("failed to generate access token")
	}

	refreshToken, err := utils.GenerateRefresh(d.ID, d.Email, "driver")
	if err != nil {
		return "", "", "", errors.New("failed to generate refresh token")
	}

	return accessToken, refreshToken, d.Status, nil
}

