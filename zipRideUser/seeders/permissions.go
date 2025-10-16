package seeders

import (
	"log"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
)

// storing permissions

func SeedPermisiions() {
	permission := []models.Permission{

		//  user managmnet
		{Name: constants.PermissionViewUsers},
		{Name: constants.PermissionAddUser},
		{Name: constants.PermissionEditUser},
		{Name: constants.PermissionDeleteUser},
		{Name: constants.PermissionBlockUser},
		{Name: constants.PermissionUnBlockUser},

		// staff managment
		{Name: constants.PermissionViewStaffs},
		{Name: constants.PermissionAddStaff},
		{Name: constants.PermissionEditStaff},
		{Name: constants.PermissionDeleteStaff},
		{Name: constants.PermissionBlockStaff},
		{Name: constants.PermissionUnBlockStaff},

		// manager managment
		{Name: constants.PermissionViewManagers},
		{Name: constants.PermissionEditManager},
		{Name: constants.PermissionDeleteManager},
		{Name: constants.PermissionBlockManager},
		{Name: constants.PermissionUnBlockManager},

		// booking managment
		{Name: constants.PermissionScheduleBooking},

		// reports
		{Name: constants.PermissionViewReports},
		{Name: constants.PermissionViewAnalytics},

		// system settings
		{Name: constants.PermissionSystemSettings},
	}

	for _, p := range permission {
		database.DB.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}

	log.Println("permission seeded successfuly")
}
