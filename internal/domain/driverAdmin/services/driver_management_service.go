package services

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// GetDriversList returns paginated list of drivers with filters
func GetDriversList(c *gin.Context) {
	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Filters
	status := c.Query("status")
	search := c.Query("search")

	query := database.DB.Model(&models.Driver{})

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if search != "" {
		query = query.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ? OR phone ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Get drivers with pagination
	var drivers []models.Driver
	if err := query.Offset(offset).Limit(limit).Find(&drivers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch drivers"})
		return
	}

	// Get document status for each driver
	var result []gin.H
	for _, driver := range drivers {
		var docs models.DriverDocuments
		database.DB.Where("driver_id = ?", driver.ID).First(&docs)

		driverData := gin.H{
			"id":             driver.ID,
			"first_name":     driver.FirstName,
			"last_name":      driver.LastName,
			"email":          driver.Email,
			"phone":          driver.Phone,
			"phone_verified": driver.PhoneVerified,
			"status":         driver.Status,
			"created_at":     driver.CreatedAt,
			"documents": gin.H{
				"license_url":   docs.LicenseURL,
				"rc_url":        docs.RCURL,
				"insurance_url": docs.InsuranceURL,
				"status":        docs.Status,
			},
		}
		result = append(result, driverData)
	}

	c.JSON(http.StatusOK, gin.H{
		"drivers": result,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// GetDriverDetails returns detailed information about a specific driver
func GetDriverDetails(c *gin.Context) {
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	var driver models.Driver
	if err := database.DB.Preload("Vehicle").First(&driver, driverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	var docs models.DriverDocuments
	database.DB.Where("driver_id = ?", driver.ID).First(&docs)

	driverData := gin.H{
		"id":             driver.ID,
		"first_name":     driver.FirstName,
		"last_name":      driver.LastName,
		"email":          driver.Email,
		"phone":          driver.Phone,
		"phone_verified": driver.PhoneVerified,
		"status":         driver.Status,
		"google_id":      driver.GoogleID,
		"created_at":     driver.CreatedAt,
		"updated_at":     driver.UpdatedAt,
		"vehicle":        driver.Vehicle,
		"documents": gin.H{
			"license_url":   docs.LicenseURL,
			"rc_url":        docs.RCURL,
			"insurance_url": docs.InsuranceURL,
			"status":        docs.Status,
		},
	}

	c.JSON(http.StatusOK, driverData)
}

// ApproveDriver approves a driver and sets status to approved
func ApproveDriver(c *gin.Context) {
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	// Get admin user ID
	val, _ := c.Get("user_id")
	adminID, _ := val.(uint)

	var driver models.Driver
	if err := database.DB.First(&driver, driverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	if driver.Status != "in_review" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "driver is not in review status"})
		return
	}

	// Update driver status
	if err := database.DB.Model(&driver).Updates(map[string]interface{}{
		"status":     "approved",
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to approve driver"})
		return
	}

	// Update document status
	database.DB.Model(&models.DriverDocuments{}).Where("driver_id = ?", driverID).Update("status", "approved")

	// Log the action
	logDriverAction(adminID, uint(driverID), "approved", "Driver approved by admin")

	c.JSON(http.StatusOK, gin.H{"message": "driver approved successfully"})
}

// RejectDriver rejects a driver and sets status to rejected
func RejectDriver(c *gin.Context) {
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	// Get admin user ID
	val, _ := c.Get("user_id")
	adminID, _ := val.(uint)

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var driver models.Driver
	if err := database.DB.First(&driver, driverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	if driver.Status == "rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "driver is already rejected"})
		return
	}

	// Update driver status
	if err := database.DB.Model(&driver).Updates(map[string]interface{}{
		"status":     "rejected",
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to reject driver"})
		return
	}

	// Update document status
	database.DB.Model(&models.DriverDocuments{}).Where("driver_id = ?", driverID).Update("status", "rejected")

	// Log the action
	logDriverAction(adminID, uint(driverID), "rejected", req.Reason)

	c.JSON(http.StatusOK, gin.H{"message": "driver rejected successfully"})
}

// SuspendDriver suspends a driver
func SuspendDriver(c *gin.Context) {
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	// Get admin user ID
	val, _ := c.Get("user_id")
	adminID, _ := val.(uint)

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var driver models.Driver
	if err := database.DB.First(&driver, driverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	if driver.Status == "suspended" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "driver is already suspended"})
		return
	}

	// Update driver status
	if err := database.DB.Model(&driver).Updates(map[string]interface{}{
		"status":     "suspended",
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to suspend driver"})
		return
	}

	// Log the action
	logDriverAction(adminID, uint(driverID), "suspended", req.Reason)

	c.JSON(http.StatusOK, gin.H{"message": "driver suspended successfully"})
}

// UnsuspendDriver removes suspension from a driver
func UnsuspendDriver(c *gin.Context) {
	driverID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid driver ID"})
		return
	}

	// Get admin user ID
	val, _ := c.Get("user_id")
	adminID, _ := val.(uint)

	var driver models.Driver
	if err := database.DB.First(&driver, driverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}

	if driver.Status != "suspended" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "driver is not suspended"})
		return
	}

	// Update driver status back to approved
	if err := database.DB.Model(&driver).Updates(map[string]interface{}{
		"status":     "approved",
		"updated_at": time.Now(),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unsuspend driver"})
		return
	}

	// Log the action
	logDriverAction(adminID, uint(driverID), "unsuspended", "Driver suspension removed")

	c.JSON(http.StatusOK, gin.H{"message": "driver unsuspended successfully"})
}

// GetDriverStats returns statistics about drivers
func GetDriverStats(c *gin.Context) {
	var stats gin.H

	// Count by status
	var statusCounts []gin.H
	database.DB.Model(&models.Driver{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts)

	// Total drivers
	var totalDrivers int64
	database.DB.Model(&models.Driver{}).Count(&totalDrivers)

	// Drivers registered today
	var todayDrivers int64
	today := time.Now().Truncate(24 * time.Hour)
	database.DB.Model(&models.Driver{}).Where("created_at >= ?", today).Count(&todayDrivers)

	// Pending review
	var pendingReview int64
	database.DB.Model(&models.Driver{}).Where("status = ?", "in_review").Count(&pendingReview)

	stats = gin.H{
		"total_drivers":    totalDrivers,
		"today_drivers":    todayDrivers,
		"pending_review":   pendingReview,
		"status_breakdown": statusCounts,
	}

	c.JSON(http.StatusOK, stats)
}

// Helper function to log driver actions
func logDriverAction(adminID, driverID uint, action, reason string) {
	// This could be expanded to store in an audit log table
	fmt.Printf("Admin %d %s driver %d: %s\n", adminID, action, driverID, reason)
}
