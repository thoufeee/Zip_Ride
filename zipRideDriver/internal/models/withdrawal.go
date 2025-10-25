package models

import "time"

type Withdrawal struct {
	ID        uint      `gorm:"primaryKey"`
	DriverID  uint      `gorm:"index;not null"`
	Amount    float64   `gorm:"not null"`
	Status    string    `gorm:"size:30;default:'pending';index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
