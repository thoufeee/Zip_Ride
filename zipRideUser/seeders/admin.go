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

	var admin models.Admin

	if err := database.DB.Where("email = ?", adminEmail).First(&admin).Error; err == nil {
		log.Println("admin already exist")
		return
	}

	// super admin role
	var supreRole models.Role
	if err := database.DB.Where("name = ?", constants.RoleSuperAdmin).First(&supreRole).Error; err != nil {
		log.Fatal("super admin role not found")
		return
	}

	// permissions
	var permissions []models.Permission
	if err := database.DB.Find(&permissions).Error; err != nil {
		log.Fatal("failed to fetch permissions")
		return
	}

	hashpass, err := utils.GenerateHash(adminPassword)
	if err != nil {
		log.Fatal("failed to hash adminpass")
	}

	newadmin := &models.Admin{
		Email:       adminEmail,
		Password:    hashpass,
		Permissions: permissions,
		RoleID:      supreRole.ID,
	}

	if err := database.DB.Create(&newadmin).Error; err != nil {
		log.Fatal("failed to create admin account")
	}

	log.Println("Admin created Successfuly")
}
