package services

import (
	"context"
	"errors"
	"time"

	"zipride/database"
	"zipride/internal/models"
	dmodels "zipride/internal/models"
	"zipride/utils"

	"gorm.io/gorm"
)

type AuthService struct {
	DB         *gorm.DB
	OtpService *OtpService
}

func NewAuthService() *AuthService {
	return &AuthService{
		DB:         database.DB,
		OtpService: NewOtpService(),
	}
}

// SignupWithEmail creates a driver record with hashed password and sends OTP to verify phone.
func (s *AuthService) SignupWithEmail(ctx context.Context, firstName, lastName, email, password, phone string) (*dmodels.Driver, error) {
	if phone == "" {
		return nil, errors.New("phone is required")
	}

	// check existing by email or phone
	var existing dmodels.Driver
	if err := s.DB.Where("email = ? OR phone = ?", email, phone).First(&existing).Error; err == nil {
		return nil, errors.New("user with given email or phone already exists")
	} else if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	hashed, err := utils.GenerateHash(password)
	if err != nil {
		return nil, err
	}

	driver := &dmodels.Driver{
		FirstName:  firstName,
		LastName:   lastName,
		Email:      &email,
		Phone:      &phone,
		Password:   hashed,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.DB.Create(driver).Error; err != nil {
		return nil, err
	}

	// send OTP for phone verification
	if err := s.OtpService.SendOTP(ctx, phone); err != nil {
		return nil, err
	}

	return driver, nil
}

func (s *AuthService) FinalizeOtpVerification(ctx context.Context, phone string) (*dmodels.Driver, error) {
	var driver dmodels.Driver
	if err := s.DB.Where("phone = ?", phone).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("driver not found")
		}
		return nil, err
	}

	driver.IsVerified = true
	driver.UpdatedAt = time.Now()
	if err := s.DB.Save(&driver).Error; err != nil {
		return nil, err
	}

	return &driver, nil
}

func (s *AuthService) LoginWithEmail(email, password string) (string, error) {
	var driver dmodels.Driver
	if err := s.DB.Where("email = ?", email).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", errors.New("driver not found")
		}
		return "", err
	}

	if !utils.CheckPass(driver.Password, password) {
		return "", errors.New("incorrect password")
	}

	if driver.Email == nil {
		return "", errors.New("driver email not found")
	}

	token, _ := utils.GenerateAccess(driver.ID, *driver.Email, driver.Role)
	return token, nil
}

func (s *AuthService) CompleteProfile(ctx context.Context, driverID uint, fullName, licenseNumber string, vehicle models.Vehicle) error {
	var driver models.Driver
	if err := s.DB.First(&driver, driverID).Error; err != nil {
		return err
	}

	// Example: update driver profile fields
	driver.FirstName = fullName
	// You may want to add LicenseNumber and Vehicle fields to your Driver model if not present
	// driver.LicenseNumber = licenseNumber
	// driver.Vehicle = vehicle

	driver.UpdatedAt = time.Now()
	if err := s.DB.Save(&driver).Error; err != nil {
		return err
	}

	// Optionally, save vehicle info if you have a vehicles table
	// vehicle.DriverID = driverID
	// if err := s.DB.Create(&vehicle).Error; err != nil {
	//     return err
	// }

	return nil
}
