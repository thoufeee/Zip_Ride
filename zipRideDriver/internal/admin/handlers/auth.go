package adminhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"

	"zipRideDriver/internal/admin/middleware"
	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"
)

type AuthHandler struct {
	cfg *config.Config
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db, rdb: rdb, log: log}
}

func (h *AuthHandler) LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/login.html", gin.H{"title": "Admin Login"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	var admin models.AdminUser
	if err := h.db.Where("email = ?", email).First(&admin).Error; err != nil {
		c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{"error": "Invalid credentials"})
		return
	}
	if !utils.CheckPassword(admin.PasswordHash, password) {
		c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{"error": "Invalid credentials"})
		return
	}
	if err := adminmiddleware.SetAdminSession(c, h.cfg, h.rdb, admin, 60*24); err != nil {
		c.HTML(http.StatusInternalServerError, "admin/login.html", gin.H{"error": "Failed to create session"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/panel/dashboard")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	adminmiddleware.ClearAdminSession(c, h.rdb)
	c.Redirect(http.StatusSeeOther, "/admin/panel/login")
}
