package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RideHandler struct {
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewRideHandler(cfg interface{}, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *RideHandler {
	return &RideHandler{db: db, rdb: rdb, log: log}
}

func (h *RideHandler) mustOwn(c *gin.Context) (uint, bool) {
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

type locReq struct {
	Lat     float64 `json:"lat" binding:"required"`
	Lng     float64 `json:"lng" binding:"required"`
	Heading *float64 `json:"heading"`
}

func (h *RideHandler) UpdateLocation(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	var req locReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	key := "driver:loc:" + strconv.Itoa(int(uid))
	val := map[string]interface{}{"lat": req.Lat, "lng": req.Lng, "heading": req.Heading, "ts": time.Now().UTC()}
	ctx := c.Request.Context()
	_ = h.rdb.HSet(ctx, key, val).Err()
	_ = h.rdb.Expire(ctx, key, 5*time.Minute).Err()
	evt := map[string]interface{}{"type": "location_update", "driver_id": uid, "payload": val}
	if b, err := json.Marshal(evt); err == nil {
		_ = h.rdb.Publish(ctx, "driver:events:"+strconv.Itoa(int(uid)), string(b)).Err()
	}
	utils.Ok(c, "location updated", nil)
}

func (h *RideHandler) ListRequests(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	var list []models.Ride
	if err := h.db.Where("driver_id = ? AND status IN ?", uid, []string{"requested", "assigned"}).Order("id desc").Limit(20).Find(&list).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to load requests")
		return
	}
	utils.Ok(c, "requests", list)
}

func (h *RideHandler) AcceptRide(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var ride models.Ride
	if err := h.db.Where("id = ? AND driver_id = ?", id, uid).First(&ride).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "ride not found")
		return
	}
	start := time.Now().UTC()
	if err := h.db.Model(&models.Ride{}).Where("id = ?", ride.ID).Updates(map[string]interface{}{"status": "accepted", "started_at": &start}).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to accept ride")
		return
	}
	ctx := context.Background()
	evt := map[string]interface{}{"type": "ride_accepted", "driver_id": uid, "ride_id": ride.ID, "ts": time.Now().UTC()}
	if b, err := json.Marshal(evt); err == nil {
		_ = h.rdb.Publish(ctx, "driver:events:"+strconv.Itoa(int(uid)), string(b)).Err()
	}
	utils.Ok(c, "accepted", nil)
}

func (h *RideHandler) CancelRide(c *gin.Context) {
	uid, ok := h.mustOwn(c)
	if !ok {
		return
	}
	id, _ := strconv.Atoi(c.Param("id"))
	var ride models.Ride
	if err := h.db.Where("id = ? AND driver_id = ?", id, uid).First(&ride).Error; err != nil {
		utils.Error(c, http.StatusNotFound, "ride not found")
		return
	}
	if err := h.db.Model(&models.Ride{}).Where("id = ?", ride.ID).Update("status", "cancelled").Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to cancel ride")
		return
	}
	ctx := context.Background()
	evt := map[string]interface{}{"type": "ride_cancelled", "driver_id": uid, "ride_id": ride.ID, "ts": time.Now().UTC()}
	if b, err := json.Marshal(evt); err == nil {
		_ = h.rdb.Publish(ctx, "driver:events:"+strconv.Itoa(int(uid)), string(b)).Err()
	}
	utils.Ok(c, "cancelled", nil)
}
