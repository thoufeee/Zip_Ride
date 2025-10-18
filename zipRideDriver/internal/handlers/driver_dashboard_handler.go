package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DriverDashboardHandler struct {
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
	cfg *config.Config
}

func NewDriverDashboardHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *DriverDashboardHandler {
	return &DriverDashboardHandler{db: db, rdb: rdb, log: log, cfg: cfg}
}

func (h *DriverDashboardHandler) mustOwn(c *gin.Context) (uint, bool) {
	uidAny, _ := c.Get("uid")
	uid, _ := uidAny.(uint)
	idStr := strings.TrimSpace(c.Param("driverId"))
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || uid == 0 || uint(id64) != uid {
		utils.Error(c, http.StatusForbidden, "forbidden")
		return 0, false
	}
	return uid, true
}

func (h *DriverDashboardHandler) Profile(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	var d models.Driver
	if err := h.db.First(&d, uid).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "driver not found")
		return
	}
	var v models.Vehicle
	_ = h.db.Where("driver_id = ?", uid).Order("id desc").First(&v).Error
	d.PasswordHash = ""
	utils.Ok(c, "profile", gin.H{"driver": d, "vehicle": v})
}

type patchStatus struct{ IsOnline bool `json:"is_online" binding:"required"` }
func (h *DriverDashboardHandler) UpdateStatus(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	var req patchStatus
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.db.Model(&models.Driver{}).Where("id = ?", uid).Update("is_online", req.IsOnline).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to update status")
		return
	}
	utils.Ok(c, "status updated", gin.H{"is_online": req.IsOnline})
}

func (h *DriverDashboardHandler) EarningsSummary(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	var lifetime, last7, last30 float64
	now := time.Now().UTC()
	seven := now.Add(-7 * 24 * time.Hour)
	thirty := now.Add(-30 * 24 * time.Hour)
	_ = h.db.Model(&models.Earning{}).Where("driver_id = ?", uid).Select("COALESCE(SUM(amount),0)").Scan(&lifetime).Error
	_ = h.db.Model(&models.Earning{}).Where("driver_id = ? AND created_at >= ?", uid, seven).Select("COALESCE(SUM(amount),0)").Scan(&last7).Error
	_ = h.db.Model(&models.Earning{}).Where("driver_id = ? AND created_at >= ?", uid, thirty).Select("COALESCE(SUM(amount),0)").Scan(&last30).Error
	var ridesCount int64
	_ = h.db.Model(&models.Ride{}).Where("driver_id = ?", uid).Count(&ridesCount).Error
	utils.Ok(c, "earnings summary", gin.H{"lifetime": lifetime, "last7": last7, "last30": last30, "rides_count": ridesCount})
}

func (h *DriverDashboardHandler) EarningsTrend(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	type row struct {
		Day   time.Time `json:"day"`
		Total float64   `json:"total"`
	}
	var rows []row
	q := `SELECT date_trunc('day', created_at) AS day, COALESCE(SUM(amount),0) AS total
		FROM earnings
		WHERE driver_id = ? AND created_at >= NOW() - INTERVAL '14 days'
		GROUP BY day
		ORDER BY day`
	if err := h.db.Raw(q, uid).Scan(&rows).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to load trend")
		return
	}
	utils.Ok(c, "trend", rows)
}

func (h *DriverDashboardHandler) RidesSummary(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	type row struct {
		Status string `json:"status"`
		Count  int64  `json:"count"`
	}
	var rows []row
	q := `SELECT status, COUNT(*) AS count
		FROM rides
		WHERE driver_id = ? AND created_at >= NOW() - INTERVAL '30 days'
		GROUP BY status`
	if err := h.db.Raw(q, uid).Scan(&rows).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to load summary")
		return
	}
	utils.Ok(c, "rides summary", rows)
}

type withdrawReq struct{ Amount float64 `json:"amount" binding:"required,gt=0"` }
func (h *DriverDashboardHandler) Withdraw(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	var req withdrawReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	w := models.Withdrawal{DriverID: uid, Amount: req.Amount, Status: "pending"}
	if err := h.db.Create(&w).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to create withdrawal")
		return
	}
	utils.Ok(c, "withdrawal requested", w)
}

func (h *DriverDashboardHandler) Dashboard(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	var todayEarnings float64
	var completedToday int64
	_ = h.db.Model(&models.Earning{}).Where("driver_id = ? AND created_at >= ?", uid, today).Select("COALESCE(SUM(amount),0)").Scan(&todayEarnings).Error
	_ = h.db.Model(&models.Ride{}).Where("driver_id = ? AND status = ? AND created_at >= ?", uid, "completed", today).Count(&completedToday).Error
	var d models.Driver
	_ = h.db.First(&d, uid).Error
	var v models.Vehicle
	_ = h.db.Where("driver_id = ?", uid).Order("id desc").First(&v).Error
	utils.Ok(c, "dashboard", gin.H{"today_earnings": todayEarnings, "rides_completed_today": completedToday, "is_online": d.IsOnline, "vehicle": v})
}
