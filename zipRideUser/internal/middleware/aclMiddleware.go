package middleware

import (
	"net/http"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// required permission checks
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {

		permi, exist := c.Get("permissions")

		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "permissions not found"})
			c.Abort()
			return
		}

		permissions, ok := permi.([]string)

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "invalid permission format"})
			c.Abort()
			return
		}

		if !utils.CheckPermission(permissions, permission) {
			c.JSON(http.StatusForbidden, gin.H{"err": "access denied"})
			c.Abort()
			return
		}

		c.Next()
	}
}
