package models

import "time"

type HelpTicket struct {
	ID          uint      `gorm:"primaryKey"`
	DriverID    uint      `gorm:"index;not null"`
	Subject     string    `gorm:"size:120;not null"`
	Description string    `gorm:"type:text;not null"`
	Status      string    `gorm:"size:30;default:'open';index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}
