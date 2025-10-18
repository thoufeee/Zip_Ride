package models

import "time"

type Ride struct {
	ID        uint      `gorm:"primaryKey"`
	DriverID  uint      `gorm:"index"`
	Pickup    string    `gorm:"size:255;not null"`
	Dropoff   string    `gorm:"size:255;not null"`
	Fare      float64   `gorm:"default:0"`
	Distance  float64   `gorm:"default:0"`
	Status    string    `gorm:"size:30;default:'pending';index"`
	StartedAt *time.Time
	EndedAt   *time.Time
	CreatedAt time.Time
}

type Earning struct {
	ID        uint      `gorm:"primaryKey"`
	DriverID  uint      `gorm:"index;not null"`
	RideID    uint      `gorm:"index"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
