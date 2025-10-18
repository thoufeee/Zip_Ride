package handlers

import (
	"net/http"
	"strings"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type HelpHandler struct {
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewHelpHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *HelpHandler {
	return &HelpHandler{db: db, rdb: rdb, log: log}
}

func (h *HelpHandler) ListFAQs(c *gin.Context) {
	var faqs []models.FAQ
	if err := h.db.Where("is_active = ?", true).Order("id asc").Find(&faqs).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to load faqs")
		return
	}
	utils.Ok(c, "faqs", faqs)
}

func (h *HelpHandler) CreateTicket(c *gin.Context) {
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
	if err := h.db.Create(&t).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to create ticket")
		return
	}
	utils.Ok(c, "ticket created", t)
}

func (h *HelpHandler) CreateReport(c *gin.Context) {
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
	if err := h.db.Create(&t).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to create report")
		return
	}
	utils.Ok(c, "report submitted", t)
}

func (h *HelpHandler) StartChat(c *gin.Context) {
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
	if err := h.db.Create(&s).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to start chat")
		return
	}
	m := models.ChatMessage{SessionID: s.ID, Sender: "driver", Message: strings.TrimSpace(req.Message), CreatedAt: time.Now().UTC()}
	if err := h.db.Create(&m).Error; err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to save message")
		return
	}
	utils.Ok(c, "chat started", gin.H{"session": s, "first_message": m})
}
