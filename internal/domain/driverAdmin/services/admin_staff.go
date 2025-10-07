package services

import (
	"fmt"
	"net/http"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// CreateDriverAdmin creates a new driver admin staff account
func CreateDriverAdmin(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
		Name     string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Password == "" || !utils.PasswordStrength(req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weak password"})
		return
	}
	var existing models.DriverAdmin
	if err := database.DB.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	hash, _ := utils.GenerateHash(req.Password)
	adm := models.DriverAdmin{Email: req.Email, Password: hash, Name: req.Name}
	if err := database.DB.Create(&adm).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": adm.ID, "email": adm.Email, "name": adm.Name})
}

// GetDriverAdmin returns details for a driver admin staff account
func GetDriverAdmin(c *gin.Context) {
	id := c.Param("id")
	var adm models.DriverAdmin
	if err := database.DB.First(&adm, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": adm.ID, "email": adm.Email, "name": adm.Name})
}

// UpdateDriverAdmin updates basic fields of a staff account
func UpdateDriverAdmin(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]any{}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if len(updates) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no updates"})
		return
	}
	if err := database.DB.Model(&models.DriverAdmin{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

// AssignRoles assigns role IDs to an admin account (replaces all roles)
func AssignRoles(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		RoleIDs []uint `json:"role_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.RoleIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role_ids required"})
		return
	}
	// delete existing
	database.DB.Where("admin_id = ?", id).Delete(&models.DriverAdminAccountRole{})
	for _, rid := range req.RoleIDs {
		database.DB.Create(&models.DriverAdminAccountRole{AdminID: toUint(id), RoleID: rid})
	}
	c.JSON(http.StatusOK, gin.H{"message": "roles assigned"})
}

// ListRoles lists all driver admin roles
func ListRoles(c *gin.Context) {
	var roles []models.DriverAdminRole
	if err := database.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list roles"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"roles": roles})
}


// SetRolePermissions sets permission keys on a role (replaces all)
func SetRolePermissions(c *gin.Context) {
	rid := c.Param("id")
	var req struct {
		PermissionKeys []string `json:"permission_keys" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.PermissionKeys) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "permission_keys required"})
		return
	}
	// find permissions by key
	var perms []models.DriverAdminPermission
	if err := database.DB.Where("key IN ?", req.PermissionKeys).Find(&perms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to query permissions"})
		return
	}
	// replace mappings
	database.DB.Where("role_id = ?", rid).Delete(&models.DriverAdminRolePermission{})
	for _, p := range perms {
		database.DB.Create(&models.DriverAdminRolePermission{RoleID: toUint(rid), PermissionID: p.ID})
	}
	c.JSON(http.StatusOK, gin.H{"message": "permissions updated"})
}

func toUint(s string) uint {
	var u uint
	_, _ = fmt.Sscan(s, &u)
	return u
}
