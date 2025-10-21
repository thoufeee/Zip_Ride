package models

import "time"

type Driver struct {
	ID            uint       `gorm:"primaryKey"`
	Name          string     `gorm:"size:100;not null"`
	Email         string     `gorm:"size:120;uniqueIndex"`
	Phone         string     `gorm:"size:20;uniqueIndex;index"`
	PasswordHash  string     `gorm:"size:255"`
	LicenseNumber string     `gorm:"size:50"`
	VehicleNumber string     `gorm:"size:50"`
	VehicleModel  string     `gorm:"size:100"`
	VehicleType   string     `gorm:"size:50"`
	Status        string     `gorm:"size:30;default:'Pending';index"`
	IsOnline      bool       `gorm:"default:false"`
	IsVerified    bool       `gorm:"default:false"`
	Rating        float64    `gorm:"default:5"`
	VerifiedAt    *time.Time `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time

	Documents []DriverDocument `gorm:"foreignKey:DriverID;constraint:OnDelete:CASCADE"`
	Rides     []Ride           `gorm:"foreignKey:DriverID"`
	Earnings  []Earning        `gorm:"foreignKey:DriverID"`
}

type DriverDocument struct {
	ID         uint       `gorm:"primaryKey"`
	DriverID   uint       `gorm:"index;not null"`
	DocType    string     `gorm:"size:50;not null"`
	DocURL     string     `gorm:"size:500;not null"`
	Verified   bool       `gorm:"default:false"`
	UploadedAt time.Time  `gorm:"autoCreateTime"`
	VerifiedAt *time.Time `gorm:"index"`
}
