package models

import "time"

type Vehicle struct {
	ID              uint      `gorm:"primaryKey"`
	DriverID        uint      `gorm:"index;not null"`
	Make            string    `gorm:"size:80"`
	Model           string    `gorm:"size:100"`
	Year            int       `gorm:"index"`
	PlateNumber     string    `gorm:"size:50;index"`
	InsuranceNumber string    `gorm:"size:100"`
	RCNumber        string    `gorm:"size:100"`
	Status          string    `gorm:"size:30;default:'active';index"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
