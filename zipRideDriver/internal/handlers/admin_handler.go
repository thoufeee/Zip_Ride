package handlers

import (
	"net/http"
	"strings"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/middleware"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/services"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// RegisterAdminRoutes wires admin-related endpoints (to be migrated from existing routes).
func RegisterAdminRoutes(r *gin.Engine, cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) {
	g := r.Group("/admin")

	// Admin login
	g.POST("/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		var admin models.AdminUser
		if err := db.Where("email = ?", strings.TrimSpace(req.Email)).First(&admin).Error; err != nil {
			utils.Error(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		if !utils.CheckPassword(admin.PasswordHash, req.Password) {
			utils.Error(c, http.StatusUnauthorized, "invalid credentials")
			return
		}
		access, err := services.GenerateToken(admin.ID, admin.Email, "admin", []string{"admin"}, cfg.JWTSecret, cfg.AccessTokenExpiry)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to issue token")
			return
		}
		refresh, err := services.GenerateToken(admin.ID, admin.Email, "admin_refresh", []string{"admin"}, cfg.JWTSecret, cfg.RefreshTokenExpiry)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to issue token")
			return
		}
		utils.Ok(c, "login success", gin.H{"access_token": access, "refresh_token": refresh})
	})

	auth := g.Group("")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	auth.Use(middleware.AdminOnly())

	auth.GET("/drivers", func(c *gin.Context) {
		var list []models.Driver
		if err := db.Order("id desc").Limit(200).Find(&list).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to list drivers")
			return
		}
		utils.Ok(c, "drivers", list)
	})

	// Approve/Reject
	auth.POST("/driver/:id/approve", func(c *gin.Context) {
		id := strings.TrimSpace(c.Param("id"))
		var body struct {
			Action string `json:"action" binding:"required,oneof=approve reject"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		status := "Approved"
		if body.Action == "reject" {
			status = "Rejected"
		}
		if err := db.Model(&models.Driver{}).Where("id = ?", id).Update("status", status).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to update status")
			return
		}
		utils.Ok(c, "updated", gin.H{"status": status})
	})

	auth.POST("/driver/:id/ban", func(c *gin.Context) {
		id := strings.TrimSpace(c.Param("id"))
		if err := db.Model(&models.Driver{}).Where("id = ?", id).Update("status", "Banned").Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to ban driver")
			return
		}
		utils.Ok(c, "banned", nil)
	})

	auth.GET("/dashboard", func(c *gin.Context) {
		var totals struct {
			Drivers  int64
			Pending  int64
			Earnings float64
		}
		_ = db.Model(&models.Driver{}).Count(&totals.Drivers).Error
		_ = db.Model(&models.Driver{}).Where("status = ?", "Pending").Count(&totals.Pending).Error
		_ = db.Model(&models.Earning{}).Select("COALESCE(SUM(amount),0)").Scan(&totals.Earnings).Error
		var lastDrivers []models.Driver
		_ = db.Order("created_at desc").Limit(10).Find(&lastDrivers).Error
		utils.Ok(c, "dashboard", gin.H{
			"totals":        totals,
			"recent_drivers": lastDrivers,
			"ts":            time.Now().UTC(),
		})
	})
}
