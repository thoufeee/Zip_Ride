package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/services"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AdminHandler struct {
	cfg *config.Config
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

// handleError provides consistent error handling for admin operations
func (h *AdminHandler) handleError(c *gin.Context, err error, template string, message string, isHTMLRequest bool) {
	h.log.Error("admin operation failed", zap.Error(err), zap.String("template", template))
	
	if isHTMLRequest {
		c.HTML(http.StatusInternalServerError, template, gin.H{
			"Error": message,
		})
	} else {
		utils.Error(c, http.StatusInternalServerError, message)
	}
}

// validateAdminAccess checks if the current user has admin access
func (h *AdminHandler) validateAdminAccess(c *gin.Context) (*models.AdminUser, bool) {
	email, err := c.Cookie("admin_token")
	if err != nil {
		return nil, false
	}

	var admin models.AdminUser
	if err := h.db.Where("email = ?", email).First(&admin).Error; err != nil {
		return nil, false
	}

	return &admin, true
}

func NewAdminHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *AdminHandler {
	return &AdminHandler{cfg: cfg, db: db, rdb: rdb, log: log}
}

func (h *AdminHandler) Login(c *gin.Context) {
	var email, password string
	
	// Check if it's a form submission (HTML) or JSON (API)
	contentType := c.GetHeader("Content-Type")
	isFormData := strings.Contains(contentType, "application/x-www-form-urlencoded") || 
				 strings.Contains(contentType, "multipart/form-data") ||
				 contentType == ""

	if isFormData {
		// Handle form data (HTML login)
		email = strings.TrimSpace(c.PostForm("email"))
		password = c.PostForm("password")
		if email == "" || password == "" {
			c.HTML(http.StatusBadRequest, "admin/login.html", gin.H{"Error": "Email and password are required"})
			return
		}
	} else {
		// Handle JSON data (API login)
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		email = strings.TrimSpace(req.Email)
		password = req.Password
	}

	var admin models.AdminUser
	if err := h.db.Where("email = ?", email).First(&admin).Error; err != nil {
		if isFormData {
			c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{"Error": "Invalid email or password"})
		} else {
			utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}
	
	if !utils.CheckPassword(admin.PasswordHash, password) {
		if isFormData {
			c.HTML(http.StatusUnauthorized, "admin/login.html", gin.H{"Error": "Invalid email or password"})
		} else {
			utils.Error(c, http.StatusUnauthorized, "invalid credentials")
		}
		return
	}

	if isFormData {
		// Set cookie for HTML login and redirect
		c.SetCookie("admin_token", admin.Email, 3600, "/", "", false, true)
		h.log.Info("admin login successful", zap.String("email", admin.Email))
		c.Redirect(http.StatusSeeOther, "/admin/dashboard")
	} else {
		// Generate JWT tokens for API login
		access, err := services.GenerateToken(admin.ID, admin.Email, "admin", []string{"admin"}, h.cfg.JWTSecret, h.cfg.AccessTokenExpiry)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to issue token")
			return
		}
		refresh, err := services.GenerateToken(admin.ID, admin.Email, "admin_refresh", []string{"admin"}, h.cfg.JWTSecret, h.cfg.RefreshTokenExpiry)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to issue token")
			return
		}
		utils.Ok(c, "login success", gin.H{"access_token": access, "refresh_token": refresh})
	}
}

