package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// Me returns driver profile and onboarding status
func Me(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)
	var d models.Driver
	if err := database.DB.First(&d, uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "driver not found"})
		return
	}
	var vehicle models.Vehicle
	_ = database.DB.Where("driver_id = ?", d.ID).First(&vehicle).Error
	var docs models.DriverDocuments
	_ = database.DB.Where("driver_id = ?", d.ID).First(&docs).Error
	c.JSON(http.StatusOK, gin.H{"driver": d, "vehicle": vehicle, "docs": docs})
}

// UpdateProfile updates basic driver profile fields
func UpdateProfile(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)
	var req struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Avatar    string `json:"avatar"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updates := map[string]any{}
	if req.FirstName != "" { updates["first_name"] = req.FirstName }
	if req.LastName != "" { updates["last_name"] = req.LastName }
	if req.Avatar != "" { updates["avatar"] = req.Avatar }
	if len(updates) > 0 {
		if err := database.DB.Model(&models.Driver{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "profile updated"})
}

// UpsertVehicle creates or updates the driver's vehicle
func UpsertVehicle(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)
	var req models.Vehicle
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.DriverID = uid
	var existing models.Vehicle
	if err := database.DB.Where("driver_id = ?", uid).First(&existing).Error; err == nil {
		// update
		req.ID = existing.ID
		if err := database.DB.Model(&existing).Updates(req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update vehicle"})
			return
		}
	} else {
		if err := database.DB.Create(&req).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vehicle"})
			return
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "vehicle saved"})
}

// UploadDocs saves or updates driver documents and marks status to in_review
func UploadDocs(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)
	var req struct {
		LicenseURL   string `json:"license_url" binding:"required"`
		RCURL        string `json:"rc_url" binding:"required"`
		InsuranceURL string `json:"insurance_url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var docs models.DriverDocuments
	if err := database.DB.Where("driver_id = ?", uid).First(&docs).Error; err == nil {
		docs.LicenseURL = req.LicenseURL
		docs.RCURL = req.RCURL
		docs.InsuranceURL = req.InsuranceURL
		docs.Status = "in_review"
		if err := database.DB.Save(&docs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update docs"})
			return
		}
	} else {
		docs = models.DriverDocuments{DriverID: uid, LicenseURL: req.LicenseURL, RCURL: req.RCURL, InsuranceURL: req.InsuranceURL, Status: "in_review"}
		if err := database.DB.Create(&docs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create docs"})
			return
		}
	}
	// Move driver to in_review when docs submitted
	_ = database.DB.Model(&models.Driver{}).Where("id = ?", uid).Update("status", "in_review").Error
	c.JSON(http.StatusOK, gin.H{"message": "documents submitted"})
}