package handlers

import (
	"net/http"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
)

// CompleteProfileHandler allows a driver to complete their profile by providing basic info, vehicle, and docs
func CompleteProfileHandler(c *gin.Context) {
	val, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	uid, ok := val.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		FirstName    string          `json:"first_name"`
		LastName     string          `json:"last_name"`
		Avatar       string          `json:"avatar"`
		Vehicle      *models.Vehicle `json:"vehicle"`
		LicenseURL   string          `json:"license_url"`
		RCURL        string          `json:"rc_url"`
		InsuranceURL string          `json:"insurance_url"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update basic profile fields if provided
	updates := map[string]any{}
	if req.FirstName != "" {
		updates["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		updates["last_name"] = req.LastName
	}
	if req.Avatar != "" {
		updates["avatar"] = req.Avatar
	}
	if len(updates) > 0 {
		if err := database.DB.Model(&models.Driver{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update profile"})
			return
		}
	}

	// Upsert vehicle if provided
	if req.Vehicle != nil {
		req.Vehicle.DriverID = uid
		var existing models.Vehicle
		if err := database.DB.Where("driver_id = ?", uid).First(&existing).Error; err == nil {
			// keep immutable fields like ID
			req.Vehicle.ID = existing.ID
			if err := database.DB.Model(&existing).Updates(req.Vehicle).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update vehicle"})
				return
			}
		} else {
			if err := database.DB.Create(req.Vehicle).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vehicle"})
				return
			}
		}
	}

	// Upsert documents if provided
	if req.LicenseURL != "" || req.RCURL != "" || req.InsuranceURL != "" {
		var docs models.DriverDocuments
		if err := database.DB.Where("driver_id = ?", uid).First(&docs).Error; err == nil {
			if req.LicenseURL != "" {
				docs.LicenseURL = req.LicenseURL
			}
			if req.RCURL != "" {
				docs.RCURL = req.RCURL
			}
			if req.InsuranceURL != "" {
				docs.InsuranceURL = req.InsuranceURL
			}
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
	}

	c.JSON(http.StatusOK, gin.H{"message": "profile submitted; awaiting admin approval"})
}
