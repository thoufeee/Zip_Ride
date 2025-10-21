package adminhandlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zipRideDriver/internal/models"
)

type RidesHandler struct{ db *gorm.DB; log *zap.Logger }

func NewRidesHandler(db *gorm.DB, log *zap.Logger) *RidesHandler { return &RidesHandler{db: db, log: log} }

// GET /admin/panel/rides?status=
func (h *RidesHandler) Index(c *gin.Context) {
	status := c.Query("status")
	var rides []models.Ride
	q := h.db.Order("created_at desc").Limit(200)
	if status != "" { q = q.Where("status = ?", status) }
	_ = q.Find(&rides).Error
	c.HTML(http.StatusOK, "admin/rides/index.html", gin.H{"Rides": rides, "FilterStatus": status})
}

// GET /admin/panel/rides/:id
func (h *RidesHandler) Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var r models.Ride
	if err := h.db.First(&r, id).Error; err != nil { c.String(http.StatusNotFound, "Ride not found"); return }
	var d models.Driver
	_ = h.db.First(&d, r.DriverID).Error
	c.HTML(http.StatusOK, "admin/rides/show.html", gin.H{"Ride": r, "Driver": d})
}

// POST /admin/panel/rides/:id/cancel
func (h *RidesHandler) Cancel(c *gin.Context) {
	id := c.Param("id")
	now := time.Now().UTC()
	_ = h.db.Model(&models.Ride{}).Where("id = ?", id).Updates(map[string]interface{}{"status": "cancelled", "EndedAt": &now}).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/rides")
}
