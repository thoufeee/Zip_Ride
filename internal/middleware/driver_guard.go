package middleware

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// getting user id from token

func GetUserID(c *gin.Context) uint {
	if userID, exists := c.Get("user_id"); exists {
		if uid, ok := userID.(uint); ok {
			return uid
		}
	}
	return 0
}

// RequireApprovedDriver ensures the authenticated user is an approved, phone-verified driver
func RequireApprovedDriver() gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "unauthorized"})
			c.Abort()
			return
		}
		uid, ok := val.(uint)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid user id"})
			c.Abort()
			return
		}
		var d models.Driver
		if err := database.DB.First(&d, uid).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "driver not found"})
			c.Abort()
			return
		}
		if !d.PhoneVerified || d.Status != "approved" {
			c.JSON(http.StatusForbidden, gin.H{"err": "driver not approved or phone not verified"})
			c.Abort()
			return
		}
		c.Next()
	}
}
