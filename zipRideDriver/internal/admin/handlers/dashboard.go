package adminhandlers

import (
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
	var totals struct {
		TotalDrivers    int64
		ActiveDrivers   int64
		PendingDrivers  int64
		TotalRides      int64
		TotalEarnings   float64
	}
	_ = h.db.Model(&models.Driver{}).Count(&totals.TotalDrivers).Error
	_ = h.db.Model(&models.Driver{}).Where("status = ?", "Approved").Count(&totals.ActiveDrivers).Error
	_ = h.db.Model(&models.Driver{}).Where("status = ?", "Pending").Count(&totals.PendingDrivers).Error
	_ = h.db.Model(&models.Ride{}).Count(&totals.TotalRides).Error
	_ = h.db.Model(&models.Earning{}).Select("COALESCE(SUM(amount),0)").Scan(&totals.TotalEarnings).Error

	c.HTML(http.StatusOK, "admin/dashboard.html", gin.H{
		"Totals": totals,
	})
}
