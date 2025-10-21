package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"
)

type DriverRegistrationHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

func NewDriverRegistrationHandler(db *gorm.DB, log *zap.Logger) *DriverRegistrationHandler {
	return &DriverRegistrationHandler{db: db, log: log}
}

// RegisterDriver handles new driver registration
// POST /api/driver/register
func (h *DriverRegistrationHandler) RegisterDriver(c *gin.Context) {
	var req struct {
		FullName      string `json:"full_name" binding:"required"`
		Email         string `json:"email" binding:"required,email"`
		Phone         string `json:"phone" binding:"required"`
		Password      string `json:"password" binding:"required,min=6"`
		LicenseNumber string `json:"license_number" binding:"required"`
		VehicleType   string `json:"vehicle_type" binding:"required"`
		VehicleModel  string `json:"vehicle_model"`
		VehicleNumber string `json:"vehicle_number"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if driver already exists
	var existingDriver models.Driver
	if err := h.db.Where("email = ? OR phone = ?", req.Email, req.Phone).First(&existingDriver).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Driver with this email or phone already exists"})
		return
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		h.log.Error("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process registration"})
		return
	}

	// Create new driver with pending status
	driver := models.Driver{
		Name:          req.FullName,
		Email:         strings.ToLower(req.Email),
		Phone:         req.Phone,
		PasswordHash:  hashedPassword,
		LicenseNumber: req.LicenseNumber,
		VehicleType:   req.VehicleType,
		VehicleModel:  req.VehicleModel,
		VehicleNumber: req.VehicleNumber,
		Status:        "Pending", // Requires admin approval
		IsVerified:    false,
		IsOnline:      false,
		Rating:        5.0, // Start with default rating
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := h.db.Create(&driver).Error; err != nil {
		h.log.Error("Failed to create driver", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete registration"})
		return
	}

	h.log.Info("New driver registered", 
		zap.Uint("driver_id", driver.ID),
		zap.String("email", driver.Email),
		zap.String("status", driver.Status))

	c.JSON(http.StatusCreated, gin.H{
		"message": "Registration successful. Your account is pending approval from admin.",
		"driver_id": driver.ID,
		"status": driver.Status,
	})
}

// CheckRegistrationStatus checks the approval status of a driver
// GET /api/driver/registration-status/:email
func (h *DriverRegistrationHandler) CheckRegistrationStatus(c *gin.Context) {
	email := strings.ToLower(c.Param("email"))

	var driver models.Driver
	if err := h.db.Where("email = ?", email).First(&driver).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check status"})
		return
	}

	response := gin.H{
		"driver_id": driver.ID,
		"email": driver.Email,
		"status": driver.Status,
		"is_verified": driver.IsVerified,
		"registered_at": driver.CreatedAt,
	}

	// Add approval/rejection time if available
	if driver.VerifiedAt != nil {
		response["verified_at"] = driver.VerifiedAt
	}

	// Add appropriate message based on status
	switch driver.Status {
	case "Pending":
		response["message"] = "Your registration is pending admin approval"
	case "Approved", "Active":
		response["message"] = "Your account has been approved. You can now log in."
	case "Rejected":
		response["message"] = "Your registration was rejected. Please contact support for more information."
	case "Suspended":
		response["message"] = "Your account has been suspended. Please contact support."
	default:
		response["message"] = "Unknown status. Please contact support."
	}

	c.JSON(http.StatusOK, response)
}

// UploadDocument allows drivers to upload required documents
// POST /api/driver/upload-document
func (h *DriverRegistrationHandler) UploadDocument(c *gin.Context) {
	driverID := c.PostForm("driver_id")
	docType := c.PostForm("doc_type") // license, vehicle_registration, insurance, etc.

	// Get the file from form
	file, err := c.FormFile("document")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate driver exists
	var driver models.Driver
	if err := h.db.First(&driver, driverID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Driver not found"})
		return
	}

	// Save file (in production, you'd upload to cloud storage)
	filename := fmt.Sprintf("driver_%s_%s_%d%s", 
		driverID, 
		docType, 
		time.Now().Unix(),
		filepath.Ext(file.Filename))
	
	uploadPath := filepath.Join("uploads", "documents", filename)
	
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		h.log.Error("Failed to save document", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload document"})
		return
	}

	// Save document record
	doc := models.DriverDocument{
		DriverID:   driver.ID,
		DocType:    docType,
		DocURL:     uploadPath,
		Verified:   false,
		UploadedAt: time.Now(),
	}

	if err := h.db.Create(&doc).Error; err != nil {
		h.log.Error("Failed to save document record", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save document"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Document uploaded successfully",
		"document_id": doc.ID,
		"status": "pending_verification",
	})
}
