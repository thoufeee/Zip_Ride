package allpermission

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// verification for roles manager || staff || admin

func Verify(c *gin.Context) {

	roleVal, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "role not found"})
		return
	}

	role := roleVal.(string)

	if role == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid role"})
		return
	}

	perms, _ := c.Get("permissions")

	c.JSON(http.StatusOK, gin.H{
		"role":        role,
		"permissions": perms,
	})
}
