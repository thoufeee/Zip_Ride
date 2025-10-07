package services

import (
	"fmt"
	"net/http"
	"time"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// SetDriverAvailability toggles driver online/offline status
func SetDriverAvailability(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)

	var req struct {
		Available bool `json:"available" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if driver is approved
	var d models.Driver
	if err := database.DB.First(&d, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}
	if d.Status != "approved" || !d.PhoneVerified {
		c.JSON(http.StatusForbidden, gin.H{"error": "driver not approved or phone not verified"})
		return
	}

	// Store availability in Redis with TTL
	key := fmt.Sprintf("driver:available:%d", uid)
	if req.Available {
		database.RDB.Set(database.Ctx, key, "1", 0) // No expiration when online
	} else {
		database.RDB.Del(database.Ctx, key)
	}

	c.JSON(http.StatusOK, gin.H{"available": req.Available, "message": "status updated"})
}

// UpdateDriverLocation updates driver's current location
func UpdateDriverLocation(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)

	var req struct {
		Latitude  float64 `json:"latitude" binding:"required"`
		Longitude float64 `json:"longitude" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if driver is available
	key := fmt.Sprintf("driver:available:%d", uid)
	exists, err := database.RDB.Exists(database.Ctx, key).Result()
	if err != nil || exists == 0 {
		c.JSON(http.StatusForbidden, gin.H{"error": "driver must be online to update location"})
		return
	}

	// Store location in Redis Geo
	geoKey := "drivers:locations"
	database.RDB.GeoAdd(database.Ctx, geoKey, &redis.GeoLocation{
		Name:      fmt.Sprintf("driver:%d", uid),
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	})

	// Update last seen timestamp
	lastSeenKey := fmt.Sprintf("driver:last_seen:%d", uid)
	database.RDB.Set(database.Ctx, lastSeenKey, time.Now().Unix(), 30*time.Minute)

	c.JSON(http.StatusOK, gin.H{"message": "location updated"})
}

// GetNearbyDrivers finds drivers within radius of given location
func GetNearbyDrivers(c *gin.Context) {
	latStr := c.Query("lat")
	lngStr := c.Query("lng")
	radiusStr := c.DefaultQuery("radius", "5") // km

	if latStr == "" || lngStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "lat and lng required"})
		return
	}

	// Convert strings to float64
	var lat, lng, radius float64
	if _, err := fmt.Sscanf(latStr, "%f", &lat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid latitude"})
		return
	}
	if _, err := fmt.Sscanf(lngStr, "%f", &lng); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid longitude"})
		return
	}
	if _, err := fmt.Sscanf(radiusStr, "%f", &radius); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid radius"})
		return
	}

	geoKey := "drivers:locations"
	locations, err := database.RDB.GeoRadius(database.Ctx, geoKey, lng, lat, &redis.GeoRadiusQuery{
		Radius: radius,
		Unit:   "km",
	}).Result()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to find nearby drivers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"drivers": locations})
}

// GetDriverStatus returns current availability and location
func GetDriverStatus(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)

	// Check availability
	availableKey := fmt.Sprintf("driver:available:%d", uid)
	available, _ := database.RDB.Exists(database.Ctx, availableKey).Result()

	// Get last location
	geoKey := "drivers:locations"
	locations, _ := database.RDB.GeoPos(database.Ctx, geoKey, fmt.Sprintf("driver:%d", uid)).Result()

	status := gin.H{
		"available": available > 0,
		"location":  nil,
	}

	if len(locations) > 0 && locations[0] != nil {
		status["location"] = gin.H{
			"latitude":  locations[0].Latitude,
			"longitude": locations[0].Longitude,
		}
	}

	c.JSON(http.StatusOK, status)
}
