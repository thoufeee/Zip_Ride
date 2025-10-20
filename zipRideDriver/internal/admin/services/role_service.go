package adminservices

import (
	"strings"

	"gorm.io/gorm"

	"zipRideDriver/internal/models"
)

var DefaultPermissions = []string{
	"admin:dashboard:view",
	"admin:drivers:view",
	"admin:drivers:approve",
	"admin:drivers:reject",
	"admin:drivers:suspend",
	"admin:vehicles:view",
	"admin:vehicles:verify",
	"admin:vehicles:assign",
	"admin:vehicles:deactivate",
	"admin:rides:view",
	"admin:rides:cancel",
	"admin:rides:monitor",
	"admin:earnings:view",
	"admin:withdrawals:approve",
	"admin:withdrawals:reject",
	"admin:users:view",
	"admin:users:block",
	"admin:users:unblock",
	"admin:help:view",
	"admin:help:reply",
	"admin:help:close",
	"admin:admins:view",
	"admin:admins:edit",
	"admin:roles:view",
	"admin:roles:edit",
	"admin:settings:edit",
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
	roleID, err := GetUserRoleID(db, userID)
	if err != nil { return false, err }
	return HasPermission(db, roleID, permissionName)
}
