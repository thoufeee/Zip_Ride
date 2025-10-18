package handlers

import (
	"net/http"
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

func RegisterHelpRoutes(r *gin.Engine, cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) {
	pub := r.Group("/api/help")
	pub.GET("/faqs", func(c *gin.Context) {
		var faqs []models.FAQ
		if err := db.Where("is_active = ?", true).Order("id asc").Find(&faqs).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to load faqs")
			return
		}
		utils.Ok(c, "faqs", faqs)
	})

	auth := r.Group("/api/help")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	auth.POST("/ticket", func(c *gin.Context) {
		uidAny, _ := c.Get("uid")
		uid, _ := uidAny.(uint)
		var req struct {
			Subject     string `json:"subject" binding:"required"`
			Description string `json:"description" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		t := models.HelpTicket{DriverID: uid, Subject: strings.TrimSpace(req.Subject), Description: strings.TrimSpace(req.Description)}
		if err := db.Create(&t).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to create ticket")
			return
		}
		utils.Ok(c, "ticket created", t)
	})

	auth.POST("/report", func(c *gin.Context) {
		uidAny, _ := c.Get("uid")
		uid, _ := uidAny.(uint)
		var req struct {
			Issue string `json:"issue" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		t := models.HelpTicket{DriverID: uid, Subject: "Report", Description: strings.TrimSpace(req.Issue)}
		if err := db.Create(&t).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to create report")
			return
		}
		utils.Ok(c, "report submitted", t)
	})

	auth.POST("/chat/start", func(c *gin.Context) {
		uidAny, _ := c.Get("uid")
		uid, _ := uidAny.(uint)
		var req struct {
			Message string `json:"message" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		s := models.ChatSession{DriverID: uid, Status: "open"}
		if err := db.Create(&s).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to start chat")
			return
		}
		m := models.ChatMessage{SessionID: s.ID, Sender: "driver", Message: strings.TrimSpace(req.Message), CreatedAt: time.Now().UTC()}
		if err := db.Create(&m).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to save message")
			return
		}
		utils.Ok(c, "chat started", gin.H{"session": s, "first_message": m})
	})
}
