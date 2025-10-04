package services

import (
	"errors"
	"gorm.io/gorm"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"
)

func RegisterDriver(firstName, lastName, email, phone, password string) (string, error) {
	if !utils.EmailCheck(email) {
		return "", errors.New("invalid email")
	}

	// Check if driver already exists
	var existing models.Driver
	err := database.DB.Where("email = ?", email).First(&existing).Error
	if err == nil {
		return "", errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	// Hash password
	hashed, err := utils.GenerateHash(password)
	if err != nil {
		return "", errors.New("failed to hash password")
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
		return "", errors.New("failed to create driver")
	}

	// Generate JWT access token
	token, err := utils.GenerateAccess(newDriver.ID, email, "driver")
	if err != nil {
		return "", errors.New("failed to generate token")
	}
	return token, nil
}

func LoginDriver(phone, password string) (string, error) {
	var d models.Driver
	err := database.DB.Where("phone = ?", phone).First(&d).Error
	if err != nil {
		return "", errors.New("driver not found")
	}

	if !utils.CheckPass(d.Password, password) {
		return "", errors.New("incorrect password")
	}

	return utils.GenerateAccess(d.ID, d.Email, "driver")
}
