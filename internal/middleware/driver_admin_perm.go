package middleware

import (
	"net/http"
	"zipride/database"

	"github.com/gin-gonic/gin"
)

// RequireDriverAdminPerm checks if the driver admin has the given permission key
func RequireDriverAdminPerm(perm string) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized"})
			c.Abort()
			return
		}
		adminID, ok := val.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid admin id"})
			c.Abort()
			return
		}
		var count int64
		database.DB.Table("driver_admin_account_roles dar").
			Joins("JOIN driver_admin_role_permissions drp ON drp.role_id = dar.role_id").
			Joins("JOIN driver_admin_permissions dp ON dp.id = drp.permission_id").
			Where("dar.admin_id = ? AND dp.key = ?", adminID, perm).
			Count(&count)
		if count == 0 {
			c.JSON(http.StatusForbidden, gin.H{"err": "missing permission"})
			c.Abort()
			return
		}
		c.Next()
	}
}
