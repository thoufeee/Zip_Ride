package handlers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"zipRideDriver/internal/config"
	"zipRideDriver/internal/models"
	"zipRideDriver/internal/services"
	"zipRideDriver/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type sendOTPReq struct {
	Phone string `json:"phone" binding:"required"`
}

type verifyOTPReq struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required"`
}

type AuthHandler struct {
	cfg *config.Config
	db  *gorm.DB
	rdb *redis.Client
	log *zap.Logger
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) *AuthHandler {
	return &AuthHandler{cfg: cfg, db: db, rdb: rdb, log: log}
}

func (h *AuthHandler) SendOTP(c *gin.Context) {
	var req sendOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	phone := strings.TrimSpace(req.Phone)
	if phone == "" {
		utils.Error(c, http.StatusBadRequest, "phone required")
		return
	}
	ctx := c.Request.Context()
	ip := c.ClientIP()
	pkey := "otp:rl:phone:" + phone
	ikey := "otp:rl:ip:" + ip
	if n, err := h.rdb.Incr(ctx, pkey).Result(); err == nil {
		if n == 1 { _ = h.rdb.Expire(ctx, pkey, time.Minute).Err() }
		if n > 3 { utils.Error(c, http.StatusTooManyRequests, "too many requests, try again later"); return }
	}
	if n, err := h.rdb.Incr(ctx, ikey).Result(); err == nil {
		if n == 1 { _ = h.rdb.Expire(ctx, ikey, time.Minute).Err() }
		if n > 10 { utils.Error(c, http.StatusTooManyRequests, "too many requests, try again later"); return }
	}
	otp, err := services.GenerateNumericOTP(6)
	if err != nil { utils.Error(c, http.StatusInternalServerError, "failed to generate otp"); return }
	key := services.OTPKey("driver", phone)
	if err := services.StoreOTP(context.Background(), h.rdb, key, otp, 5*time.Minute); err != nil {
		utils.Error(c, http.StatusInternalServerError, "failed to store otp"); return }
	tw := services.NewTwilio(h.cfg, h.log)
	go func() { _ = tw.SendOTP(context.Background(), phone, otp) }()
	utils.Ok(c, "otp sent", nil)
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	var req verifyOTPReq
	if err := c.ShouldBindJSON(&req); err != nil { utils.Error(c, http.StatusBadRequest, err.Error()); return }
	phone := strings.TrimSpace(req.Phone)
	ok, err := services.VerifyOTP(context.Background(), h.rdb, services.OTPKey("driver", phone), strings.TrimSpace(req.OTP))
	if err != nil || !ok { utils.Error(c, http.StatusUnauthorized, "invalid or expired otp"); return }
	var d models.Driver
	created := false
	if err := h.db.Where("phone = ?", phone).First(&d).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			d = models.Driver{Phone: phone, Status: "Pending", IsVerified: true}
			if err := h.db.Create(&d).Error; err != nil { utils.Error(c, http.StatusInternalServerError, "failed to create driver"); return }
			created = true
		} else {
			utils.Error(c, http.StatusInternalServerError, "failed to load driver"); return
		}
	} else {
		if !d.IsVerified { _ = h.db.Model(&models.Driver{}).Where("id = ?", d.ID).Update("is_verified", true).Error }
	}
	access, err := services.GenerateToken(d.ID, d.Email, "driver", []string{"driver"}, h.cfg.JWTSecret, h.cfg.AccessTokenExpiry)
	if err != nil { utils.Error(c, http.StatusInternalServerError, "failed to issue token"); return }
	refresh, err := services.GenerateToken(d.ID, d.Email, "driver_refresh", nil, h.cfg.JWTSecret, h.cfg.RefreshTokenExpiry)
	if err != nil { utils.Error(c, http.StatusInternalServerError, "failed to issue token"); return }
	utils.Ok(c, "verified", gin.H{"access_token": access, "refresh_token": refresh, "onboarding_required": created})
}

func (h *AuthHandler) Login(c *gin.Context) { h.SendOTP(c) }

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var body struct{ Refresh string `json:"refresh_token" binding:"required"` }
	if err := c.ShouldBindJSON(&body); err != nil { utils.Error(c, http.StatusBadRequest, err.Error()); return }
	claims, err := services.ParseToken(strings.TrimSpace(body.Refresh), h.cfg.JWTSecret)
	if err != nil || claims.Subject != "driver_refresh" { utils.Error(c, http.StatusUnauthorized, "invalid refresh token"); return }
	access, err := services.GenerateToken(claims.UserID, claims.Email, "driver", []string{"driver"}, h.cfg.JWTSecret, h.cfg.AccessTokenExpiry)
	if err != nil { utils.Error(c, http.StatusInternalServerError, "failed to issue access token"); return }
	utils.Ok(c, "refreshed", gin.H{"access_token": access})
}

func (h *AuthHandler) Logout(c *gin.Context) { utils.Ok(c, "logged out", nil) }
