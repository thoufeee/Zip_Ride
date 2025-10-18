package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/middleware"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func RegisterVehicleRoutes(r *gin.Engine, cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) {
	g := r.Group("/api/driver")
	g.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	mustOwn := func(c *gin.Context) (uint, bool) {
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

	g.GET("/:driverId/vehicles", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		var list []models.Vehicle
		if err := db.Where("driver_id = ?", uid).Order("id desc").Find(&list).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to load vehicles")
			return
		}
		utils.Ok(c, "vehicles", list)
	})

	type vehicleReq struct {
		Make            string `json:"make" binding:"required"`
		Model           string `json:"model" binding:"required"`
		Year            int    `json:"year" binding:"required"`
		PlateNumber     string `json:"plate_number" binding:"required"`
		InsuranceNumber string `json:"insurance_number" binding:"required"`
		RCNumber        string `json:"rc_number" binding:"required"`
	}

	g.POST("/:driverId/vehicles", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		var req vehicleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		v := models.Vehicle{
			DriverID:        uid,
			Make:            strings.TrimSpace(req.Make),
			Model:           strings.TrimSpace(req.Model),
			Year:            req.Year,
			PlateNumber:     strings.TrimSpace(req.PlateNumber),
			InsuranceNumber: strings.TrimSpace(req.InsuranceNumber),
			RCNumber:        strings.TrimSpace(req.RCNumber),
			Status:          "active",
		}
		if err := db.Create(&v).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to create vehicle")
			return
		}
		utils.Ok(c, "created", v)
	})

	g.PUT("/:driverId/vehicles/:id", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		var req vehicleReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		updates := map[string]interface{}{
			"make":             strings.TrimSpace(req.Make),
			"model":            strings.TrimSpace(req.Model),
			"year":             req.Year,
			"plate_number":     strings.TrimSpace(req.PlateNumber),
			"insurance_number": strings.TrimSpace(req.InsuranceNumber),
			"rc_number":        strings.TrimSpace(req.RCNumber),
			"updated_at":       time.Now().UTC(),
		}
		if err := db.Model(&models.Vehicle{}).Where("id = ? AND driver_id = ?", id, uid).Updates(updates).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to update vehicle")
			return
		}
		utils.Ok(c, "updated", nil)
	})

	g.DELETE("/:driverId/vehicles/:id", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		if err := db.Where("id = ? AND driver_id = ?", id, uid).Delete(&models.Vehicle{}).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to delete vehicle")
			return
		}
		utils.Ok(c, "deleted", nil)
	})

	type docReq struct {
		DocType string `json:"doc_type" binding:"required"`
		DocURL  string `json:"doc_url" binding:"required"`
	}

	g.GET("/:driverId/documents", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		var docs []models.DriverDocument
		if err := db.Where("driver_id = ?", uid).Order("uploaded_at desc").Find(&docs).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to load documents")
			return
		}
		utils.Ok(c, "documents", docs)
	})

	g.POST("/:driverId/documents", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		var req docReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		d := models.DriverDocument{DriverID: uid, DocType: strings.TrimSpace(req.DocType), DocURL: strings.TrimSpace(req.DocURL), UploadedAt: time.Now().UTC()}
		if err := db.Create(&d).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to upload document")
			return
		}
		utils.Ok(c, "uploaded", d)
	})

	type docStatusReq struct{ Verified bool `json:"verified" binding:"required"` }
	g.PATCH("/:driverId/documents/:id/status", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		var req docStatusReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		updates := map[string]interface{}{"verified": req.Verified}
		if req.Verified { updates["verified_at"] = time.Now().UTC() }
		if err := db.Model(&models.DriverDocument{}).Where("id = ? AND driver_id = ?", id, uid).Updates(updates).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to update document")
			return
		}
		utils.Ok(c, "document updated", nil)
	})
}
