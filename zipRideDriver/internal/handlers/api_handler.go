package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type APIHandler struct {
	cfg *config.Config
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewAPIHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *APIHandler {
	return &APIHandler{
		cfg: cfg,
		db:  db,
		rdb: rdb,
		log: log,
	}
}

// GET /api/vehicles - Fetch all vehicles (requires JWT)
func (h *APIHandler) GetAllVehicles(c *gin.Context) {
	var vehicles []models.Vehicle

	if err := h.db.Order("created_at desc").Find(&vehicles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to fetch vehicles",
		})
		return
	}

	// Transform to match expected response format
	var response []gin.H
	for _, v := range vehicles {
		response = append(response, gin.H{
			"id":       v.ID,
			"name":     v.Make + " " + v.Model,
			"number":   v.PlateNumber,
			"type":     "car", // Default type since VehicleType doesn't exist in model
			"isActive": v.Status == "active",
			"imageUrl": "/uploads/vehicles/" + v.PlateNumber + ".jpg",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "success",
		"vehicles": response,
	})
}

// POST /api/vehicles - Add new vehicle (requires JWT)
func (h *APIHandler) CreateVehicle(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Number   string `json:"number" binding:"required"`
		Type     string `json:"type" binding:"required"`
		ImageURL string `json:"imageUrl"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
		})
		return
	}

	// Extract make and model from name
	nameParts := strings.SplitN(req.Name, " ", 2)
	make := nameParts[0]
	model := ""
	if len(nameParts) > 1 {
		model = nameParts[1]
	}

	// Get driver ID from JWT context
	userID, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	// Convert userID to uint
	var driverID uint
	switch v := userID.(type) {
	case uint:
		driverID = v
	case float64:
		driverID = uint(v)
	case string:
		if id, err := strconv.ParseUint(v, 10, 32); err == nil {
			driverID = uint(id)
		} else {
			driverID = 1 // fallback
		}
	default:
		driverID = 1 // fallback
	}

	vehicle := models.Vehicle{
		Make:        make,
		Model:       model,
		PlateNumber: req.Number,
		Status:      "active",
		DriverID:    driverID,
	}

	if err := h.db.Create(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to create vehicle",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Vehicle added successfully",
		"vehicle": gin.H{
			"id":       vehicle.ID,
			"name":     req.Name,
			"number":   req.Number,
			"type":     req.Type,
			"imageUrl": req.ImageURL,
		},
	})
}

// PUT /api/vehicles/:id - Update vehicle (requires JWT)
func (h *APIHandler) UpdateVehicle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid vehicle ID",
		})
		return
	}

	var req struct {
		Name     string `json:"name"`
		Number   string `json:"number"`
		IsActive bool   `json:"isActive"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
		})
		return
	}

	var vehicle models.Vehicle
	if err := h.db.First(&vehicle, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Vehicle not found",
		})
		return
	}

	// Update fields
	if req.Name != "" {
		nameParts := strings.SplitN(req.Name, " ", 2)
		vehicle.Make = nameParts[0]
		if len(nameParts) > 1 {
			vehicle.Model = nameParts[1]
		}
	}
	if req.Number != "" {
		vehicle.PlateNumber = req.Number
	}
	if req.IsActive {
		vehicle.Status = "active"
	} else {
		vehicle.Status = "inactive"
	}

	if err := h.db.Save(&vehicle).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Failed to update vehicle",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Vehicle updated successfully",
	})
}

// GET /api/documents - Fetch all documents (requires JWT)
func (h *APIHandler) GetAllDocuments(c *gin.Context) {
	// Get user ID from JWT context
	userID, exists := c.Get("uid")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "error",
			"message": "User not authenticated",
		})
		return
	}

	// Find driver by user ID
	var driver models.Driver
	if err := h.db.Where("user_id = ?", userID).First(&driver).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Driver not found",
		})
		return
	}

	// Mock documents for now - you can implement actual document model
	documents := []gin.H{
		{
			"id":         1,
			"type":       "Driving License",
			"status":     "Valid",
			"expiryDate": "2026-03-01",
		},
		{
			"id":         2,
			"type":       "RC Book",
			"status":     "Valid",
			"expiryDate": "2028-05-01",
		},
		{
			"id":         3,
			"type":       "Insurance",
			"status":     "Valid",
			"expiryDate": "2025-12-31",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "success",
		"documents": documents,
	})
}

// PUT /api/documents/:id - Update document (requires JWT)
func (h *APIHandler) UpdateDocument(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status     string `json:"status" binding:"required"`
		ExpiryDate string `json:"expiryDate" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
		})
		return
	}

	// Mock update - implement actual document update logic
	h.log.Info("Document update requested", zap.String("documentId", id))
	
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Document updated successfully",
	})
}

// POST /api/help/chat/start - Start live chat (requires JWT)
func (h *APIHandler) StartLiveChat(c *gin.Context) {
	var req struct {
		DriverID string `json:"driverId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
		})
		return
	}

	// Generate chat session ID
	chatSessionID := "CHAT-" + time.Now().Format("20060102") + "-" + strconv.FormatInt(time.Now().Unix(), 10)

	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"chatSessionId": chatSessionID,
		"message":       "Chat started with support",
	})
}

// POST /api/ride/accept - Accept a ride (requires JWT)
func (h *APIHandler) AcceptRide(c *gin.Context) {
	var req struct {
		DriverID   string `json:"driverId" binding:"required"`
		RideID     string `json:"rideId" binding:"required"`
		Status     string `json:"status" binding:"required"`
		AcceptedAt string `json:"acceptedAt" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
		})
		return
	}

	// Mock ride acceptance - implement actual ride logic
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ride accepted successfully",
		"rideDetails": gin.H{
			"rideId":           req.RideID,
			"passengerName":    "Mock Passenger",
			"pickup":           "Mock Pickup Location",
			"drop":             "Mock Drop Location",
			"fare":             450.0,
			"paymentMethod":    "Cash",
			"passengerContact": "+91 9876543210",
			"passengerLocation": gin.H{
				"latitude":  11.2587,
				"longitude": 75.7803,
			},
			"status": "Ongoing",
		},
	})
}

// POST /api/ride/cancel - Cancel a ride (requires JWT)
func (h *APIHandler) CancelRide(c *gin.Context) {
	var req struct {
		RideID      string `json:"rideId" binding:"required"`
		DriverID    string `json:"driverId" binding:"required"`
		Reason      string `json:"reason" binding:"required"`
		CancelledAt string `json:"cancelledAt" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
		})
		return
	}

	// Mock ride cancellation - implement actual ride logic
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Ride cancelled successfully",
	})
}

// POST /api/chat/send - Send chat message (requires JWT)
func (h *APIHandler) SendChatMessage(c *gin.Context) {
	var req struct {
		RideID     string `json:"rideId" binding:"required"`
		SenderID   string `json:"senderId" binding:"required"`
		ReceiverID string `json:"receiverId" binding:"required"`
		Message    string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid request data",
		})
		return
	}

	// Mock chat message - implement actual chat logic
	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
