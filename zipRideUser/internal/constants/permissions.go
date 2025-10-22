package constants

// permission for admins

const (
	PermissionAccessAdminDash = "ACCESS_ADMINDASH"

	// user managment
	PermissionViewUsers   = "VIEW_USERS"
	PermissionAddUser     = "ADD_USER"
	PermissionEditUser    = "EDIT_USER"
	PermissionDeleteUser  = "DELETE_USER"
	PermissionBlockUser   = "BLOCK_USER"
	PermissionUnBlockUser = "UNBLOCK_USER"

	// staff managment
	PermissionViewStaffs   = "VIEW_STAFFS"
	PermissionAddStaff     = "ADD_STAFF"
	PermissionEditStaff    = "EDIT_STAFF"
	PermissionDeleteStaff  = "DELETE_STAFF"
	PermissionBlockStaff   = "BLOCK_STAFF"
	PermissionUnBlockStaff = "UNBLOCK_STAFF"

	// booking managment
	PermissionBookingManagment = "BOOKING_MANAGEMENT"
	PermissionScheduleBooking  = "SCHEDULED_BOOKINGS"

	// view ride && eranings && user analytics
	PermissionViewAnalytics = "VIEW_ANALYTICS"

	// prize pool and commission
	PermissionPrizePool = "ACCESS_PRIZEPOOL"

	// system settings (pricing rules && commisions && app settings)
	PermissionSystemSettings = "SYSTEM_SETTINGS"

	// all permissions
	ViewAllPermissions = "VIEW_ALLPERMISSIONS"
)
