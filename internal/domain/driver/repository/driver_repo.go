package repository

import (
	"gorm.io/gorm"
	"zipride/internal/domain/driver/models"
)

type DriverRepo struct {
	db *gorm.DB
}

func NewDriverRepo(db *gorm.DB) *DriverRepo {
	return &DriverRepo{db: db}
}

func (r *DriverRepo) Create(driver *models.Driver) error {
	return r.db.Create(driver).Error
}

func (r *DriverRepo) FindByEmail(email string) (*models.Driver, error) {
	var driver models.Driver
	err := r.db.Where("email = ?", email).First(&driver).Error
	return &driver, err
}

func (r *DriverRepo) UpdateStatus(driverID string, status string) error {
	return r.db.Model(&models.Driver{}).Where("id = ?", driverID).Update("status", status).Error
}
