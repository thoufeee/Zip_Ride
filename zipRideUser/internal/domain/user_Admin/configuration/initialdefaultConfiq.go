package configuration

import (
	"log"
	"zipride/database"
	"zipride/internal/models"
)

// default configuration

func InitialDefaultConfig() {
	defaultConfig := &models.WebConfig{
		SiteName:          "Zip Ride",
		Currency:          "INR",
		CurrencySymbol:    "₹",
		PaymentGateway:    "000",
		PaymentPublicKey:  "000",
		PaymentSecertKey:  "000",
		ContactEmail:      "zipRide@gmail.com",
		ContactPhone:      "+910000000000",
		ContactAddress:    "PIPLINE,Kochi",
		MainteanceMode:    false,
		MainteanceMessage: "under Maintainenece",
	}

	if err := database.DB.FirstOrCreate(&defaultConfig).Error; err != nil {
		log.Println("failed to create configuration")
		return
	}
}
