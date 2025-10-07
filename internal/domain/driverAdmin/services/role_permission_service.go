package services

import (
	"net/http"
	"strconv"
	"time"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateRole creates a new driver admin role
func CreateRole(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Permissions []uint `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role name already exists
	var existingRole models.DriverAdminRole
	if err := database.DB.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "role name already exists"})
		return
	}

	// Create role
	role := models.DriverAdminRole{
		Name:        req.Name,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create role"})
		return
	}

	// Assign permissions
	if len(req.Permissions) > 0 {
		for _, permissionID := range req.Permissions {
			rolePermission := models.DriverAdminRolePermission{
				RoleID:       role.ID,
				PermissionID: permissionID,
				CreatedAt:    time.Now(),
			}
			database.DB.Create(&rolePermission)
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "role created successfully",
		"role": gin.H{
			"id":          role.ID,
			"name":        role.Name,
			"description": role.Description,
		},
	})
}

// GetRolesList returns list of all roles with their permissions
func GetRolesList(c *gin.Context) {
	var roles []models.DriverAdminRole
	if err := database.DB.Preload("RolePermissions.Permission").Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch roles"})
		return
	}

	var result []gin.H
	for _, role := range roles {
		permissions := make([]gin.H, len(role.RolePermissions))
		for i, rp := range role.RolePermissions {
			permissions[i] = gin.H{
				"id":          rp.Permission.ID,
				"name":        rp.Permission.Name,
				"description": rp.Permission.Description,
			}
		}

		roleData := gin.H{
			"id":          role.ID,
			"name":        role.Name,
			"description": role.Description,
			"permissions": permissions,
			"created_at":  role.CreatedAt,
		}
		result = append(result, roleData)
	}

	c.JSON(http.StatusOK, gin.H{"roles": result})
}

// UpdateRole updates a role and its permissions
func UpdateRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Permissions []uint `json:"permissions"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if role exists
	var role models.DriverAdminRole
	if err := database.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	// Check name uniqueness if changing
	if req.Name != "" && req.Name != role.Name {
		var existingRole models.DriverAdminRole
		if err := database.DB.Where("name = ? AND id != ?", req.Name, roleID).First(&existingRole).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "role name already exists"})
			return
		}
	}

	// Update role fields
	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["updated_at"] = time.Now()

	if err := database.DB.Model(&role).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update role"})
		return
	}

	// Update permissions if provided
	if req.Permissions != nil {
		// Remove existing permissions
		database.DB.Where("role_id = ?", roleID).Delete(&models.DriverAdminRolePermission{})

		// Add new permissions
		for _, permissionID := range req.Permissions {
			rolePermission := models.DriverAdminRolePermission{
				RoleID:       uint(roleID),
				PermissionID: permissionID,
				CreatedAt:    time.Now(),
			}
			database.DB.Create(&rolePermission)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "role updated successfully"})
}

// DeleteRole deletes a role (soft delete by setting inactive)
func DeleteRole(c *gin.Context) {
	roleID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	// Check if role exists
	var role models.DriverAdminRole
	if err := database.DB.First(&role, roleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "role not found"})
		return
	}

	// Check if role is assigned to any admin
	var count int64
	database.DB.Model(&models.DriverAdminAccountRole{}).Where("role_id = ?", roleID).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot delete role that is assigned to admin accounts"})
		return
	}

	// Delete role and its permissions
	database.DB.Where("role_id = ?", roleID).Delete(&models.DriverAdminRolePermission{})
	database.DB.Delete(&role)

	c.JSON(http.StatusOK, gin.H{"message": "role deleted successfully"})
}

// GetPermissionsList returns all available permissions
func GetPermissionsList(c *gin.Context) {
	var permissions []models.DriverAdminPermission
	if err := database.DB.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch permissions"})
		return
	}

	var result []gin.H
	for _, perm := range permissions {
		permData := gin.H{
			"id":          perm.ID,
			"name":        perm.Name,
			"description": perm.Description,
			"created_at":  perm.CreatedAt,
		}
		result = append(result, permData)
	}

	c.JSON(http.StatusOK, gin.H{"permissions": result})
}

// AssignRoleToAdmin assigns a role to a driver admin
func AssignRoleToAdmin(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	var req struct {
		RoleID uint `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if admin exists
	var admin models.DriverAdmin
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	// Check if role exists
	var role models.DriverAdminRole
	if err := database.DB.First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	// Remove existing roles
	database.DB.Where("driver_admin_id = ?", adminID).Delete(&models.DriverAdminAccountRole{})

	// Assign new role
	accountRole := models.DriverAdminAccountRole{
		DriverAdminID: uint(adminID),
		RoleID:        req.RoleID,
		AssignedAt:    time.Now(),
	}

	if err := database.DB.Create(&accountRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "role assigned successfully"})
}

// GetAdminPermissions returns permissions for a specific admin
func GetAdminPermissions(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	// Get admin with roles and permissions
	var admin models.DriverAdmin
	if err := database.DB.Preload("AccountRoles.Role.RolePermissions.Permission").First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	// Collect all permissions
	permissionMap := make(map[uint]gin.H)
	for _, accountRole := range admin.AccountRoles {
		for _, rolePermission := range accountRole.Role.RolePermissions {
			permissionMap[rolePermission.Permission.ID] = gin.H{
				"id":          rolePermission.Permission.ID,
				"name":        rolePermission.Permission.Name,
				"description": rolePermission.Permission.Description,
			}
		}
	}

	// Convert map to slice
	var permissions []gin.H
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	c.JSON(http.StatusOK, gin.H{
		"admin_id":    admin.ID,
		"permissions": permissions,
	})
}
