package adminmiddleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zipRideDriver/internal/admin/services"
	"zipRideDriver/internal/models"
)

// ACLMiddleware enforces a required permission for SSR admin routes.
func ACLMiddleware(db *gorm.DB, requiredPermission string) gin.HandlerFunc {
	perm := strings.TrimSpace(requiredPermission)
	return func(c *gin.Context) {
		adm, exists := c.Get("admin")
		if !exists {
			c.Redirect(http.StatusSeeOther, "/admin/panel/login")
			c.Abort()
			return
		}
		admin := adm.(models.AdminUser)
		
		// Super Admin bypass - check both role field and super_admin variations
		if strings.EqualFold(admin.Role, "super_admin") || 
		   strings.EqualFold(admin.Role, "superadmin") || 
		   strings.EqualFold(admin.Role, "admin") {
			c.Next()
			return
		}
		
		// For development/testing - allow access if no specific permission system is set up
		if perm == "" {
			c.Next()
			return
		}
		
		// Check permissions through role system
		ok, err := adminservices.HasPermissionForUser(db, admin.ID, perm)
		if err != nil {
			// If permission checking fails, allow super_admin users through
			if strings.EqualFold(admin.Role, "super_admin") {
				c.Next()
				return
			}
			c.String(http.StatusForbidden, "Access Denied: permission check failed")
			c.Abort()
			return
		}
		
		if !ok {
			c.String(http.StatusForbidden, "Access Denied: insufficient privileges")
			c.Abort()
			return
		}
		
		c.Next()
	}
}
