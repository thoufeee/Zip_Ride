package adminhandlers

import (
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"zipRideDriver/internal/models"
)

type DashboardHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewDashboardHandler(db *gorm.DB, log *zap.Logger) *DashboardHandler {
	return &DashboardHandler{db: db, log: log}
}

func (h *DashboardHandler) DashboardPage(c *gin.Context) {
	var totalDrivers, activeDrivers, pendingDrivers int64
	var totalRides, completedRides, ongoingRides int64
	var totalEarnings, todayEarnings, weeklyEarnings float64

	// Driver statistics
	h.db.Model(&models.Driver{}).Count(&totalDrivers)
	h.db.Model(&models.Driver{}).Where("is_online = ?", true).Count(&activeDrivers)
	h.db.Model(&models.Driver{}).Where("status = ?", "Pending").Count(&pendingDrivers)
	
	// Ride statistics
	h.db.Model(&models.Ride{}).Count(&totalRides)
	h.db.Model(&models.Ride{}).Where("status = ?", "completed").Count(&completedRides)
	h.db.Model(&models.Ride{}).Where("status = ?", "ongoing").Count(&ongoingRides)
	
	// Earnings statistics
	h.db.Model(&models.Earning{}).Select("COALESCE(SUM(amount), 0)").Scan(&totalEarnings)
	h.db.Model(&models.Earning{}).Where("DATE(created_at) = CURRENT_DATE").Select("COALESCE(SUM(amount), 0)").Scan(&todayEarnings)
	h.db.Model(&models.Earning{}).Where("created_at >= DATE_SUB(NOW(), INTERVAL 7 DAY)").Select("COALESCE(SUM(amount), 0)").Scan(&weeklyEarnings)

	// Fetch pending withdrawals count
	var pendingWithdrawals int64
	h.db.Model(&models.Withdrawal{}).Where("status = ?", "pending").Count(&pendingWithdrawals)

	// Fetch open tickets count
	var openTickets int64
	h.db.Model(&models.HelpTicket{}).Where("status IN ?", []string{"open", "pending"}).Count(&openTickets)

	// Fetch recent drivers (last 5)
	var recentDrivers []models.Driver
	h.db.Order("created_at desc").Limit(5).Find(&recentDrivers)

	// Fetch recent rides (last 10)
	var recentRides []models.Ride
	h.db.Preload("Driver").Preload("Rider").Order("created_at desc").Limit(10).Find(&recentRides)
	
	// Weekly ride trend data (for charts)
	type DayData struct {
		Day   string `json:"day"`
		Count int    `json:"count"`
	}
	var weeklyTrend []DayData
	// This would be populated with actual data from DB

	labels := make([]string, len(weeklyTrend))
	counts := make([]int, len(weeklyTrend))
	for i, d := range weeklyTrend {
		labels[i] = d.Day
		counts[i] = d.Count
	}
	labelsJSON, _ := json.Marshal(labels)
	countsJSON, _ := json.Marshal(counts)

	c.HTML(http.StatusOK, "admin/dashboard.html", gin.H{
		"Totals": gin.H{
			"TotalDrivers":   totalDrivers,
			"ActiveDrivers":  activeDrivers,
			"PendingDrivers": pendingDrivers,
			"TotalRides":     totalRides,
			"CompletedRides": completedRides,
			"OngoingRides":   ongoingRides,
			"TotalEarnings":  totalEarnings,
			"TodayEarnings":  todayEarnings,
			"WeeklyEarnings": weeklyEarnings,
		},
		"PendingWithdrawals": pendingWithdrawals,
		"OpenTickets":        openTickets,
		"RecentDrivers":      recentDrivers,
		"RecentRides":        recentRides,
		"WeeklyTrend":        weeklyTrend,
		"WeeklyTrendLabelsJSON": template.JS(string(labelsJSON)),
		"WeeklyTrendCountsJSON": template.JS(string(countsJSON)),
	})
}
