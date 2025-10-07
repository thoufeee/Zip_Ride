package services

import (
	"net/http"
	"strconv"
	"time"
	"zipride/database"
	"zipride/internal/models"
	"zipride/utils"

	"github.com/gin-gonic/gin"
)

// CreateDriverAdminStaff creates a new driver admin staff member
func CreateDriverAdminStaff(c *gin.Context) {
	var req struct {
		FirstName string `json:"first_name" binding:"required"`
		LastName  string `json:"last_name" binding:"required"`
		Email     string `json:"email" binding:"required,email"`
		Password  string `json:"password" binding:"required,min=6"`
		RoleID    uint   `json:"role_id" binding:"required"`
		IsActive  bool   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingAdmin models.DriverAdmin
	if err := database.DB.Where("email = ?", req.Email).First(&existingAdmin).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	// Verify role exists
	var role models.DriverAdminRole
	if err := database.DB.First(&role, req.RoleID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
		return
	}

	// Hash password
	hashedPassword, err := utils.GenerateHash(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create admin account
	admin := models.DriverAdmin{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsActive:     req.IsActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := database.DB.Create(&admin).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create admin account"})
		return
	}

	// Assign role
	accountRole := models.DriverAdminAccountRole{
		DriverAdminID: admin.ID,
		RoleID:        req.RoleID,
		AssignedAt:    time.Now(),
	}

	if err := database.DB.Create(&accountRole).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to assign role"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "admin staff created successfully",
		"admin": gin.H{
			"id":         admin.ID,
			"first_name": admin.FirstName,
			"last_name":  admin.LastName,
			"email":      admin.Email,
			"is_active":  admin.IsActive,
			"role_id":    req.RoleID,
		},
	})
}

// GetDriverAdminStaffList returns list of driver admin staff
func GetDriverAdminStaffList(c *gin.Context) {
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Filters
	active := c.Query("active")
	search := c.Query("search")

	query := database.DB.Model(&models.DriverAdmin{})

	// Apply filters
	if active != "" {
		query = query.Where("is_active = ?", active == "true")
	}
	if search != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get admins with roles
	var admins []models.DriverAdmin
	if err := query.Preload("AccountRoles.Role").Offset(offset).Limit(limit).Find(&admins).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch admin staff"})
		return
	}

	// Format response
	var result []gin.H
	for _, admin := range admins {
		roles := make([]gin.H, len(admin.AccountRoles))
		for i, ar := range admin.AccountRoles {
			roles[i] = gin.H{
				"id":   ar.Role.ID,
				"name": ar.Role.Name,
			}
		}

		adminData := gin.H{
			"id":         admin.ID,
			"first_name": admin.FirstName,
			"last_name":  admin.LastName,
			"email":      admin.Email,
			"is_active":  admin.IsActive,
			"roles":      roles,
			"created_at": admin.CreatedAt,
		}
		result = append(result, adminData)
	}

	c.JSON(http.StatusOK, gin.H{
		"staff": result,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// UpdateDriverAdminStaff updates driver admin staff information
func UpdateDriverAdminStaff(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		IsActive  *bool  `json:"is_active"`
		RoleID    uint   `json:"role_id"`
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

	// Check email uniqueness if changing
	if req.Email != "" && req.Email != admin.Email {
		var existingAdmin models.DriverAdmin
		if err := database.DB.Where("email = ? AND id != ?", req.Email, adminID).First(&existingAdmin).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
			return
		}
	}

	// Update fields
	updates := make(map[string]interface{})
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}
	updates["updated_at"] = time.Now()

	if err := database.DB.Model(&admin).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update admin"})
		return
	}

	// Update role if provided
	if req.RoleID != 0 {
		// Verify role exists
		var role models.DriverAdminRole
		if err := database.DB.First(&role, req.RoleID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role ID"})
			return
		}

		// Remove existing roles and assign new one
		database.DB.Where("driver_admin_id = ?", adminID).Delete(&models.DriverAdminAccountRole{})

		accountRole := models.DriverAdminAccountRole{
			DriverAdminID: uint(adminID),
			RoleID:        req.RoleID,
			AssignedAt:    time.Now(),
		}
		database.DB.Create(&accountRole)
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin staff updated successfully"})
}

// DeleteDriverAdminStaff deletes a driver admin staff member
func DeleteDriverAdminStaff(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	// Check if admin exists
	var admin models.DriverAdmin
	if err := database.DB.First(&admin, adminID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "admin not found"})
		return
	}

	// Soft delete (set is_active to false)
	if err := database.DB.Model(&admin).Update("is_active", false).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate admin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "admin staff deactivated successfully"})
}

// ChangeDriverAdminPassword changes password for driver admin staff
func ChangeDriverAdminPassword(c *gin.Context) {
	adminID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid admin ID"})
		return
	}

	var req struct {
		NewPassword string `json:"new_password" binding:"required,min=6"`
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

	// Hash new password
	hashedPassword, err := utils.GenerateHash(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Update password
	if err := database.DB.Model(&admin).Updates(map[string]interface{}{
		"password_hash": hashedPassword,
		"updated_at":    time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password updated successfully"})
}
