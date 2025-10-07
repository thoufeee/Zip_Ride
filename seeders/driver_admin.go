package seeders

import (
	"log"
	"os"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"
)

// SeedDriverAdmin seeds driver admin permissions, roles and a default admin account
func SeedDriverAdmin() {
	perms := []models.DriverAdminPermission{
		{Key: constants.PermissionDriverView, Description: "View drivers and documents"},
		{Key: constants.PermissionDriverApprove, Description: "Approve drivers"},
		{Key: constants.PermissionDriverReject, Description: "Reject drivers"},
		{Key: constants.PermissionDriverSuspend, Description: "Suspend drivers"},
	}
	for _, p := range perms {
		database.DB.Where(models.DriverAdminPermission{Key: p.Key}).FirstOrCreate(&p)
	}

	role := models.DriverAdminRole{Name: "driver_admin"}
	database.DB.Where(models.DriverAdminRole{Name: role.Name}).FirstOrCreate(&role)

	var permRows []models.DriverAdminPermission
	database.DB.Find(&permRows)
	for _, p := range permRows {
		database.DB.Where(models.DriverAdminRolePermission{RoleID: role.ID, PermissionID: p.ID}).FirstOrCreate(&models.DriverAdminRolePermission{RoleID: role.ID, PermissionID: p.ID})
	}

	// seed a default driver admin account
	dAdminEmail := os.Getenv("DRIVER_ADMIN_EMAIL")
	if dAdminEmail == "" {
		dAdminEmail = "driver.admin@example.com"
	}
	var dadmin models.DriverAdmin
	database.DB.Where(models.DriverAdmin{Email: dAdminEmail}).First(&dadmin)
	if dadmin.ID == 0 {
		pwdStr := os.Getenv("DRIVER_ADMIN_PASSWORD")
		if pwdStr == "" {
			pwdStr = "Admin@123"
		}
		pwd, _ := utils.GenerateHash(pwdStr)
		dadmin = models.DriverAdmin{
			Email:        dAdminEmail,
			PasswordHash: pwd,
			FirstName:    "Driver",
			LastName:     "Admin",
			Name:         "Driver Admin",
			IsActive:     true,
		}
		if err := database.DB.Create(&dadmin).Error; err != nil {
			log.Println("failed to seed driver admin:", err)
		} else {
			log.Println("Driver admin created/ensured")
		}
	} else {
		// Ensure password field migration (if older seeds used Password)
		if dadmin.PasswordHash == "" && dadmin.Password != "" {
			// best-effort: move hashed value from old Password field
			dadmin.PasswordHash = dadmin.Password
			dadmin.Password = ""
			dadmin.IsActive = true
			database.DB.Save(&dadmin)
		}
	}
	database.DB.Where(models.DriverAdminAccountRole{AdminID: dadmin.ID, RoleID: role.ID}).FirstOrCreate(&models.DriverAdminAccountRole{AdminID: dadmin.ID, RoleID: role.ID})
}
