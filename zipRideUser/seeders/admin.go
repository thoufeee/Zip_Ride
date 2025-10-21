package seeders

import (
	"encoding/json"
	"log"
	"os"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"gorm.io/datatypes"
)

// adding admin

func SeedAdmin() {

	adminEmail := os.Getenv("EMAIL")
	adminPassword := os.Getenv("PASSWORD")

	var admin models.Admin

	if err := database.DB.Where("email = ?", adminEmail).First(&admin).Error; err == nil {
		log.Println("admin already exist")
		return
	}

	var permissions []models.Permission
	if err := database.DB.Find(&permissions).Error; err != nil {
		log.Fatal("failed to fetch permissions")
		return
	}

	var perms []string
	for _, p := range permissions {
		perms = append(perms, p.Name)
	}

	permsjson, err := json.Marshal(perms)

	if err != nil {
		log.Fatal("failed to marshal permissions")
	}

	hashpass, err := utils.GenerateHash(adminPassword)
	if err != nil {
		log.Fatal("failed to hash adminpass")
	}

	newadmin := &models.Admin{
		Email:       adminEmail,
		Password:    hashpass,
		Permissions: datatypes.JSON(permsjson),
		Role:        constants.RoleSuperAdmin,
	}

	if err := database.DB.Create(newadmin).Error; err != nil {
		log.Fatal("failed to create admin account")
	}

	log.Println("Admin created Successfuly")
}
