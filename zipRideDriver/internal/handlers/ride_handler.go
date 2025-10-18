package handlers

import (
	"context"
	"encoding/json"
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

func RegisterRideRoutes(r *gin.Engine, cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) {
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

	type locReq struct {
		Lat     float64 `json:"lat" binding:"required"`
		Lng     float64 `json:"lng" binding:"required"`
		Heading *float64 `json:"heading"`
	}
	g.POST("/:driverId/location", func(c *gin.Context) {
		uid, ok := mustOwn(c)
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
		_ = rdb.HSet(ctx, key, val).Err()
		_ = rdb.Expire(ctx, key, 5*time.Minute).Err()
		evt := map[string]interface{}{"type": "location_update", "driver_id": uid, "payload": val}
		if b, err := json.Marshal(evt); err == nil {
			_ = rdb.Publish(ctx, "driver:events:"+strconv.Itoa(int(uid)), string(b)).Err()
		}
		utils.Ok(c, "location updated", nil)
	})

	g.GET("/:driverId/requests", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		var list []models.Ride
		if err := db.Where("driver_id = ? AND status IN ?", uid, []string{"requested", "assigned"}).Order("id desc").Limit(20).Find(&list).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to load requests")
			return
		}
		utils.Ok(c, "requests", list)
	})

	g.POST("/:driverId/rides/:id/accept", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		var ride models.Ride
		if err := db.Where("id = ? AND driver_id = ?", id, uid).First(&ride).Error; err != nil {
			utils.Error(c, http.StatusNotFound, "ride not found")
			return
		}
		start := time.Now().UTC()
		if err := db.Model(&models.Ride{}).Where("id = ?", ride.ID).Updates(map[string]interface{}{"status": "accepted", "started_at": &start}).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to accept ride")
			return
		}
		ctx := context.Background()
		evt := map[string]interface{}{"type": "ride_accepted", "driver_id": uid, "ride_id": ride.ID, "ts": time.Now().UTC()}
		if b, err := json.Marshal(evt); err == nil {
			_ = rdb.Publish(ctx, "driver:events:"+strconv.Itoa(int(uid)), string(b)).Err()
		}
		utils.Ok(c, "accepted", nil)
	})

	g.POST("/:driverId/rides/:id/cancel", func(c *gin.Context) {
		uid, ok := mustOwn(c)
		if !ok {
			return
		}
		id, _ := strconv.Atoi(c.Param("id"))
		var ride models.Ride
		if err := db.Where("id = ? AND driver_id = ?", id, uid).First(&ride).Error; err != nil {
			utils.Error(c, http.StatusNotFound, "ride not found")
			return
		}
		if err := db.Model(&models.Ride{}).Where("id = ?", ride.ID).Update("status", "cancelled").Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to cancel ride")
			return
		}
		ctx := context.Background()
		evt := map[string]interface{}{"type": "ride_cancelled", "driver_id": uid, "ride_id": ride.ID, "ts": time.Now().UTC()}
		if b, err := json.Marshal(evt); err == nil {
			_ = rdb.Publish(ctx, "driver:events:"+strconv.Itoa(int(uid)), string(b)).Err()
		}
		utils.Ok(c, "cancelled", nil)
	})
}
