package routes

import (
	"zipride/internal/userauth/handlers"

	"github.com/gin-gonic/gin"
)

// public routes

func PublicRoutes(c *gin.Engine) {
	api := c.Group("/")

	// api.POST("signup-otp", handlers.SendOtpHandler)
	// api.POST("verify-otp", handlers.VerifyOtpHandler)
	// api.POST("signup", handlers.RegisterUser)

	api.POST("googlesignup", handlers.GoogleSignup)
	api.POST("googlesigin", handlers.GoogleSigin)

	api.POST("signup", handlers.SignUp)
	api.POST("signin", handlers.SignIn)

}
