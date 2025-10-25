package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"
)

type PublicHandler struct {
	cfg *config.Config
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewPublicHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *PublicHandler {
	return &PublicHandler{
		cfg: cfg,
		db:  db,
		rdb: rdb,
		log: log,
	}
}

// RegisterDriver handles driver registration (legacy endpoint)
// This is a simplified version - use /api/driver/register for full registration
func (h *PublicHandler) RegisterDriver(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Phone    string `json:"phone" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if driver exists
	var existingDriver models.Driver
	if err := h.db.Where("email = ? OR phone = ?", req.Email, req.Phone).First(&existingDriver).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Driver with this email or phone already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		h.log.Error("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}

	// Create driver
	driver := models.Driver{
		Name:         req.Name,
		Email:        strings.ToLower(req.Email),
		Phone:        req.Phone,
		PasswordHash: hashedPassword,
		Status:       "Pending", // Requires admin approval
		IsVerified:   false,
		IsOnline:     false,
		Rating:       5.0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := h.db.Create(&driver).Error; err != nil {
		h.log.Error("Failed to create driver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete registration"})
		return
	}

	h.log.Info("Driver registered",
		zap.Uint("driver_id", driver.ID),
		zap.String("email", driver.Email))

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Registration successful. Awaiting admin approval.",
		"driver_id": driver.ID,
	})
}

// LoginDriver handles driver login
func (h *PublicHandler) LoginDriver(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find driver
	var driver models.Driver
	if err := h.db.Where("email = ?", strings.ToLower(req.Email)).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			return
		}
		h.log.Error("Database error during login", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Login failed"})
		return
	}

	// Check if driver is approved
	if driver.Status != "Approved" && driver.Status != "Active" {
		c.JSON(http.StatusForbidden, gin.H{
			"error":   "Account not approved",
			"status":  driver.Status,
			"message": "Your account is pending admin approval",
		})
		return
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, driver.PasswordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(driver.ID, driver.Email, h.cfg.JWTSecret)
	if err != nil {
		h.log.Error("Failed to generate token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Update online status
	h.db.Model(&driver).Update("is_online", true)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"driver": gin.H{
			"id":          driver.ID,
			"name":        driver.Name,
			"email":       driver.Email,
			"phone":       driver.Phone,
			"status":      driver.Status,
			"is_verified": driver.IsVerified,
		},
	})
}

// Health check endpoint
func (h *PublicHandler) Health(c *gin.Context) {
	// Check database connection
	var count int64
	dbHealthy := h.db.Raw("SELECT 1").Count(&count).Error == nil

	// Check Redis connection
	ctx := c.Request.Context()
	_, redisErr := h.rdb.Ping(ctx).Result()
	redisHealthy := redisErr == nil

	status := "healthy"
	statusCode := http.StatusOK

	if !dbHealthy || !redisHealthy {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":    status,
		"timestamp": time.Now().Unix(),
		"services": gin.H{
			"database": dbHealthy,
			"redis":    redisHealthy,
		},
	})
}
