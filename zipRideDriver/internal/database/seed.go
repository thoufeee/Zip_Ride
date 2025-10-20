package database

import (
	"strings"

	"zipRideDriver/internal/config"
	adminservices "zipRideDriver/internal/admin/services"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SeedDefaults(db *gorm.DB, cfg *config.Config, log *zap.Logger) error {
	adminEmail := strings.TrimSpace(cfg.DriverAdminEmail)
	adminPass := strings.TrimSpace(cfg.DriverAdminPassword)
	if adminEmail == "" || adminPass == "" {
		return nil
	}
	var admin models.AdminUser
	if err := db.Where("email = ?", adminEmail).First(&admin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hash, herr := utils.HashPassword(adminPass)
			if herr != nil {
				return herr
			}
			admin = models.AdminUser{Name: "ZipRide Admin", Email: adminEmail, PasswordHash: hash, Role: "super_admin"}
			if err := db.Create(&admin).Error; err != nil {
				return err
			}
			log.Info("seeded default admin", zap.String("email", adminEmail))
		}
	} else {
		// ensure role field on existing default admin
		if admin.Role == "" {
			_ = db.Model(&admin).Update("role", "super_admin").Error
		}
	}
	// ensure admin role exists and map to user
	var role models.Role
	if err := db.Where("name = ?", "admin").First(&role).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			role = models.Role{Name: "admin"}
			if err := db.Create(&role).Error; err != nil {
				return err
			}
		}
	}
	var ur models.UserRole
	if err := db.Where("user_id = ? AND role_id = ?", admin.ID, role.ID).First(&ur).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			ur = models.UserRole{UserID: admin.ID, RoleID: role.ID}
			if err := db.Create(&ur).Error; err != nil {
				return err
			}
		}
	}

	// Ensure default permissions exist and grant all to admin role
	if err := adminservices.EnsurePermissions(db, adminservices.DefaultPermissions); err != nil {
		return err
	}
	// grant all permissions to admin role
	for _, pname := range adminservices.DefaultPermissions {
		var perm models.Permission
		if err := db.Where("name = ?", strings.TrimSpace(pname)).First(&perm).Error; err == nil {
			var rp models.RolePermission
			if err := db.Where("role_id = ? AND permission_id = ?", role.ID, perm.ID).First(&rp).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					rp = models.RolePermission{RoleID: role.ID, PermissionID: perm.ID}
					_ = db.Create(&rp).Error
				}
			}
		}
	}
	return nil
}
