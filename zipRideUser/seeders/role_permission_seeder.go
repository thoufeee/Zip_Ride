package seeders

import (
	"log"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
)

// giving permission to each role automatically

func SeedRolePermissions() {

	var allPerms []models.Permission

	database.DB.Find(&allPerms)

	//  giving permission to super admin
	var SuperAdmin models.Role

	database.DB.Where("name = ?", constants.RoleSuperAdmin).First(&SuperAdmin)
	database.DB.Model(&SuperAdmin).Association("Permissions").Replace(allPerms)

	//  giving permission to manager
	var manager models.Role

	database.DB.Where("name = ?", constants.RoleManager).First(&manager)
	var managerPermission []models.Permission
	database.DB.Where("name IN ?", []string{
		constants.PermissionViewUsers, constants.PermissionViewAnalytics, constants.PermissionViewReports,
	}).Find(&managerPermission)
	database.DB.Model(&manager).Association("Permissions").Replace(managerPermission)

	// giving permission to staff
	var staff models.Role

	database.DB.Where("name = ?", constants.RoleStaff).First(&staff)
	var staffPermission []models.Permission
	database.DB.Where("name IN ?", []string{
		constants.PermissionViewUsers,
	}).Find(&staffPermission)
	database.DB.Model(&staff).Association("Permissions").Replace(staffPermission)

	log.Println("permission assegned to roles successfuly")

}
