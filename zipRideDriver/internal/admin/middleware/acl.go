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
		// Super Admin bypass
		if strings.EqualFold(admin.Role, "super_admin") {
			c.Next(); return
		}
		ok, err := adminservices.HasPermissionForUser(db, admin.ID, perm)
		if err != nil || !ok {
			c.String(http.StatusForbidden, "Access Denied: insufficient privileges")
			c.Abort()
			return
		}
		c.Next()
	}
}
