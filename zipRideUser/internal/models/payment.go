package models

import (
	"time"

	"github.com/google/uuid"
)

// payment

type Payment struct {
	ID              uuid.UUID `gorm:"type:uuid;primaryKey;" json:"id"`
	BookingID       uint      `json:"booking_id,omitempty"`
	UserID          uint      `json:"user_id"`
	DriverID        uint      `json:"driver_id,omitempty"`
	SubScribitionID string    `json:"subscribition_id"`

	TotalAmount  float64 `json:"total_amount"`
	DriverAmount float64 `json:"driver_amount"`
	Commission   float64 `json:"commission"`
	Discount     float64 `json:"discount,omitempty"`

	Currency       string `json:"currency"`
	PaymentGateway string `json:"payment_gateway"`
	PaymentMethod  string `json:"payment_method"`
	TransactionID  string `json:"transaction_id"`
	Status         string `json:"status"`
	PaymentType    string `json:"paymenttype"`

	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`

	PaymentDate  time.Time  `json:"payment_date"`
	RefundAmount float64    `json:"refund_amount"`
	RefundStatus string     `json:"refund_status"`
	RefundDate   *time.Time `json:"refund_date,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
