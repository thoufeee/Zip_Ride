package routes

import (
	"zipride/internal/admin/services"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

// admin routes

func AdminRoutes(c *gin.Engine) {

	admin := c.Group("/admin")

	admin.Use(middleware.Auth())

	admin.GET("/allusers", services.GetAllUsers)
}
