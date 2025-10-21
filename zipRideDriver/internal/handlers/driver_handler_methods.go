package handlers

import (
	"net/http"
	"strings"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type DriverHandler struct {
	cfg *config.Config
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewDriverHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *DriverHandler {
	return &DriverHandler{cfg: cfg, db: db, rdb: rdb, log: log}
}

func (h *DriverHandler) Onboarding(c *gin.Context) {
	uidAny, exists := c.Get("uid")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "unauthorized")
		return
	}
	uid, _ := uidAny.(uint)
	var req onboardingReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	updates := map[string]interface{}{
		"name":           strings.TrimSpace(req.Name),
		"email":          strings.TrimSpace(req.Email),
		"license_number": strings.TrimSpace(req.LicenseNumber),
	}
	if err := h.db.Model(&models.Driver{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to update profile")
		return
	}
	v := models.Vehicle{
		DriverID:        uid,
		Make:            strings.TrimSpace(req.VehicleMake),
		Model:           strings.TrimSpace(req.VehicleModel),
		Year:            req.VehicleYear,
		PlateNumber:     strings.TrimSpace(req.PlateNumber),
		InsuranceNumber: strings.TrimSpace(req.InsuranceNumber),
		RCNumber:        strings.TrimSpace(req.RCNumber),
		Status:          "active",
	}
	if err := h.db.Create(&v).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to save vehicle")
		return
	}
	utils.Ok(c, "onboarding complete", nil)
}
