package adminhandlers

import (
	"net/http"
	"encoding/csv"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"zipRideDriver/internal/models"
)

type EarningsHandler struct{ db *gorm.DB; log *zap.Logger }

func NewEarningsHandler(db *gorm.DB, log *zap.Logger) *EarningsHandler { return &EarningsHandler{db: db, log: log} }

// GET /admin/panel/earnings
func (h *EarningsHandler) Index(c *gin.Context) {
	var total float64
	_ = h.db.Model(&models.Earning{}).Select("COALESCE(SUM(amount),0)").Scan(&total).Error
	var items []models.Earning
	_ = h.db.Order("created_at desc").Limit(200).Find(&items).Error
	c.HTML(http.StatusOK, "admin/earnings/index.html", gin.H{"Total": total, "Earnings": items})
}

// GET /admin/panel/withdrawals
func (h *EarningsHandler) Withdrawals(c *gin.Context) {
	var pending []models.Withdrawal
	_ = h.db.Where("status = ?", "pending").Order("created_at asc").Find(&pending).Error
	c.HTML(http.StatusOK, "admin/earnings/withdrawals.html", gin.H{"Pending": pending})
}

// POST /admin/panel/withdrawals/:id/approve
func (h *EarningsHandler) ApproveWithdrawal(c *gin.Context) {
	id := c.Param("id")
	_ = h.db.Model(&models.Withdrawal{}).Where("id = ?", id).Update("status", "approved").Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/withdrawals")
}

// POST /admin/panel/withdrawals/:id/reject
func (h *EarningsHandler) RejectWithdrawal(c *gin.Context) {
	id := c.Param("id")
	_ = h.db.Model(&models.Withdrawal{}).Where("id = ?", id).Update("status", "rejected").Error
	c.Redirect(http.StatusSeeOther, "/admin/panel/withdrawals")
}

// GET /admin/panel/earnings/export
func (h *EarningsHandler) ExportCSV(c *gin.Context) {
    var items []models.Earning
    _ = h.db.Order("created_at desc").Find(&items).Error
    c.Header("Content-Type", "text/csv")
    c.Header("Content-Disposition", "attachment; filename=earnings.csv")
    w := csv.NewWriter(c.Writer)
    _ = w.Write([]string{"ID", "DriverID", "RideID", "Amount", "CreatedAt"})
    for _, e := range items {
        _ = w.Write([]string{
            strconv.FormatUint(uint64(e.ID), 10),
            strconv.FormatUint(uint64(e.DriverID), 10),
            strconv.FormatUint(uint64(e.RideID), 10),
            strconv.FormatFloat(e.Amount, 'f', 2, 64),
            e.CreatedAt.Format("2006-01-02 15:04:05"),
        })
    }
    w.Flush()
}
