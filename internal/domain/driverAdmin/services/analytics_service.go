package services

import (
	"net/http"
	"time"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// GetDriverAnalytics returns comprehensive driver analytics
func GetDriverAnalytics(c *gin.Context) {
	// Date range filters
	startDate := c.DefaultQuery("start_date", time.Now().AddDate(0, 0, -30).Format("2006-01-02"))
	endDate := c.DefaultQuery("end_date", time.Now().Format("2006-01-02"))

	// Parse dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_date format"})
		return
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_date format"})
		return
	}

	// Total drivers
	var totalDrivers int64
	database.DB.Model(&models.Driver{}).Count(&totalDrivers)

	// Drivers by status
	var statusStats []gin.H
	database.DB.Model(&models.Driver{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusStats)

	// Registration trends (daily)
	var registrationTrends []gin.H
	database.DB.Model(&models.Driver{}).
		Select("DATE(created_at) as date, count(*) as count").
		Where("created_at BETWEEN ? AND ?", start, end.Add(24*time.Hour)).
		Group("DATE(created_at)").
		Order("date").
		Scan(&registrationTrends)

	// Phone verification rate
	var verifiedCount int64
	database.DB.Model(&models.Driver{}).Where("phone_verified = ?", true).Count(&verifiedCount)
	verificationRate := float64(verifiedCount) / float64(totalDrivers) * 100

	// Document completion rate
	var docsCompleted int64
	database.DB.Model(&models.DriverDocuments{}).
		Where("license_url != '' AND rc_url != '' AND insurance_url != ''").
		Count(&docsCompleted)
	docCompletionRate := float64(docsCompleted) / float64(totalDrivers) * 100

	// Recent activity (last 7 days)
	var recentActivity []gin.H
	database.DB.Model(&models.Driver{}).
		Select("DATE(created_at) as date, count(*) as registrations").
		Where("created_at >= ?", time.Now().AddDate(0, 0, -7)).
		Group("DATE(created_at)").
		Order("date").
		Scan(&recentActivity)

	// Top performing statuses
	var topStatuses []gin.H
	database.DB.Model(&models.Driver{}).
		Select("status, count(*) as count").
		Group("status").
		Order("count DESC").
		Limit(5).
		Scan(&topStatuses)

	analytics := gin.H{
		"overview": gin.H{
			"total_drivers":        totalDrivers,
			"phone_verified":       verifiedCount,
			"verification_rate":    verificationRate,
			"docs_completed":       docsCompleted,
			"doc_completion_rate":  docCompletionRate,
		},
		"status_breakdown":    statusStats,
		"registration_trends": registrationTrends,
		"recent_activity":     recentActivity,
		"top_statuses":        topStatuses,
		"date_range": gin.H{
			"start_date": startDate,
			"end_date":   endDate,
		},
	}

	c.JSON(http.StatusOK, analytics)
}

// GetDriverPerformanceMetrics returns performance metrics for drivers
func GetDriverPerformanceMetrics(c *gin.Context) {
	// This would typically include ride metrics, ratings, etc.
	// For now, we'll return basic driver status metrics

	var metrics gin.H

	// Approval rate
	var totalApplications int64
	var approvedCount int64
	database.DB.Model(&models.Driver{}).Count(&totalApplications)
	database.DB.Model(&models.Driver{}).Where("status = ?", "approved").Count(&approvedCount)
	
	approvalRate := float64(0)
	if totalApplications > 0 {
		approvalRate = float64(approvedCount) / float64(totalApplications) * 100
	}

	// Average time to approval (for approved drivers)
	var avgApprovalTime float64
	database.DB.Model(&models.Driver{}).
		Select("AVG(EXTRACT(EPOCH FROM (updated_at - created_at))/3600) as avg_hours").
		Where("status = ?", "approved").
		Scan(&avgApprovalTime)

	// Rejection reasons (if we had a rejection_reason field)
	var rejectionCount int64
	database.DB.Model(&models.Driver{}).Where("status = ?", "rejected").Count(&rejectionCount)

	// Suspension rate
	var suspensionCount int64
	database.DB.Model(&models.Driver{}).Where("status = ?", "suspended").Count(&suspensionCount)

	metrics = gin.H{
		"approval_rate":      approvalRate,
		"avg_approval_time":  avgApprovalTime, // hours
		"total_applications": totalApplications,
		"approved_count":     approvedCount,
		"rejection_count":    rejectionCount,
		"suspension_count":   suspensionCount,
		"rejection_rate":     float64(rejectionCount) / float64(totalApplications) * 100,
		"suspension_rate":    float64(suspensionCount) / float64(totalApplications) * 100,
	}

	c.JSON(http.StatusOK, metrics)
}

// GetDocumentAnalytics returns document-related analytics
func GetDocumentAnalytics(c *gin.Context) {
	var analytics gin.H

	// Document upload status
	var docStats []gin.H
	database.DB.Model(&models.DriverDocuments{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&docStats)

	// Document completion by type
	var licenseCount, rcCount, insuranceCount int64
	database.DB.Model(&models.DriverDocuments{}).Where("license_url != ''").Count(&licenseCount)
	database.DB.Model(&models.DriverDocuments{}).Where("rc_url != ''").Count(&rcCount)
	database.DB.Model(&models.DriverDocuments{}).Where("insurance_url != ''").Count(&insuranceCount)

	// Total drivers for percentage calculation
	var totalDrivers int64
	database.DB.Model(&models.Driver{}).Count(&totalDrivers)

	analytics = gin.H{
		"document_status": docStats,
		"completion_by_type": gin.H{
			"license": gin.H{
				"count":      licenseCount,
				"percentage": float64(licenseCount) / float64(totalDrivers) * 100,
			},
			"rc": gin.H{
				"count":      rcCount,
				"percentage": float64(rcCount) / float64(totalDrivers) * 100,
			},
			"insurance": gin.H{
				"count":      insuranceCount,
				"percentage": float64(insuranceCount) / float64(totalDrivers) * 100,
			},
		},
		"total_drivers": totalDrivers,
	}

	c.JSON(http.StatusOK, analytics)
}

// GetSystemHealthMetrics returns system health and performance metrics
func GetSystemHealthMetrics(c *gin.Context) {
	var health gin.H

	// Active drivers (online)
	var activeDrivers int64
	// This would typically check Redis for online drivers
	// For now, we'll use approved drivers as a proxy
	database.DB.Model(&models.Driver{}).Where("status = ?", "approved").Count(&activeDrivers)

	// Pending reviews
	var pendingReviews int64
	database.DB.Model(&models.Driver{}).Where("status = ?", "in_review").Count(&pendingReviews)

	// Recent registrations (last 24 hours)
	var recentRegistrations int64
	database.DB.Model(&models.Driver{}).
		Where("created_at >= ?", time.Now().Add(-24*time.Hour)).
		Count(&recentRegistrations)

	// Admin activity (if we had audit logs)
	var totalAdmins int64
	database.DB.Model(&models.DriverAdmin{}).Where("is_active = ?", true).Count(&totalAdmins)

	health = gin.H{
		"active_drivers":        activeDrivers,
		"pending_reviews":       pendingReviews,
		"recent_registrations":  recentRegistrations,
		"active_admins":         totalAdmins,
		"system_status":         "healthy", // This could be more sophisticated
		"last_updated":          time.Now(),
	}

	c.JSON(http.StatusOK, health)
}
