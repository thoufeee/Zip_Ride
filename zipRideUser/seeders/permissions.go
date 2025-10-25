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

		{Name: constants.PermissionAccessAdminDash},

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

		// booking managment
		{Name: constants.PermissionScheduleBooking},
		{Name: constants.PermissionBookingManagment},

		// reports
		{Name: constants.PermissionViewAnalytics},

		// system settings
		{Name: constants.PermissionSystemSettings},

		// permissions
		{Name: constants.ViewAllPermissions},

		// prize pool
		{Name: constants.PermissionPrizePool},

		// subscription
		{Name: constants.PermissionSubscription},
	}

	for _, p := range permission {
		database.DB.FirstOrCreate(&p, models.Permission{Name: p.Name})
	}

	log.Println("permission seeded successfuly")
}
