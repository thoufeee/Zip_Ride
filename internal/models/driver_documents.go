package models

import (
	"time"
)

type DriverDocuments struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	DriverID     uint      `gorm:"not null;uniqueIndex" json:"driver_id"`
	LicenseURL   string    `json:"license_url"`
	RCURL        string    `json:"rc_url"`
	InsuranceURL string    `json:"insurance_url"`
	Status       string    `gorm:"default:pending" json:"status"` // pending, in_review, approved, rejected
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
