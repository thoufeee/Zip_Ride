package adminhandlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zipRideDriver/internal/models"
)

type HelpHandler struct{ db *gorm.DB; log *zap.Logger }

func NewHelpHandler(db *gorm.DB, log *zap.Logger) *HelpHandler { return &HelpHandler{db: db, log: log} }

// GET /admin/panel/help
func (h *HelpHandler) Index(c *gin.Context) {
	var tickets []models.HelpTicket
	_ = h.db.Order("created_at desc").Limit(200).Find(&tickets).Error
	c.HTML(http.StatusOK, "admin/help/index.html", gin.H{"Tickets": tickets})
}

// GET /admin/panel/help/:id
func (h *HelpHandler) Show(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var t models.HelpTicket
	if err := h.db.First(&t, id).Error; err != nil { c.String(http.StatusNotFound, "Ticket not found"); return }
	var d models.Driver
	_ = h.db.First(&d, t.DriverID).Error
	// load or create chat session for this driver
	var sess models.ChatSession
	if err := h.db.Where("driver_id = ? AND status = ?", d.ID, "open").First(&sess).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			sess = models.ChatSession{DriverID: d.ID, Status: "open"}
			_ = h.db.Create(&sess).Error
		}
	}
	var msgs []models.ChatMessage
	_ = h.db.Where("session_id = ?", sess.ID).Order("id asc").Find(&msgs).Error
	c.HTML(http.StatusOK, "admin/help/show.html", gin.H{"Ticket": t, "Driver": d, "Session": sess, "Messages": msgs})
}

// POST /admin/panel/help/:id/reply
func (h *HelpHandler) Reply(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var t models.HelpTicket
	if err := h.db.First(&t, id).Error; err != nil { c.String(http.StatusNotFound, "Ticket not found"); return }
	message := strings.TrimSpace(c.PostForm("message"))
	if message == "" { c.Redirect(http.StatusSeeOther, "/admin/panel/help/"+strconv.Itoa(id)); return }
	// ensure session
	var sess models.ChatSession
	if err := h.db.Where("driver_id = ? AND status = ?", t.DriverID, "open").First(&sess).Error; err != nil {
		if err == gorm.ErrRecordNotFound { sess = models.ChatSession{DriverID: t.DriverID, Status: "open"}; _ = h.db.Create(&sess).Error }
	}
	_ = h.db.Create(&models.ChatMessage{SessionID: sess.ID, Sender: "admin", Message: message, CreatedAt: time.Now().UTC()}).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/help/"+strconv.Itoa(id))
}

// POST /admin/panel/help/:id/close
func (h *HelpHandler) Close(c *gin.Context) {
	id := c.Param("id")
	_ = h.db.Model(&models.HelpTicket{}).Where("id = ?", id).Update("status", "closed").Error
	// close any open chat session for this driver
	var t models.HelpTicket
	if err := h.db.First(&t, id).Error; err == nil {
		_ = h.db.Model(&models.ChatSession{}).Where("driver_id = ? AND status = ?", t.DriverID, "open").Update("status", "closed").Error
	}
	c.Redirect(http.StatusSeeOther, "/admin/panel/help")
}
