package database

import (
	"strings"

	"zipRideDriver/internal/config"
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
			admin = models.AdminUser{Name: "ZipRide Admin", Email: adminEmail, PasswordHash: hash}
			if err := db.Create(&admin).Error; err != nil {
				return err
			}
			log.Info("seeded default admin", zap.String("email", adminEmail))
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
	return nil
}
