package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"zipride/database"
	"zipride/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UploadDocument handles file upload for driver documents
func UploadDocument(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)

	// Get document type from query
	docType := c.Query("type") // license, rc, insurance
	if docType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "document type required"})
		return
	}

	// Validate document type
	allowedTypes := []string{"license", "rc", "insurance"}
	if !contains(allowedTypes, docType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid document type"})
		return
	}

	// Get uploaded file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file required"})
		return
	}
	defer file.Close()

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowedExts := []string{".jpg", ".jpeg", ".png", ".pdf"}
	if !contains(allowedExts, ext) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid file type"})
		return
	}

	// Validate file size (5MB max)
	if header.Size > 5*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file too large"})
		return
	}

	// Create upload directory if not exists
	uploadDir := "uploads/documents"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("%s_%s_%d%s", docType, uuid.New().String(), time.Now().Unix(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filepath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}

	// Update driver documents
	url := fmt.Sprintf("/uploads/documents/%s", filename)
	var docs models.DriverDocuments
	if err := database.DB.Where("driver_id = ?", uid).First(&docs).Error; err != nil {
		// Create new document record
		docs = models.DriverDocuments{
			DriverID: uid,
			Status:   "pending",
		}
	}

	// Update specific document URL
	switch docType {
	case "license":
		docs.LicenseURL = url
	case "rc":
		docs.RCURL = url
	case "insurance":
		docs.InsuranceURL = url
	}

	// If all documents uploaded, set status to in_review
	if docs.LicenseURL != "" && docs.RCURL != "" && docs.InsuranceURL != "" {
		docs.Status = "in_review"
		// Update driver status
		database.DB.Model(&models.Driver{}).Where("id = ?", uid).Update("status", "in_review")
	}

	if err := database.DB.Save(&docs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save document info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "document uploaded successfully",
		"url":     url,
		"type":    docType,
	})
}

// GetDocumentStatus returns current document upload status
func GetDocumentStatus(c *gin.Context) {
	val, _ := c.Get("user_id")
	uid, _ := val.(uint)

	var docs models.DriverDocuments
	if err := database.DB.Where("driver_id = ?", uid).First(&docs).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "no documents found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"license_url":   docs.LicenseURL,
		"rc_url":        docs.RCURL,
		"insurance_url": docs.InsuranceURL,
		"status":        docs.Status,
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
