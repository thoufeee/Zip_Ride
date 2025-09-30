package middleware

import (
	"net/http"
	"os"
	"strings"
	"zipride/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// auth

func Auth(roles ...string) gin.HandlerFunc {
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

		for _, role := range roles {
			if claims.Role == role {
				c.Set("user_id", claims.UserId)
				c.Set("email", claims.Email)
				c.Set("role", claims.Role)
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"err": "access denied"})
		c.Abort()

	}
}
