package adminmiddleware

import (
	"net/http"
	"context"
	"time"
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
)

const AdminSessionCookie = "admin_session"
const adminSessionPrefix = "admin_ssr:"

func randomSID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil { return "", err }
	return hex.EncodeToString(b), nil
}

// SSRAuthMiddleware loads admin session from Redis using SID cookie and sets admin into context.
func SSRAuthMiddleware(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		sid, err := c.Cookie(AdminSessionCookie)
		if err != nil || sid == "" {
			c.Next()
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		key := adminSessionPrefix + sid
		adminIDStr, err := rdb.Get(ctx, key).Result()
		if err != nil {
			c.Next(); return
		}
		var admin models.AdminUser
		if err := db.First(&admin, adminIDStr).Error; err == nil {
			c.Set("admin", admin)
		}
		c.Next()
	}
}

// RequireAdmin ensures an admin is logged in; otherwise redirects to login.
func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("admin"); !exists {
			c.Redirect(http.StatusSeeOther, "/admin/panel/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

// SetAdminSession creates a Redis session and sets SID cookie.
func SetAdminSession(c *gin.Context, cfg *config.Config, rdb *redis.Client, admin models.AdminUser, expiryMinutes int) error {
	sid, err := randomSID()
	if err != nil { return err }
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := rdb.Set(ctx, adminSessionPrefix+sid, admin.ID, time.Duration(expiryMinutes)*time.Minute).Err(); err != nil {
		return err
	}
	c.SetCookie(AdminSessionCookie, sid, expiryMinutes*60, "/", "", false, true)
	return nil
}

// ClearAdminSession deletes Redis session and clears cookie
func ClearAdminSession(c *gin.Context, rdb *redis.Client) {
	sid, _ := c.Cookie(AdminSessionCookie)
	if sid != "" {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		rdb.Del(ctx, adminSessionPrefix+sid)
		cancel()
	}
	c.SetCookie(AdminSessionCookie, "", -1, "/", "", false, true)
}
