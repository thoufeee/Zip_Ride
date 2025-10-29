package models

import "gorm.io/gorm"

// web configuration

type WebConfig struct {
	gorm.Model
	SiteName          string `json:"site_name"`
	Currency          string `json:"currency"`
	CurrencySymbol    string `json:"currency_symbol"`
	PaymentGateway    string `json:"payment_gateway"`
	PaymentPublicKey  string `json:"payment_publicKey"`
	PaymentSecertKey  string `json:"payment_secretKey"`
	ContactEmail      string `json:"contact_email"`
	ContactPhone      string `json:"contact_phone"`
	ContactAddress    string `json:"contact_address"`
	MainteanceMode    bool   `json:"maintenance_mode"`
	MainteanceMessage string `json:"maintenance_message"`
}
