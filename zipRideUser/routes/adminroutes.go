package routes

import (
	"zipride/internal/constants"

	"zipride/internal/domain/user_Admin/allpermission"
	staffmanagment "zipride/internal/domain/user_Admin/staffManagment"
	"zipride/internal/domain/user_Admin/staffManagment/controllers"
	services "zipride/internal/domain/user_Admin/userManagment/controllers"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

// admin routes

func SuperAdminRoutes(c *gin.Engine) {

	admin := c.Group("/admin")

	admin.Use(middleware.JwtValidation())
	admin.Use(middleware.RoleCheck(constants.RoleSuperAdmin, constants.RoleManager, constants.RoleStaff))

	// verify route
	admin.GET("/verify", allpermission.Verify)

	// admin || staff || manager profile
	admin.GET("/profile", staffmanagment.StaffProfile)

	//user management
	admin.GET("/allusers", middleware.RequirePermission(constants.PermissionViewUsers), services.GetAllUsers)
	admin.POST("/createuser", middleware.RequirePermission(constants.PermissionAddUser), services.AddUser)
	admin.PUT("/user/:id", middleware.RequirePermission(constants.PermissionEditUser), services.UpdateUser)
	admin.DELETE("/user/:id", middleware.RequirePermission(constants.PermissionDeleteUser), services.DeleteUser)
	admin.PUT("/userblock/:id", middleware.RequirePermission(constants.PermissionBlockUser), services.BlockUser)
	admin.PUT("/userunblock/:id", middleware.RequirePermission(constants.PermissionUnBlockUser), services.UnBlockUser)

	//Staff management
	admin.GET("/allstaffs", middleware.RequirePermission(constants.PermissionViewStaffs), controllers.GETAllStaff)
	admin.POST("/createstaff", middleware.RequirePermission(constants.PermissionAddStaff), controllers.CreateStaff)
	admin.PUT("/staffupdate/:id", middleware.RequirePermission(constants.PermissionEditStaff), controllers.UpdateStaff)
	admin.DELETE("/staffdelete/:id", middleware.RequirePermission(constants.PermissionDeleteStaff), controllers.UpdateStaff)

	//route for all permissions && roles
	admin.GET("/allpermissions", middleware.RequirePermission(constants.ViewAllPermissions), allpermission.Permissions)
	admin.GET("/allroles", middleware.RequirePermission(constants.PermissionViewAllRoles), allpermission.AllRoles)

}
