package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func permsCacheKey(userID uint) string { return "acl:admin:perms:" + strconv.Itoa(int(userID)) }

func RequirePermissions(db *gorm.DB, rdb *redis.Client, log *zap.Logger, required ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		subjAny, _ := c.Get("subject")
		subj, _ := subjAny.(string)
		if subj != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		uidAny, _ := c.Get("uid")
		uid, ok := uidAny.(uint)
		if !ok || uid == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid user"})
			return
		}

		ctx := context.Background()
		key := permsCacheKey(uid)
		cached, err := rdb.Get(ctx, key).Result()
		has := map[string]bool{}
		if err == nil && cached != "" {
			for _, p := range strings.Split(cached, ",") {
				p = strings.TrimSpace(p)
				if p != "" {
					has[p] = true
				}
			}
		} else {
			var rows []struct{ Name string }
			q := `SELECT DISTINCT p.name AS name
        FROM user_roles ur
        JOIN roles r ON r.id = ur.role_id
        JOIN role_permissions rp ON rp.role_id = r.id
        JOIN permissions p ON p.id = rp.permission_id
        WHERE ur.user_id = ?`
			if err := db.Raw(q, uid).Scan(&rows).Error; err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "failed to load permissions"})
				return
			}
			list := make([]string, 0, len(rows))
			for _, r := range rows {
				has[r.Name] = true
				list = append(list, r.Name)
			}
			_ = rdb.Set(ctx, key, strings.Join(list, ","), 5*time.Minute).Err()
		}

		for _, need := range required {
			if !has[need] {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
				return
			}
		}
		c.Next()
	}
}
