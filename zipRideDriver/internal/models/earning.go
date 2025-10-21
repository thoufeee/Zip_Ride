package models

import "time"

type Earning struct {
	ID        uint      `gorm:"primaryKey"`
	DriverID  uint      `gorm:"index;not null"`
	RideID    uint      `gorm:"index"`
	Amount    float64   `gorm:"not null"`
	Type      string    `gorm:"size:50"` // ride_fare, tip, bonus
	Status    string    `gorm:"size:30;default:'pending'"` // pending, paid, withdrawn
	CreatedAt time.Time
	UpdatedAt time.Time
	
	Driver Driver `gorm:"foreignKey:DriverID"`
	Ride   *Ride  `gorm:"foreignKey:RideID"`
}
