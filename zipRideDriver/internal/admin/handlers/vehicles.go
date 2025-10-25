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

type VehiclesHandler struct{ db *gorm.DB; log *zap.Logger }

func NewVehiclesHandler(db *gorm.DB, log *zap.Logger) *VehiclesHandler { return &VehiclesHandler{db: db, log: log} }

type vehicleView struct {
	Vehicle models.Vehicle
	Driver  models.Driver
}

// GET /admin/panel/vehicles
func (h *VehiclesHandler) Index(c *gin.Context) {
	var vs []models.Vehicle
	_ = h.db.Order("created_at desc").Limit(200).Find(&vs).Error
	// Load drivers for display
	driverMap := map[uint]models.Driver{}
	var result []vehicleView
	for _, v := range vs {
		if _, ok := driverMap[v.DriverID]; !ok {
			var d models.Driver
			_ = h.db.First(&d, v.DriverID).Error
			driverMap[v.DriverID] = d
		}
		result = append(result, vehicleView{Vehicle: v, Driver: driverMap[v.DriverID]})
	}
	c.HTML(http.StatusOK, "admin/vehicles/index.html", gin.H{"Vehicles": result})
}

// GET /admin/panel/vehicles/:id
func (h *VehiclesHandler) Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var v models.Vehicle
	if err := h.db.First(&v, id).Error; err != nil { c.String(http.StatusNotFound, "Vehicle not found"); return }
	var d models.Driver
	_ = h.db.First(&d, v.DriverID).Error
	c.HTML(http.StatusOK, "admin/vehicles/show.html", gin.H{"Vehicle": v, "Driver": d})
}

// POST /admin/panel/vehicles/:id/verify
func (h *VehiclesHandler) Verify(c *gin.Context) {
	id := c.Param("id")
	_ = h.db.Model(&models.Vehicle{}).Where("id = ?", id).Update("status", "verified").Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/vehicles")
}

// POST /admin/panel/vehicles/:id/assign (form: driver_id)
func (h *VehiclesHandler) Assign(c *gin.Context) {
	id := c.Param("id")
	driverID := strings.TrimSpace(c.PostForm("driver_id"))
	if driverID == "" { c.String(http.StatusBadRequest, "driver_id required"); return }
	_ = h.db.Model(&models.Vehicle{}).Where("id = ?", id).Update("driver_id", driverID).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/vehicles/"+id)
}

// POST /admin/panel/vehicles/:id/deactivate
func (h *VehiclesHandler) Deactivate(c *gin.Context) {
	id := c.Param("id")
	_ = h.db.Model(&models.Vehicle{}).Where("id = ?", id).Update("status", "inactive").Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/vehicles")
}