func (h *AdminHandler) ListDrivers(c *gin.Context) {
	// Check if it's an HTML request (browser) or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	// Get query parameters for filtering
	status := c.Query("status")
	page := c.DefaultQuery("page", "1")

	var drivers []models.Driver
	query := h.db.Order("created_at desc")

	// Apply status filter if provided
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// For API requests, apply pagination
	if !isHTMLRequest {
		var offset int
		if pageNum := c.Query("page"); pageNum != "" {
			if p, err := strconv.Atoi(pageNum); err == nil && p > 0 {
				offset = (p - 1) * 50
			}
		}
		query = query.Offset(offset).Limit(50)
	} else {
		// For HTML requests, limit to reasonable number
		query = query.Limit(200)
	}

	if err := query.Find(&drivers).Error; err != nil {
		if isHTMLRequest {
			c.HTML(http.StatusInternalServerError, "admin/drivers.html", gin.H{
				"Error": "Failed to load drivers",
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to list drivers")
		}
		return
	}

	if isHTMLRequest {
		// Render HTML template
		c.HTML(http.StatusOK, "admin/drivers.html", gin.H{
			"Drivers":      drivers,
			"FilterStatus": status,
			"CurrentPage":  page,
		})
	} else {
		// Return JSON for API
		utils.Ok(c, "drivers", drivers)
	}
}

func (h *AdminHandler) ApproveDriver(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	
	// Check if it's an HTML request or API request
	contentType := c.GetHeader("Content-Type")
	isHTMLRequest := !strings.Contains(contentType, "application/json")

	var status string
	if isHTMLRequest {
		// For HTML requests, approve directly
		status = "Approved"
	} else {
		// For API requests, check JSON body
		var body struct {
			Action string `json:"action" binding:"required,oneof=approve reject"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		status = "Approved"
		if body.Action == "reject" {
			status = "Rejected"
		}
	}

	if err := h.db.Model(&models.Driver{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		if isHTMLRequest {
			c.Redirect(http.StatusSeeOther, "/admin/drivers?error=failed_to_update")
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to update status")
		}
		return
	}

	h.log.Info("driver status updated", zap.String("driver_id", id), zap.String("status", status))

	if isHTMLRequest {
		c.Redirect(http.StatusSeeOther, "/admin/drivers")
	} else {
		utils.Ok(c, "updated", gin.H{"status": status})
	}
}

func (h *AdminHandler) SuspendDriver(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	
	if err := h.db.Model(&models.Driver{}).Where("id = ?", id).Update("status", "Suspended").Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/drivers?error=failed_to_suspend")
		return
	}

	h.log.Info("driver suspended", zap.String("driver_id", id))
	c.Redirect(http.StatusSeeOther, "/admin/drivers")
}

func (h *AdminHandler) BanDriver(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if err := h.db.Model(&models.Driver{}).Where("id = ?", id).Update("status", "Banned").Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to ban driver")
		return
	}
	utils.Ok(c, "banned", nil)
}

// ListRides handles both HTML and API requests for ride listing
func (h *AdminHandler) ListRides(c *gin.Context) {
	// Check if it's an HTML request (browser) or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	// Get query parameters for filtering
	status := c.Query("status")
	driverID := c.Query("driver_id")
	page := c.DefaultQuery("page", "1")

	var rides []models.Ride
	query := h.db.Preload("Driver").Preload("Rider").Order("created_at desc")

	// Apply status filter if provided
	if status != "" && status != "all" {
		query = query.Where("status = ?", status)
	}

	// Apply driver filter if provided
	if driverID != "" && driverID != "all" {
		query = query.Where("driver_id = ?", driverID)
	}

	// For API requests, apply pagination
	if !isHTMLRequest {
		var offset int
		if pageNum := c.Query("page"); pageNum != "" {
			if p, err := strconv.Atoi(pageNum); err == nil && p > 0 {
				offset = (p - 1) * 50
			}
		}
		query = query.Offset(offset).Limit(50)
	} else {
		// For HTML requests, limit to reasonable number
		query = query.Limit(200)
	}

	if err := query.Find(&rides).Error; err != nil {
		if isHTMLRequest {
			c.HTML(http.StatusInternalServerError, "admin/rides.html", gin.H{
				"Error": "Failed to load rides",
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to list rides")
		}
		return
	}

	if isHTMLRequest {
		// Get drivers for filter dropdown
		var drivers []models.Driver
		h.db.Where("status = ?", "Approved").Order("name asc").Find(&drivers)

		// Render HTML template
		c.HTML(http.StatusOK, "admin/rides.html", gin.H{
			"Rides":        rides,
			"Drivers":      drivers,
			"FilterStatus": status,
			"FilterDriver": driverID,
			"CurrentPage":  page,
		})
	} else {
		// Return JSON for API
		utils.Ok(c, "rides", rides)
	}
}

// ShowRide displays detailed information about a specific ride
func (h *AdminHandler) ShowRide(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin/rides.html", gin.H{
			"Error": "Invalid ride ID",
		})
		return
	}

	var ride models.Ride
	if err := h.db.Preload("Driver").Preload("Rider").First(&ride, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/rides.html", gin.H{
			"Error": "Ride not found",
		})
		return
	}

	// Check if it's an HTML request or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	if isHTMLRequest {
		c.HTML(http.StatusOK, "admin/ride_detail.html", gin.H{
			"Ride": ride,
		})
	} else {
		utils.Ok(c, "ride", ride)
	}
}

// CancelRide allows admin to cancel a ride
func (h *AdminHandler) CancelRide(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	
	if err := h.db.Model(&models.Ride{}).Where("id = ?", id).Update("status", "cancelled").Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/rides?error=failed_to_cancel")
		return
	}

	h.log.Info("ride cancelled by admin", zap.String("ride_id", id))
	c.Redirect(http.StatusSeeOther, "/admin/rides")
}

// ListEarnings handles both HTML and API requests for earnings listing
func (h *AdminHandler) ListEarnings(c *gin.Context) {
	// Check if it's an HTML request (browser) or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	// Get query parameters for filtering
	driverID := c.Query("driver_id")

	type EarningSummary struct {
		DriverID      uint    `json:"driver_id"`
		DriverName    string  `json:"driver_name"`
		DriverPhone   string  `json:"driver_phone"`
		TotalEarnings float64 `json:"total_earnings"`
		RideCount     int64   `json:"ride_count"`
		AvgEarning    float64 `json:"avg_earning"`
	}

	var earnings []EarningSummary
	query := h.db.Table("earnings").
		Select("drivers.id as driver_id, drivers.name as driver_name, drivers.phone as driver_phone, SUM(earnings.amount) as total_earnings, COUNT(earnings.id) as ride_count, AVG(earnings.amount) as avg_earning").
		Joins("JOIN drivers ON drivers.id = earnings.driver_id").
		Group("drivers.id, drivers.name, drivers.phone").
		Order("total_earnings DESC")

	// Apply driver filter if provided
	if driverID != "" && driverID != "all" {
		query = query.Where("drivers.id = ?", driverID)
	}

	if err := query.Scan(&earnings).Error; err != nil {
		if isHTMLRequest {
			c.HTML(http.StatusInternalServerError, "admin/earnings.html", gin.H{
				"Error": "Failed to load earnings data",
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to list earnings")
		}
		return
	}

	if isHTMLRequest {
		// Get drivers for filter dropdown
		var drivers []models.Driver
		h.db.Where("status = ?", "Approved").Order("name asc").Find(&drivers)

		// Calculate totals
		var totalEarnings float64
		var totalRides int64
		for _, earning := range earnings {
			totalEarnings += earning.TotalEarnings
			totalRides += earning.RideCount
		}

		// Calculate average per driver
		var avgPerDriver float64
		if len(earnings) > 0 {
			avgPerDriver = totalEarnings / float64(len(earnings))
		}

		// Render HTML template
		c.HTML(http.StatusOK, "admin/earnings.html", gin.H{
			"Earnings":      earnings,
			"Drivers":       drivers,
			"FilterDriver":  driverID,
			"TotalEarnings": totalEarnings,
			"TotalRides":    totalRides,
			"DriverCount":   len(earnings),
			"AvgPerDriver":  avgPerDriver,
		})
	} else {
		// Return JSON for API
		utils.Ok(c, "earnings", earnings)
	}
}

// ExportEarnings generates CSV export of earnings data
func (h *AdminHandler) ExportEarnings(c *gin.Context) {
	type EarningSummary struct {
		DriverID      uint    `json:"driver_id"`
		DriverName    string  `json:"driver_name"`
		DriverPhone   string  `json:"driver_phone"`
		TotalEarnings float64 `json:"total_earnings"`
		RideCount     int64   `json:"ride_count"`
		AvgEarning    float64 `json:"avg_earning"`
	}

	var earnings []EarningSummary
	if err := h.db.Table("earnings").
		Select("drivers.id as driver_id, drivers.name as driver_name, drivers.phone as driver_phone, SUM(earnings.amount) as total_earnings, COUNT(earnings.id) as ride_count, AVG(earnings.amount) as avg_earning").
		Joins("JOIN drivers ON drivers.id = earnings.driver_id").
		Group("drivers.id, drivers.name, drivers.phone").
		Order("total_earnings DESC").
		Scan(&earnings).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to export earnings")
		return
	}

	// Set CSV headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=earnings_report.csv")

	// Write CSV header
	c.String(http.StatusOK, "Driver ID,Driver Name,Phone,Total Earnings,Ride Count,Average Earning\n")

	// Write CSV data
	for _, earning := range earnings {
		c.String(http.StatusOK, "%d,%s,%s,%.2f,%d,%.2f\n",
			earning.DriverID,
			earning.DriverName,
			earning.DriverPhone,
			earning.TotalEarnings,
			earning.RideCount,
			earning.AvgEarning,
		)
	}
}

func (h *AdminHandler) Dashboard(c *gin.Context) {
	var totalDrivers int64
	var activeDrivers int64
	var totalRides int64
	var completedRides int64
	var totalEarnings float64

	_ = h.db.Model(&models.Driver{}).Count(&totalDrivers).Error
	_ = h.db.Model(&models.Driver{}).Where("status = ? OR is_online = ?", "Approved", true).Count(&activeDrivers).Error
	_ = h.db.Model(&models.Ride{}).Count(&totalRides).Error
	_ = h.db.Model(&models.Ride{}).Where("status = ?", "completed").Count(&completedRides).Error
	_ = h.db.Model(&models.Earning{}).Select("COALESCE(SUM(amount),0)").Scan(&totalEarnings).Error

	// Check if client wants JSON response
	if c.GetHeader("Accept") == "application/json" || c.Query("format") == "json" {
		utils.Ok(c, "dashboard", gin.H{
			"total_drivers":   totalDrivers,
			"active_drivers":  activeDrivers,
			"total_rides":     totalRides,
			"completed_rides": completedRides,
			"total_earnings":  totalEarnings,
		})
		return
	}

	// Render HTML template
	c.HTML(http.StatusOK, "admin/dashboard.html", gin.H{
		"TotalDrivers":   totalDrivers,
		"ActiveDrivers":  activeDrivers,
		"TotalRides":     totalRides,
		"CompletedRides": completedRides,
		"TotalEarnings":  totalEarnings,
	})
}

// ListReports handles both HTML and API requests for reports and analytics
func (h *AdminHandler) ListReports(c *gin.Context) {
	// Check if it's an HTML request (browser) or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	// Get query parameters for date filtering
	period := c.DefaultQuery("period", "30") // Default to 30 days

	type ReportData struct {
		TotalRides      int64   `json:"total_rides"`
		TotalEarnings   float64 `json:"total_earnings"`
		TotalDrivers    int64   `json:"total_drivers"`
		ActiveDrivers   int64   `json:"active_drivers"`
		CompletedRides  int64   `json:"completed_rides"`
		CancelledRides  int64   `json:"cancelled_rides"`
		AvgRideValue    float64 `json:"avg_ride_value"`
		CompletionRate  float64 `json:"completion_rate"`
		TopDrivers      []struct {
			Name       string  `json:"name"`
			Rides      int64   `json:"rides"`
			Earnings   float64 `json:"earnings"`
			AvgPerRide float64 `json:"avg_per_ride"`
		} `json:"top_drivers"`
		DailyStats []struct {
			Date     string  `json:"date"`
			Rides    int64   `json:"rides"`
			Earnings float64 `json:"earnings"`
		} `json:"daily_stats"`
		StatusBreakdown []struct {
			Status string `json:"status"`
			Count  int64  `json:"count"`
		} `json:"status_breakdown"`
	}

	var report ReportData

	// Basic statistics
	h.db.Model(&models.Ride{}).Count(&report.TotalRides)
	h.db.Model(&models.Earning{}).Select("COALESCE(SUM(amount),0)").Scan(&report.TotalEarnings)
	h.db.Model(&models.Driver{}).Count(&report.TotalDrivers)
	h.db.Model(&models.Driver{}).Where("status = ?", "Approved").Count(&report.ActiveDrivers)
	h.db.Model(&models.Ride{}).Where("status = ?", "completed").Count(&report.CompletedRides)
	h.db.Model(&models.Ride{}).Where("status = ?", "cancelled").Count(&report.CancelledRides)

	// Calculate average ride value and completion rate
	if report.CompletedRides > 0 {
		report.AvgRideValue = report.TotalEarnings / float64(report.CompletedRides)
	}
	if report.TotalRides > 0 {
		report.CompletionRate = (float64(report.CompletedRides) / float64(report.TotalRides)) * 100
	}

	// Top performing drivers with average per ride calculation
	type TopDriverRaw struct {
		Name     string
		Rides    int64
		Earnings float64
	}
	var topDriversRaw []TopDriverRaw
	h.db.Table("earnings").
		Select("drivers.name, COUNT(earnings.id) as rides, SUM(earnings.amount) as earnings").
		Joins("JOIN drivers ON drivers.id = earnings.driver_id").
		Group("drivers.id, drivers.name").
		Order("earnings DESC").
		Limit(5).
		Scan(&topDriversRaw)

	// Calculate average per ride for each driver
	for _, driver := range topDriversRaw {
		avgPerRide := float64(0)
		if driver.Rides > 0 {
			avgPerRide = driver.Earnings / float64(driver.Rides)
		}
		report.TopDrivers = append(report.TopDrivers, struct {
			Name       string  `json:"name"`
			Rides      int64   `json:"rides"`
			Earnings   float64 `json:"earnings"`
			AvgPerRide float64 `json:"avg_per_ride"`
		}{
			Name:       driver.Name,
			Rides:      driver.Rides,
			Earnings:   driver.Earnings,
			AvgPerRide: avgPerRide,
		})
	}

	// Daily statistics for the last 7 days (simplified for SQLite compatibility)
	h.db.Raw(`
		SELECT 
			date(created_at) as date,
			COUNT(*) as rides,
			COALESCE(SUM(fare), 0) as earnings
		FROM rides 
		WHERE created_at >= datetime('now', '-7 days')
		GROUP BY date(created_at)
		ORDER BY date DESC
	`).Scan(&report.DailyStats)

	// Status breakdown
	h.db.Model(&models.Ride{}).
		Select("status, COUNT(*) as count").
		Group("status").
		Scan(&report.StatusBreakdown)

	if isHTMLRequest {
		// Render HTML template
		c.HTML(http.StatusOK, "admin/reports.html", gin.H{
			"Report": report,
			"Period": period,
		})
	} else {
		// Return JSON for API
		utils.Ok(c, "reports", report)
	}
}

// ExportReports generates comprehensive CSV report
func (h *AdminHandler) ExportReports(c *gin.Context) {
	// Set CSV headers
	c.Header("Content-Type", "text/csv")
	c.Header("Content-Disposition", "attachment; filename=platform_report.csv")

	// Write comprehensive report
	c.String(http.StatusOK, "ZipRide Platform Report\n\n")
	
	// Basic metrics
	var totalRides, totalDrivers, activeDrivers, completedRides int64
	var totalEarnings float64
	
	h.db.Model(&models.Ride{}).Count(&totalRides)
	h.db.Model(&models.Driver{}).Count(&totalDrivers)
	h.db.Model(&models.Driver{}).Where("status = ?", "Approved").Count(&activeDrivers)
	h.db.Model(&models.Ride{}).Where("status = ?", "completed").Count(&completedRides)
	h.db.Model(&models.Earning{}).Select("COALESCE(SUM(amount),0)").Scan(&totalEarnings)

	c.String(http.StatusOK, "Platform Overview\n")
	c.String(http.StatusOK, "Total Rides,%d\n", totalRides)
	c.String(http.StatusOK, "Total Earnings,%.2f\n", totalEarnings)
	c.String(http.StatusOK, "Total Drivers,%d\n", totalDrivers)
	c.String(http.StatusOK, "Active Drivers,%d\n", activeDrivers)
	c.String(http.StatusOK, "Completed Rides,%d\n", completedRides)
	
	// Top drivers
	type TopDriver struct {
		Name     string
		Rides    int64
		Earnings float64
	}
	var topDrivers []TopDriver
	h.db.Table("earnings").
		Select("drivers.name, COUNT(earnings.id) as rides, SUM(earnings.amount) as earnings").
		Joins("JOIN drivers ON drivers.id = earnings.driver_id").
		Group("drivers.id, drivers.name").
		Order("earnings DESC").
		Limit(10).
		Scan(&topDrivers)

	c.String(http.StatusOK, "\nTop Drivers\n")
	c.String(http.StatusOK, "Name,Rides,Earnings\n")
	for _, driver := range topDrivers {
		c.String(http.StatusOK, "%s,%d,%.2f\n", driver.Name, driver.Rides, driver.Earnings)
	}
}

// AdminSettings displays admin profile and settings
func (h *AdminHandler) AdminSettings(c *gin.Context) {
	// Check if it's an HTML request (browser) or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	// Get admin email from cookie
	email, err := c.Cookie("admin_token")
	if err != nil {
		if isHTMLRequest {
			c.Redirect(http.StatusSeeOther, "/admin/login")
		} else {
			utils.Error(c, http.StatusUnauthorized, "not authenticated")
		}
		return
	}

	// Find admin user
	var admin models.AdminUser
	if err := h.db.Where("email = ?", email).First(&admin).Error; err != nil {
		if isHTMLRequest {
			c.HTML(http.StatusInternalServerError, "admin/settings.html", gin.H{
				"Error": "Failed to load admin profile",
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to load admin profile")
		}
		return
	}

	if isHTMLRequest {
		// Render HTML template
		c.HTML(http.StatusOK, "admin/settings.html", gin.H{
			"Admin": admin,
		})
	} else {
		// Return JSON for API (without password hash)
		utils.Ok(c, "admin", gin.H{
			"id":    admin.ID,
			"name":  admin.Name,
			"email": admin.Email,
			"role":  admin.Role,
		})
	}
}

// ChangePassword handles admin password change
func (h *AdminHandler) ChangePassword(c *gin.Context) {
	// Check if it's a form submission or JSON request
	contentType := c.GetHeader("Content-Type")
	isFormData := strings.Contains(contentType, "application/x-www-form-urlencoded") || 
				 strings.Contains(contentType, "multipart/form-data") ||
				 contentType == ""

	var newPassword, currentPassword string
	
	if isFormData {
		// Handle form data
		newPassword = strings.TrimSpace(c.PostForm("new_password"))
		currentPassword = strings.TrimSpace(c.PostForm("current_password"))
	} else {
		// Handle JSON data
		var req struct {
			NewPassword     string `json:"new_password" binding:"required,min=6"`
			CurrentPassword string `json:"current_password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		newPassword = req.NewPassword
		currentPassword = req.CurrentPassword
	}

	if newPassword == "" {
		if isFormData {
			c.HTML(http.StatusBadRequest, "admin/settings.html", gin.H{
				"Error": "New password is required",
			})
		} else {
			utils.Error(c, http.StatusBadRequest, "new password is required")
		}
		return
	}

	if len(newPassword) < 6 {
		if isFormData {
			c.HTML(http.StatusBadRequest, "admin/settings.html", gin.H{
				"Error": "Password must be at least 6 characters long",
			})
		} else {
			utils.Error(c, http.StatusBadRequest, "password must be at least 6 characters long")
		}
		return
	}

	// Get admin email from cookie
	email, err := c.Cookie("admin_token")
	if err != nil {
		if isFormData {
			c.Redirect(http.StatusSeeOther, "/admin/login")
		} else {
			utils.Error(c, http.StatusUnauthorized, "not authenticated")
		}
		return
	}

	// Find and verify current admin
	var admin models.AdminUser
	if err := h.db.Where("email = ?", email).First(&admin).Error; err != nil {
		if isFormData {
			c.HTML(http.StatusInternalServerError, "admin/settings.html", gin.H{
				"Error": "Admin not found",
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "admin not found")
		}
		return
	}

	// Verify current password if provided
	if currentPassword != "" {
		if !utils.CheckPassword(admin.PasswordHash, currentPassword) {
			if isFormData {
				c.HTML(http.StatusUnauthorized, "admin/settings.html", gin.H{
					"Error": "Current password is incorrect",
					"Admin": admin,
				})
			} else {
				utils.Error(c, http.StatusUnauthorized, "current password is incorrect")
			}
			return
		}
	}

	// Hash new password
	newPasswordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		if isFormData {
			c.HTML(http.StatusInternalServerError, "admin/settings.html", gin.H{
				"Error": "Failed to hash password",
				"Admin": admin,
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to hash password")
		}
		return
	}

	// Update password in database
	if err := h.db.Model(&admin).Update("password_hash", newPasswordHash).Error; err != nil {
		if isFormData {
			c.HTML(http.StatusInternalServerError, "admin/settings.html", gin.H{
				"Error": "Failed to update password",
				"Admin": admin,
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to update password")
		}
		return
	}

	h.log.Info("admin password changed", zap.String("admin_email", admin.Email))

	if isFormData {
		c.HTML(http.StatusOK, "admin/settings.html", gin.H{
			"Admin":   admin,
			"Success": "Password updated successfully",
		})
	} else {
		utils.Ok(c, "password updated successfully", nil)
	}
}

// ListUsers handles both HTML and API requests for user/rider management
func (h *AdminHandler) ListUsers(c *gin.Context) {
	// Check if it's an HTML request (browser) or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	// Get query parameters for filtering
	status := c.Query("status")
	page := c.DefaultQuery("page", "1")

	var users []models.Rider
	query := h.db.Order("created_at desc")

	// Apply status filter if provided
	if status == "blocked" {
		query = query.Where("is_blocked = ?", true)
	} else if status == "active" {
		query = query.Where("is_blocked = ?", false)
	}

	// For API requests, apply pagination
	if !isHTMLRequest {
		var offset int
		if pageNum := c.Query("page"); pageNum != "" {
			if p, err := strconv.Atoi(pageNum); err == nil && p > 0 {
				offset = (p - 1) * 50
			}
		}
		query = query.Offset(offset).Limit(50)
	} else {
		// For HTML requests, limit to reasonable number
		query = query.Limit(200)
	}

	if err := query.Find(&users).Error; err != nil {
		if isHTMLRequest {
			c.HTML(http.StatusInternalServerError, "admin/users.html", gin.H{
				"Error": "Failed to load users",
			})
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to list users")
		}
		return
	}

	if isHTMLRequest {
		// Calculate statistics
		var totalUsers, activeUsers, blockedUsers int64
		h.db.Model(&models.Rider{}).Count(&totalUsers)
		h.db.Model(&models.Rider{}).Where("is_blocked = ?", false).Count(&activeUsers)
		h.db.Model(&models.Rider{}).Where("is_blocked = ?", true).Count(&blockedUsers)

		// Calculate active rate
		var activeRate float64
		if totalUsers > 0 {
			activeRate = (float64(activeUsers) / float64(totalUsers)) * 100
		}

		// Render HTML template
		c.HTML(http.StatusOK, "admin/users.html", gin.H{
			"Users":        users,
			"FilterStatus": status,
			"CurrentPage":  page,
			"TotalUsers":   totalUsers,
			"ActiveUsers":  activeUsers,
			"BlockedUsers": blockedUsers,
			"ActiveRate":   activeRate,
		})
	} else {
		// Return JSON for API
		utils.Ok(c, "users", users)
	}
}

// BlockUser blocks or unblocks a user
func (h *AdminHandler) BlockUser(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	action := c.Query("action") // "block" or "unblock"
	
	if action != "block" && action != "unblock" {
		action = "block" // default action
	}

	isBlocked := action == "block"
	
	if err := h.db.Model(&models.Rider{}).Where("id = ?", id).Update("is_blocked", isBlocked).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/admin/users?error=failed_to_update")
		return
	}

	actionText := "blocked"
	if !isBlocked {
		actionText = "unblocked"
	}

	h.log.Info("user status updated", zap.String("user_id", id), zap.String("action", actionText))
	c.Redirect(http.StatusSeeOther, "/admin/users")
}

// ShowUser displays detailed information about a specific user
func (h *AdminHandler) ShowUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.HTML(http.StatusBadRequest, "admin/users.html", gin.H{
			"Error": "Invalid user ID",
		})
		return
	}

	var user models.Rider
	if err := h.db.First(&user, id).Error; err != nil {
		c.HTML(http.StatusNotFound, "admin/users.html", gin.H{
			"Error": "User not found",
		})
		return
	}

	// Get user's ride history
	var rides []models.Ride
	h.db.Where("rider_id = ?", user.ID).Preload("Driver").Order("created_at desc").Limit(10).Find(&rides)

	// Calculate account age
	accountAge := "Less than a day"
	if days := int(time.Since(user.CreatedAt).Hours() / 24); days > 0 {
		if days == 1 {
			accountAge = "1 day"
		} else {
			accountAge = fmt.Sprintf("%d days", days)
		}
	}

	// Check if it's an HTML request or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	if isHTMLRequest {
		c.HTML(http.StatusOK, "admin/user_detail.html", gin.H{
			"User":       user,
			"Rides":      rides,
			"AccountAge": accountAge,
		})
	} else {
		utils.Ok(c, "user", gin.H{
			"user":        user,
			"rides":       rides,
			"account_age": accountAge,
		})
	}
}

// ListComplaints handles complaints management (placeholder implementation)
func (h *AdminHandler) ListComplaints(c *gin.Context) {
	// Check if it's an HTML request (browser) or API request
	acceptHeader := c.GetHeader("Accept")
	isHTMLRequest := strings.Contains(acceptHeader, "text/html") || acceptHeader == ""

	if isHTMLRequest {
		// Render HTML template with placeholder data
		c.HTML(http.StatusOK, "admin/complaints.html", gin.H{
			"Message":     "Complaints management feature coming soon",
			"CurrentPage": "complaints",
		})
	} else {
		// Return JSON for API
		utils.Ok(c, "complaints", gin.H{
			"message": "Complaints management feature coming soon",
		})
	}
}
