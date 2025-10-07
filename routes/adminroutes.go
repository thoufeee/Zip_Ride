package routes

import (
	"zipride/internal/constants"
	"zipride/internal/domain/user_Admin/controllers"
	"zipride/internal/domain/user_Admin/services"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

// admin routes

func AdminRoutes(c *gin.Engine) {

	admin := c.Group("/admin")

	admin.Use(middleware.Auth(constants.RoleSuperAdmin))

	admin.GET("/allusers", services.GetAllUsers)

	//Staff management
	admin.GET("/staffs", controllers.GETAllStaff)
	admin.POST("/createstaff", controllers.CreateStaff)
	admin.PUT("/staff/:id", controllers.UpdateStaff)
	admin.DELETE("/staff/:id", controllers.UpdateStaff)

	//user management
	admin.GET("/users",services.GetAllUsers)
	admin.POST("/createuser",services.AddUser)
	admin.PUT("/user/:id",services.UpdateUser)
	admin.DELETE("/user/:id",services.DeleteUser)
}
