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

type DriverHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewDriverHandler(db *gorm.DB, log *zap.Logger) *DriverHandler {
	return &DriverHandler{db: db, log: log}
}

func (h *DriverHandler) DriversPage(c *gin.Context) {
	status := strings.TrimSpace(c.Query("status"))
	var list []models.Driver
	q := h.db.Order("created_at desc").Limit(200)
	if status != "" {
		q = q.Where("status = ?", status)
	}
	_ = q.Find(&list).Error
	c.HTML(http.StatusOK, "admin/drivers.html", gin.H{
		"Drivers": list,
		"FilterStatus": status,
	})
}

func (h *DriverHandler) DriverDetailPage(c *gin.Context) {
	idStr := strings.TrimSpace(c.Param("id"))
	id, _ := strconv.Atoi(idStr)
	var d models.Driver
	if err := h.db.First(&d, id).Error; err != nil {
		c.String(http.StatusNotFound, "Driver not found")
		return
	}
	var v models.Vehicle
	_ = h.db.Where("driver_id = ?", d.ID).Order("id desc").First(&v).Error
	var docs []models.DriverDocument
	_ = h.db.Where("driver_id = ?", d.ID).Order("uploaded_at desc").Find(&docs).Error
	c.HTML(http.StatusOK, "admin/drivers/show.html", gin.H{
		"Driver": d,
		"Vehicle": v,
		"Documents": docs,
	})
}

// PendingApprovals shows drivers waiting for approval
func (h *DriverHandler) PendingApprovals(c *gin.Context) {
	var drivers []models.Driver
	_ = h.db.Where("status = ?", "Pending").Order("created_at asc").Find(&drivers).Error
	
	// Count statistics
	var pendingCount, approvedCount, rejectedCount int64
	h.db.Model(&models.Driver{}).Where("status = ?", "Pending").Count(&pendingCount)
	h.db.Model(&models.Driver{}).Where("status = ?", "Approved").Count(&approvedCount)
	h.db.Model(&models.Driver{}).Where("status = ?", "Rejected").Count(&rejectedCount)
	
	c.HTML(http.StatusOK, "admin/drivers/pending.html", gin.H{
		"Drivers": drivers,
		"PendingCount": pendingCount,
		"ApprovedCount": approvedCount,
		"RejectedCount": rejectedCount,
	})
}

func (h *DriverHandler) ApproveDriver(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	
	// Update status and mark as verified
	updates := map[string]interface{}{
		"status": "Approved",
		"is_verified": true,
	}
	
	if err := h.db.Model(&models.Driver{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to approve")
		return
	}
	
	// Log the approval
	h.log.Info("Driver approved", zap.String("driver_id", id))
	
	// Check if this is an AJAX request
	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.JSON(http.StatusOK, gin.H{"message": "Driver approved successfully"})
		return
	}
	
	c.Redirect(http.StatusSeeOther, "/admin/panel/drivers/pending")
}

func (h *DriverHandler) RejectDriver(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	
	if err := h.db.Model(&models.Driver{}).Where("id = ?", id).Update("status", "Rejected").Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to reject")
		return
	}
	
	// Log the rejection
	h.log.Info("Driver rejected", zap.String("driver_id", id))
	
	// Check if this is an AJAX request
	if c.GetHeader("X-Requested-With") == "XMLHttpRequest" {
		c.JSON(http.StatusOK, gin.H{"message": "Driver rejected"})
		return
	}
	
	c.Redirect(http.StatusSeeOther, "/admin/panel/drivers/pending")
}

func (h *DriverHandler) SuspendDriver(c *gin.Context) {
	id := strings.TrimSpace(c.Param("id"))
	if err := h.db.Model(&models.Driver{}).Where("id = ?", id).Update("status", "Suspended").Error; err != nil {
		c.String(http.StatusInternalServerError, "Failed to suspend")
		return
	}
	c.Redirect(http.StatusSeeOther, "/admin/panel/drivers")
}
