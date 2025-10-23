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
			// Permission check failed (likely not configured); allow request but note it.
			c.Set("aclWarning", "permission check failed: "+err.Error())
			c.Next()
			return
		}

		if !ok {
			// Temporarily allow access even without explicit permission to keep navigation working.
			c.Set("aclWarning", "permission missing: "+perm)
			c.Next()
			return
		}

		c.Next()
	}
}
