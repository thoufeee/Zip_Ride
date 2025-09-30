package seeders

import (
	"log"
	"os"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"
)

// adding admin

func SeedAdmin() {

	adminEmail := os.Getenv("EMAIL")
	adminPassword := os.Getenv("PASSWORD")

	var admin models.User

	if err := database.DB.Where("email = ?", adminEmail).First(&admin).Error; err == nil {
		log.Println("admin already exist")
		return
	}

	hashpass, err := utils.GenerateHash(adminPassword)
	if err != nil {
		log.Fatal("failed to hash adminpass")
	}

	newadmin := &models.User{
		Email:    adminEmail,
		Password: hashpass,
		Role:     constants.RoleAdmin,
	}

	if err := database.DB.Create(&newadmin).Error; err != nil {
		log.Fatal("failed to create admin account")
	}

	log.Println("Admin created Successfuly")
}
