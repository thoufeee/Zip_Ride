package adminhandlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"zipRideDriver/internal/models"
)

type RoleHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewRoleHandler(db *gorm.DB, log *zap.Logger) *RoleHandler {
	return &RoleHandler{db: db, log: log}
}

func (h *RoleHandler) RolesPage(c *gin.Context) {
	var roles []models.Role
	_ = h.db.Order("name asc").Find(&roles).Error
	var perms []models.Permission
	_ = h.db.Order("name asc").Find(&perms).Error
	c.HTML(http.StatusOK, "admin/roles.html", gin.H{"Roles": roles, "Permissions": perms})
}

func (h *RoleHandler) NewRolePage(c *gin.Context) {
	var perms []models.Permission
	_ = h.db.Order("name asc").Find(&perms).Error
	c.HTML(http.StatusOK, "admin/role_new.html", gin.H{"Permissions": perms})
}

func (h *RoleHandler) CreateRole(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	if name == "" {
		c.String(http.StatusBadRequest, "name required")
		return
	}
	role := models.Role{Name: name}
	if err := h.db.Create(&role).Error; err != nil {
		c.String(http.StatusInternalServerError, "failed to create role")
		return
	}
	permNames := c.PostFormArray("permissions")
	var perms []models.Permission
	if len(permNames) > 0 {
		_ = h.db.Where("name IN ?", permNames).Find(&perms).Error
		for _, p := range perms {
			_ = h.db.Create(&models.RolePermission{RoleID: role.ID, PermissionID: p.ID}).Error
		}
	}
	c.Redirect(http.StatusSeeOther, "/admin/panel/roles")
}

func (h *RoleHandler) EditRolePage(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role
	if err := h.db.First(&role, id).Error; err != nil {
		c.String(http.StatusNotFound, "role not found")
		return
	}
	var perms []models.Permission
	_ = h.db.Order("name asc").Find(&perms).Error
	var rps []models.RolePermission
	_ = h.db.Where("role_id = ?", role.ID).Find(&rps).Error
	selected := map[uint]bool{}
	for _, rp := range rps { selected[rp.PermissionID] = true }
	c.HTML(http.StatusOK, "admin/role_edit.html", gin.H{"Role": role, "Permissions": perms, "Selected": selected})
}

func (h *RoleHandler) UpdateRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var role models.Role
	if err := h.db.First(&role, id).Error; err != nil {
		c.String(http.StatusNotFound, "role not found")
		return
	}
	name := strings.TrimSpace(c.PostForm("name"))
	if name != "" { _ = h.db.Model(&role).Update("name", name).Error }
	_ = h.db.Where("role_id = ?", role.ID).Delete(&models.RolePermission{}).Error
	permNames := c.PostFormArray("permissions")
	var perms []models.Permission
	if len(permNames) > 0 {
		_ = h.db.Where("name IN ?", permNames).Find(&perms).Error
		for _, p := range perms {
			_ = h.db.Create(&models.RolePermission{RoleID: role.ID, PermissionID: p.ID}).Error
		}
	}
	c.Redirect(http.StatusSeeOther, "/admin/panel/roles")
}

func (h *RoleHandler) DeleteRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	_ = h.db.Where("role_id = ?", id).Delete(&models.RolePermission{}).Error
	_ = h.db.Delete(&models.Role{}, id).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/roles")
}
