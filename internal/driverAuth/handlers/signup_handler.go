package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"zipride/internal/driverAuth/services"
)

var authSvc = services.NewAuthService()

// EmailSignupHandler handles signup with email + password + phone
func EmailSignupHandler(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
		Phone     string `json:"phone" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Trim inputs to avoid whitespace issues
	req.FirstName = strings.TrimSpace(req.FirstName)
	req.LastName = strings.TrimSpace(req.LastName)
	req.Email = strings.TrimSpace(req.Email)
	req.Phone = strings.TrimSpace(req.Phone)

	driver, err := authSvc.SignupWithEmail(context.Background(), req.FirstName, req.LastName, req.Email, req.Password, req.Phone)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	// Return minimal response
	c.JSON(http.StatusCreated, gin.H{
		"driver_id":      driver.ID,
		"phone_verified": driver.IsVerified, // <-- fix here
		"next_step":      "verify_phone",
		"created_at":     driver.CreatedAt.Format(time.RFC3339),
	})
}
