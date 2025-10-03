package handlers

import (
	"net/http"
	"os"
	"strings"
	"zipride/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// user logout
func UserLogout(c *gin.Context) {
	header := c.GetHeader("authorization")

	if !strings.HasPrefix(header, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "missing or invalid token"})
		return
	}

	refresh := strings.TrimPrefix(header, "Bearer ")

	token, err := jwt.ParseWithClaims(refresh, &utils.Claims{}, func(t *jwt.Token) (any, error) {
		return []byte(os.Getenv("REFRESH_KEY")), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "invalid token"})
		return
	}

	claims, ok := token.Claims.(*utils.Claims)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "unathorized"})
		return
	}

	if err := utils.DeleteRefresh(claims.TokenId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "logout failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"res": "Loged out Completed"})
}
