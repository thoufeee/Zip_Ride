package models

import (
	"time"

	"gorm.io/gorm"
)

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
	ID        string         `gorm:"type:uuid;primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	UserName  string         `json:"user_name"`
	UserEmail string         `json:"user_email"`
	PlanID    string         `gorm:"type:uuid"`
	PlanName  string         `json:"plan_name"`
	StartDate time.Time      `json:"start_date"`
	EndDate   time.Time      `json:"end_date"`
	Status    string         `gorm:"type:varchar(20)" json:"status"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
