package seeders

import (
	"log"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
)

// seeding roles

func SeedRoles() {

	roles := []string{constants.RoleSuperAdmin, constants.RoleManager, constants.RoleStaff, constants.RoleUser}

	for _, name := range roles {
		database.DB.FirstOrCreate(&models.Role{}, models.Role{Name: name})
	}

	log.Println("role seeded successfuly")
}
