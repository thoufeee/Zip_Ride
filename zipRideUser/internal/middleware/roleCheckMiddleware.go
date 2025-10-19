package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// rolechecks
func RoleCheck(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleface, exist := c.Get("role")

		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "role not found"})
			c.Abort()
			return
		}

		role, ok := roleface.(string)

		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"err": "invalid role format"})
			c.Abort()
			return
		}

		for _, r := range roles {
			if strings.EqualFold(role, r) {
				c.Next()
				return
			}

			fmt.Println("Role from token/context:", role)
			fmt.Println("Allowed roles:", roles)

		}

		c.JSON(http.StatusForbidden, gin.H{"err": "access deniedddddddddddddddd"})
		c.Abort()
	}
}
