package adminservices

import (
	"strings"

	"gorm.io/gorm"

	"zipRideDriver/internal/models"
)

var DefaultPermissions = []string{
	// Dashboard & Analytics
	"admin:dashboard:view",
	"admin:analytics:view",
	
	// Driver Management
	"admin:drivers:view",
	"admin:drivers:create",
	"admin:drivers:edit",
	"admin:drivers:delete",
	"admin:drivers:approve",
	"admin:drivers:reject",
	"admin:drivers:suspend",
	"admin:drivers:block",
	"admin:drivers:unblock",
	"admin:drivers:view_pending",
	"admin:drivers:view_documents",
	"admin:drivers:verify_documents",
	
	// Vehicle Management
	"admin:vehicles:view",
	"admin:vehicles:create",
	"admin:vehicles:edit",
	"admin:vehicles:delete",
	"admin:vehicles:verify",
	"admin:vehicles:assign",
	"admin:vehicles:deactivate",
	
	// Ride Management
	"admin:rides:view",
	"admin:rides:cancel",
	"admin:rides:monitor",
	"admin:rides:edit",
	"admin:rides:view_details",
	
	// Earnings & Financial
	"admin:earnings:view",
	"admin:earnings:export",
	"admin:earnings:process_payouts",
	"admin:withdrawals:view",
	"admin:withdrawals:approve",
	"admin:withdrawals:reject",
	"admin:withdrawals:process",
	
	// User Management
	"admin:users:view",
	"admin:users:create",
	"admin:users:edit",
	"admin:users:delete",
	"admin:users:block",
	"admin:users:unblock",
	
	// Complaints & Support
	"admin:complaints:view",
	"admin:complaints:resolve",
	"admin:complaints:assign",
	"admin:help:view",
	"admin:help:reply",
	"admin:help:close",
	"admin:help:escalate",
	
	// Schedule Management
	"admin:schedule:view",
	"admin:schedule:edit",
	"admin:schedule:manage_shifts",
	
	// Admin & Role Management
	"admin:admins:view",
	"admin:admins:create",
	"admin:admins:edit",
	"admin:admins:delete",
	"admin:roles:view",
	"admin:roles:create",
	"admin:roles:edit",
	"admin:roles:delete",
	"admin:permissions:manage",
	
	// System Settings
	"admin:settings:view",
	"admin:settings:edit",
	"admin:settings:manage_config",
	
	// Reports
	"admin:reports:view",
	"admin:reports:generate",
	"admin:reports:export",
}

// Predefined roles with their typical permissions
var PredefinedRoles = map[string][]string{
	"super_admin": []string{"*"}, // All permissions
	"fleet_manager": []string{
		"admin:dashboard:view",
		"admin:drivers:*",
		"admin:vehicles:*",
		"admin:rides:*",
		"admin:schedule:*",
	},
	"finance_admin": []string{
		"admin:dashboard:view",
		"admin:earnings:*",
		"admin:withdrawals:*",
		"admin:reports:view",
		"admin:reports:generate",
		"admin:reports:export",
	},
	"support_staff": []string{
		"admin:dashboard:view",
		"admin:drivers:view",
		"admin:rides:view",
		"admin:users:view",
		"admin:complaints:*",
		"admin:help:*",
	},
	"viewer": []string{
		"admin:dashboard:view",
		"admin:drivers:view",
		"admin:vehicles:view",
		"admin:rides:view",
		"admin:earnings:view",
		"admin:users:view",
		"admin:reports:view",
	},
}

func EnsurePermissions(db *gorm.DB, names []string) error {
	for _, n := range names {
		name := strings.TrimSpace(n)
		if name == "" { continue }
		var p models.Permission
		if err := db.Where("name = ?", name).First(&p).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				p = models.Permission{Name: name}
				if err := db.Create(&p).Error; err != nil { return err }
			}
		}
	}
	return nil
}

func GetUserRoleID(db *gorm.DB, userID uint) (uint, error) {
	var ur models.UserRole
	if err := db.Where("user_id = ?", userID).First(&ur).Error; err != nil { return 0, err }
	return ur.RoleID, nil
}

func HasPermission(db *gorm.DB, roleID uint, permissionName string) (bool, error) {
	var perm models.Permission
	if err := db.Where("name = ?", strings.TrimSpace(permissionName)).First(&perm).Error; err != nil { return false, err }
	var rp models.RolePermission
	if err := db.Where("role_id = ? AND permission_id = ?", roleID, perm.ID).First(&rp).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return false, nil }
		return false, err
	}
	return true, nil
}

func HasPermissionForUser(db *gorm.DB, userID uint, permissionName string) (bool, error) {
	// First check if user has a direct role assignment
	roleID, err := GetUserRoleID(db, userID)
	if err != nil {
		// If no UserRole found, check if the user is a super_admin via AdminUser.Role field
		var admin models.AdminUser
		if err := db.First(&admin, userID).Error; err != nil {
			return false, err
		}
		// Super admins have all permissions
		if strings.EqualFold(admin.Role, "super_admin") || 
		   strings.EqualFold(admin.Role, "superadmin") ||
		   strings.EqualFold(admin.Role, "admin") {
			return true, nil
		}
		return false, err
	}
	return HasPermission(db, roleID, permissionName)
}
