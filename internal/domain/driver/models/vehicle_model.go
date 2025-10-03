package models

import "github.com/google/uuid"

type Vehicle struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DriverID    uuid.UUID
	Make        string `json:"make"`
	Model       string `json:"model"`
	PlateNumber string `json:"plate_number" binding:"required"`
	Year        int    `json:"year"`
}
