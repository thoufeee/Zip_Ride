package routes

import (
	"zipride/handlers"

	"github.com/gin-gonic/gin"
)

// account creating route

func User(c *gin.Engine) {
	api := c.Group("/")

	// api.POST("signup-otp", handlers.SendOtpHandler)
	// api.POST("verify-otp", handlers.VerifyOtpHandler)
	// api.POST("signup", handlers.RegisterUser)

	api.POST("signup", handlers.GoogleSignup)
	api.POST("sigin", handlers.GoogleSigin)

}
