package constants

// permission for admins

const (

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

	// manager managment
	PermissionViewManagers   = "VIEW_MANAGERS"
	PermissionAddManager     = "ADD_MANAGER"
	PermissionEditManager    = "EDIT_MANAGER"
	PermissionDeleteManager  = "DELETE_MANAGER"
	PermissionBlockManager   = "BLOCK_MANAGER"
	PermissionUnBlockManager = "UNBLOCK_MANAGER"

	//  Reports
	PermissionViewReports = "VIEW_REPORTS"
	// view ride && eranings && user analytics
	PermissionViewAnalytics = "VIEW_ANALYTICS"

	// system settings (pricing rules && commisions && app settings)
	PermissionSystemSettings = "SYSTEM_SETTINGS"
)
