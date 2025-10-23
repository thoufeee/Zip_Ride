package routes

import (
	"zipride/internal/constants"
	"zipride/internal/domain/booking_module/handlers"
	ratingdriver "zipride/internal/domain/booking_module/handlers/Rating_Driver"
	"zipride/internal/domain/chat"
	"zipride/internal/domain/user/services"
	"zipride/internal/middleware"
	authHandlers "zipride/internal/user_Auth/handlers"

	"github.com/gin-gonic/gin"
)

// user routes
func UserRoutes(c *gin.Engine) {
	user := c.Group("/user")

	//user Forget password

	user.Use(middleware.JwtValidation())
	user.Use(middleware.RoleCheck(constants.RoleUser))

	//user profile set up
	user.GET("/profile", services.GetUserProfile)
	user.PUT("/update", services.UpdateUserProfile)
	user.DELETE("/delete", services.DeleteUserProfile)

	//booking module
	user.POST("/estimate", handlers.EstimateBooking)
	user.POST("/now", handlers.CreateBookingNow)
	user.POST("/later", handlers.CreateBookingLater)
	user.POST("/cancel", handlers.CancelBookingHandler)
	user.GET("/history", handlers.GetBookingHistoryHandler)

	//Chat
	user.GET("/chat/:booking_id", chat.ConnectUser)
	//testing
	// user.GET("/driver/bookings",kafka.GetKafkaBookings)

	//feedback
	user.POST("/rate-driver/:booking_id", ratingdriver.RateDriver)

	// logout
	user.POST("/logout", authHandlers.UserLogout)

}
