package database

import (
	"strings"

	adminservices "zipRideDriver/internal/admin/services"
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
	// Normalize any incorrectly stored role values
	_ = db.Model(&models.AdminUser{}).Where("role = ?", "superadmin").Update("role", "super_admin").Error

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

	// Create a test admin user for simple cookie-based login
	var testAdmin models.AdminUser
	if err := db.Where("email = ?", "admin@zipride.com").First(&testAdmin).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			hash, herr := utils.HashPassword("admin123")
			if herr != nil {
				log.Error("failed to hash test admin password", zap.Error(herr))
			} else {
				testAdmin = models.AdminUser{Name: "Test Admin", Email: "admin@zipride.com", PasswordHash: hash, Role: "super_admin"}
				if err := db.Create(&testAdmin).Error; err != nil {
					log.Error("failed to create test admin", zap.Error(err))
				} else {
					log.Info("test admin created", zap.String("email", "admin@zipride.com"), zap.String("password", "admin123"))
				}
			}
		}
	}

	// Create sample riders for testing user management
	sampleRiders := []models.Rider{
		{Name: "Alice Johnson", Email: "alice.johnson@example.com", Phone: "+1234567800", IsBlocked: false},
		{Name: "Bob Smith", Email: "bob.smith@example.com", Phone: "+1234567801", IsBlocked: false},
		{Name: "Carol Davis", Email: "carol.davis@example.com", Phone: "+1234567802", IsBlocked: true},
		{Name: "David Wilson", Email: "david.wilson@example.com", Phone: "+1234567803", IsBlocked: false},
		{Name: "Eve Brown", Email: "eve.brown@example.com", Phone: "+1234567804", IsBlocked: false},
	}

	for _, rider := range sampleRiders {
		var existingRider models.Rider
		if err := db.Where("email = ?", rider.Email).First(&existingRider).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				if err := db.Create(&rider).Error; err != nil {
					log.Error("failed to create sample rider", zap.Error(err), zap.String("email", rider.Email))
				}
			}
		}
	}

	return nil
}
