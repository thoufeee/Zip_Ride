package routes

import (
	"zipride/internal/user_Auth/handlers"

	"github.com/gin-gonic/gin"
)

// public routes

func PublicRoutes(c *gin.Engine) {
	api := c.Group("/")

	// otp based signup
	api.POST("signup-otp", handlers.SendOtpHandler)
	api.POST("verify-otp", handlers.VerifyOtpHandler)
	api.POST("otp-signup", handlers.RegisterUser)

	// otp based signin
	api.POST("signin-otp", handlers.OtpSignin)
	api.POST("signin-verifyotp", handlers.VerifyOTP)

	// google signup && sigin
	api.POST("googlesignup", handlers.GoogleSignup)
	api.POST("googlesigin", handlers.GoogleSigin)

	// token based authentication
	api.POST("signup", handlers.SignUp)
	api.POST("signin", handlers.SignIn)

	// forgott password
	api.POST("/forget-password", handlers.ForgetPassword)
	api.POST("/veryfy-otp", handlers.VerifyForgotOTP)
	api.POST("/reset-password", handlers.ResetPassword)

}
