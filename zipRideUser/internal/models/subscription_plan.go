package models

import "time"

// subscription plan

type SubscriptionPlan struct {
	ID                string    `gorm:"type:uuid;primaryKey" json:"id"`
	PlanName          string    `json:"planname"`
	Description       string    `json:"description"`
	DurationDays      int       `json:"duration_days"`
	Price             float64   `json:"price"`
	ComissionDiscount float64   `json:"comission_discount"`
	CreatedAt         time.Time `gorm:"autoCreateTime"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime"`
	DeletedAt         time.Time `gorm:"autoDeleteTime"`
}

// user subscription

type UserSubscription struct {
	ID        string `gorm:"type:uuid;primarykey" json:"id"`
	UserID    uint   `gorm:"not null;index"`
	PlanID    string `gorm:"type:uuid"`
	StartDate time.Time
	EndDate   time.Time
	Status    string
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoCreateTime"`
	DeletedAt time.Time `gorm:"autoCreateTime"`
}
