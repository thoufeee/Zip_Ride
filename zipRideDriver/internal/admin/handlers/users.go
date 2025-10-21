package adminhandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zipRideDriver/internal/models"
)

type UsersHandler struct{ db *gorm.DB; log *zap.Logger }

func NewUsersHandler(db *gorm.DB, log *zap.Logger) *UsersHandler { return &UsersHandler{db: db, log: log} }

// GET /admin/panel/users
func (h *UsersHandler) Index(c *gin.Context) {
	var users []models.Rider
	_ = h.db.Order("created_at desc").Limit(200).Find(&users).Error
	c.HTML(http.StatusOK, "admin/users/index.html", gin.H{"Users": users})
}

// POST /admin/panel/users/:id/block
func (h *UsersHandler) Block(c *gin.Context) {
	id := c.Param("id")
	_ = h.db.Model(&models.Rider{}).Where("id = ?", id).Update("is_blocked", true).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/users")
}

// POST /admin/panel/users/:id/unblock
func (h *UsersHandler) Unblock(c *gin.Context) {
	id := c.Param("id")
	_ = h.db.Model(&models.Rider{}).Where("id = ?", id).Update("is_blocked", false).Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/users")
}
