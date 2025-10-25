package adminhandlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"
)

type AdminsHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewAdminsHandler(db *gorm.DB, log *zap.Logger) *AdminsHandler {
	return &AdminsHandler{db: db, log: log}
}

type adminView struct {
	Admin models.AdminUser
	Role  string
}

func (h *AdminsHandler) AdminsPage(c *gin.Context) {
	var admins []models.AdminUser
	_ = h.db.Order("email asc").Find(&admins).Error
	var roles []models.Role
	_ = h.db.Order("name asc").Find(&roles).Error
	roleMap := map[uint]string{}
	for _, r := range roles { roleMap[r.ID] = r.Name }
	var urs []models.UserRole
	_ = h.db.Find(&urs).Error
	userRole := map[uint]uint{}
	for _, ur := range urs { userRole[ur.UserID] = ur.RoleID }
	var view []adminView
	for _, a := range admins {
		view = append(view, adminView{Admin: a, Role: roleMap[userRole[a.ID]]})
	}
	c.HTML(http.StatusOK, "admin/admins.html", gin.H{"Admins": view, "Roles": roles})
}

func (h *AdminsHandler) CreateAdmin(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	roleID, _ := strconv.Atoi(c.PostForm("role_id"))
	if name == "" || email == "" || password == "" || roleID == 0 {
		c.String(http.StatusBadRequest, "all fields required")
		return
	}
	hash, err := utils.HashPassword(password)
	if err != nil { c.String(http.StatusInternalServerError, "hash error"); return }
	admin := models.AdminUser{Name: name, Email: email, PasswordHash: hash}
	if err := h.db.Create(&admin).Error; err != nil { c.String(http.StatusInternalServerError, "create admin failed"); return }
	_ = h.db.Create(&models.UserRole{UserID: admin.ID, RoleID: uint(roleID)}).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/admins")
}

func (h *AdminsHandler) UpdateAdmin(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var admin models.AdminUser
	if err := h.db.First(&admin, id).Error; err != nil { c.String(http.StatusNotFound, "not found"); return }
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	roleID, _ := strconv.Atoi(c.PostForm("role_id"))
	updates := map[string]interface{}{}
	if name != "" { updates["name"] = name }
	if email != "" { updates["email"] = email }
	if password != "" { if hash, err := utils.HashPassword(password); err == nil { updates["password_hash"] = hash } }
	if len(updates) > 0 { _ = h.db.Model(&admin).Updates(updates).Error }
	if roleID > 0 {
		var ur models.UserRole
		if err := h.db.Where("user_id = ?", admin.ID).First(&ur).Error; err == nil {
			_ = h.db.Model(&ur).Update("role_id", roleID).Error
		} else {
			_ = h.db.Create(&models.UserRole{UserID: admin.ID, RoleID: uint(roleID)}).Error
		}
	}
	c.Redirect(http.StatusSeeOther, "/admin/panel/admins")
}

func (h *AdminsHandler) DeleteAdmin(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = h.db.Where("user_id = ?", id).Delete(&models.UserRole{}).Error
	_ = h.db.Delete(&models.AdminUser{}, id).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/admins")
}
