package handlers

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"zipride/utils"
)

// RefreshToken exchanges a valid refresh token for new access and refresh tokens
func RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshKey := []byte(os.Getenv("REFRESH_KEY"))
	token, err := jwt.ParseWithClaims(req.RefreshToken, &utils.Claims{}, func(t *jwt.Token) (any, error) { return refreshKey, nil })
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	claims, ok := token.Claims.(*utils.Claims)
	if !ok || claims.TokenId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh claims"})
		return
	}

	// Check token id is still valid in Redis
	if !utils.CheckToken(claims.TokenId) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token revoked"})
		return
	}

	// Rotate: delete old and issue new
	_ = utils.DeleteRefresh(claims.TokenId)

	access, err := utils.GenerateAccess(claims.UserId, claims.Email, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue access"})
		return
	}
	newRefresh, err := utils.GenerateRefresh(claims.UserId, claims.Email, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to issue refresh"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": access, "refresh_token": newRefresh})
}

// Logout revokes the provided refresh token
func Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshKey := []byte(os.Getenv("REFRESH_KEY"))
	token, err := jwt.ParseWithClaims(req.RefreshToken, &utils.Claims{}, func(t *jwt.Token) (any, error) { return refreshKey, nil })
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}
	claims, ok := token.Claims.(*utils.Claims)
	if !ok || claims.TokenId == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh claims"})
		return
	}
	_ = utils.DeleteRefresh(claims.TokenId)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}
