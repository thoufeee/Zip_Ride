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
	"zipRideDriver/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type sendOTPReq struct{ Phone string `json:"phone" binding:"required"` }
type verifyOTPReq struct{ Phone string `json:"phone" binding:"required"`; OTP string `json:"otp" binding:"required"` }
type onboardingReq struct {
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email"`
	LicenseNumber  string `json:"license_number" binding:"required"`
	VehicleMake    string `json:"vehicle_make" binding:"required"`
	VehicleModel   string `json:"vehicle_model" binding:"required"`
	VehicleYear    int    `json:"vehicle_year" binding:"required"`
	PlateNumber    string `json:"plate_number" binding:"required"`
	InsuranceNumber string `json:"insurance_number" binding:"required"`
	RCNumber       string `json:"rc_number" binding:"required"`
}

func RegisterDriverRoutes(r *gin.Engine, cfg *config.Config, db *gorm.DB, rdb *redis.Client, log *zap.Logger) {
	g := r.Group("/driver")

	g.POST("/send-otp", func(c *gin.Context) {
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
 		// rate limit per phone and per IP
 		ctx := c.Request.Context()
 		ip := c.ClientIP()
 		pkey := "otp:rl:phone:" + phone
 		ikey := "otp:rl:ip:" + ip
 		if n, err := rdb.Incr(ctx, pkey).Result(); err == nil {
 			if n == 1 { _ = rdb.Expire(ctx, pkey, time.Minute).Err() }
 			if n > 3 { utils.Error(c, http.StatusTooManyRequests, "too many requests, try again later"); return }
 		}
 		if n, err := rdb.Incr(ctx, ikey).Result(); err == nil {
 			if n == 1 { _ = rdb.Expire(ctx, ikey, time.Minute).Err() }
 			if n > 10 { utils.Error(c, http.StatusTooManyRequests, "too many requests, try again later"); return }
 		}
		otp, err := services.GenerateNumericOTP(6)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to generate otp")
			return
		}
		key := services.OTPKey("driver", phone)
		if err := services.StoreOTP(context.Background(), rdb, key, otp, 5*time.Minute); err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to store otp")
			return
		}
		tw := services.NewTwilio(cfg, log)
		go func() { _ = tw.SendOTP(context.Background(), phone, otp) }()
		utils.Ok(c, "otp sent", nil)
	})

	g.POST("/verify-otp", func(c *gin.Context) {
		var req verifyOTPReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		phone := strings.TrimSpace(req.Phone)
		ok, err := services.VerifyOTP(context.Background(), rdb, services.OTPKey("driver", phone), strings.TrimSpace(req.OTP))
		if err != nil || !ok {
			utils.Error(c, http.StatusUnauthorized, "invalid or expired otp")
			return
		}
		var d models.Driver
		created := false
		if err := db.Where("phone = ?", phone).First(&d).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				d = models.Driver{Phone: phone, Status: "Pending", IsVerified: true}
				if err := db.Create(&d).Error; err != nil {
					utils.Error(c, http.StatusInternalServerError, "failed to create driver")
					return
				}
				created = true
			} else {
				utils.Error(c, http.StatusInternalServerError, "failed to load driver")
				return
			}
		} else {
			if !d.IsVerified {
				_ = db.Model(&models.Driver{}).Where("id = ?", d.ID).Update("is_verified", true).Error
			}
		}
		access, err := services.GenerateToken(d.ID, d.Email, "driver", []string{"driver"}, cfg.JWTSecret, cfg.AccessTokenExpiry)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to issue token")
			return
		}
		refresh, err := services.GenerateToken(d.ID, d.Email, "driver_refresh", nil, cfg.JWTSecret, cfg.RefreshTokenExpiry)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to issue token")
			return
		}
		utils.Ok(c, "verified", gin.H{"access_token": access, "refresh_token": refresh, "onboarding_required": created})
	})

	g.POST("/login", func(c *gin.Context) {
		var req sendOTPReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		c.Request.URL.Path = "/driver/send-otp"
		r.HandleContext(c)
	})

	g.POST("/refresh-token", func(c *gin.Context) {
		var body struct{ Refresh string `json:"refresh_token" binding:"required"` }
		if err := c.ShouldBindJSON(&body); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		claims, err := services.ParseToken(strings.TrimSpace(body.Refresh), cfg.JWTSecret)
		if err != nil || claims.Subject != "driver_refresh" {
			utils.Error(c, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		access, err := services.GenerateToken(claims.UserID, claims.Email, "driver", []string{"driver"}, cfg.JWTSecret, cfg.AccessTokenExpiry)
		if err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to issue access token")
			return
		}
		utils.Ok(c, "refreshed", gin.H{"access_token": access})
	})

	g.POST("/logout", func(c *gin.Context) {
		utils.Ok(c, "logged out", nil)
	})

	auth := g.Group("")
	auth.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	// Onboarding (protected)
	auth.POST("/onboarding", func(c *gin.Context) {
		uidAny, exists := c.Get("uid")
		if !exists {
			utils.Error(c, http.StatusUnauthorized, "unauthorized")
			return
		}
		uid, _ := uidAny.(uint)
		var req onboardingReq
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.Error(c, http.StatusBadRequest, err.Error())
			return
		}
		updates := map[string]interface{}{
			"name":           strings.TrimSpace(req.Name),
			"email":          strings.TrimSpace(req.Email),
			"license_number": strings.TrimSpace(req.LicenseNumber),
		}
		if err := db.Model(&models.Driver{}).Where("id = ?", uid).Updates(updates).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to update profile")
			return
		}
		v := models.Vehicle{
			DriverID:       uid,
			Make:           strings.TrimSpace(req.VehicleMake),
			Model:          strings.TrimSpace(req.VehicleModel),
			Year:           req.VehicleYear,
			PlateNumber:    strings.TrimSpace(req.PlateNumber),
			InsuranceNumber: strings.TrimSpace(req.InsuranceNumber),
			RCNumber:       strings.TrimSpace(req.RCNumber),
			Status:         "active",
		}
		if err := db.Create(&v).Error; err != nil {
			utils.Error(c, http.StatusInternalServerError, "failed to save vehicle")
			return
		}
		utils.Ok(c, "onboarding complete", nil)
	})
}
