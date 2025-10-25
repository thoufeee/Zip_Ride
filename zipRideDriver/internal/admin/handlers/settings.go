package adminhandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"
)

type SettingsHandler struct{ db *gorm.DB; log *zap.Logger }

func NewSettingsHandler(db *gorm.DB, log *zap.Logger) *SettingsHandler { return &SettingsHandler{db: db, log: log} }

// GET /admin/panel/settings
func (h *SettingsHandler) SettingsPage(c *gin.Context) {
	adm, _ := c.Get("admin")
	admin := adm.(models.AdminUser)
	c.HTML(http.StatusOK, "admin/settings/index.html", gin.H{"Admin": admin})
}

// POST /admin/panel/settings/password
func (h *SettingsHandler) UpdatePassword(c *gin.Context) {
	adm, _ := c.Get("admin")
	admin := adm.(models.AdminUser)
	cur := c.PostForm("current_password")
	newp := strings.TrimSpace(c.PostForm("new_password"))
	conf := strings.TrimSpace(c.PostForm("confirm_password"))
	if newp == "" || newp != conf {
		c.HTML(http.StatusBadRequest, "admin/settings/index.html", gin.H{"Admin": admin, "Error": "Passwords do not match"})
		return
	}
	if !utils.CheckPassword(admin.PasswordHash, cur) {
		c.HTML(http.StatusBadRequest, "admin/settings/index.html", gin.H{"Admin": admin, "Error": "Current password incorrect"})
		return
	}
	hash, err := utils.HashPassword(newp)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin/settings/index.html", gin.H{"Admin": admin, "Error": "Failed to hash password"})
		return
	}
	if err := h.db.Model(&models.AdminUser{}).Where("id = ?", admin.ID).Update("password_hash", hash).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "admin/settings/index.html", gin.H{"Admin": admin, "Error": "Failed to update password"})
		return
	}
	h.log.Info("admin password updated", zap.Uint("admin_id", admin.ID))
	c.HTML(http.StatusOK, "admin/settings/index.html", gin.H{"Admin": admin, "Success": "Password updated"})
}
