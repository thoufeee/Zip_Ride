package models

import "time"

type ChatSession struct {
	ID        uint      `gorm:"primaryKey"`
	DriverID  uint      `gorm:"index;not null"`
	Status    string    `gorm:"size:30;default:'open';index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type ChatMessage struct {
	ID         uint      `gorm:"primaryKey"`
	SessionID  uint      `gorm:"index;not null"`
	Sender     string    `gorm:"size:20;not null"`
	Message    string    `gorm:"type:text;not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
