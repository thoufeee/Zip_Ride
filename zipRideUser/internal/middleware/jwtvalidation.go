package middleware

import (
	"net/http"
	"os"
	"strings"
	"zipride/database"
	"zipride/internal/constants"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// jwt token check

func JwtValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		if !strings.HasPrefix(header, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
			c.Abort()
			return
		}

		tokenstr := strings.TrimPrefix(header, "Bearer ")
		jwtkey := []byte(os.Getenv("JWT_SECRET"))

		token, err := jwt.ParseWithClaims(tokenstr, &utils.Claims{}, func(t *jwt.Token) (any, error) {
			return jwtkey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid or token expired"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*utils.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid token claims"})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserId)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		// if admin || staff || manager load permissions

		if claims.Role != constants.RoleUser {
			var admin models.Admin

			if err := database.DB.Preload("Role.Permissions").Preload("Permissions").First(&admin, claims.UserId).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"err": "failed to load permissions"})
				c.Abort()
				return
			}

			perms := utils.MergePermissions(utils.PermissionToString(admin.Role.Permissions),
				utils.PermissionToString(admin.Permissions),
			)

			c.Set("permissions", perms)
			c.Set("admin", admin)
		}

		c.Next()

	}
}
