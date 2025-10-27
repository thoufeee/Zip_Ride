package routes

import (
	"zipride/internal/constants"

	"zipride/internal/domain/user_Admin/allpermission"
	"zipride/internal/domain/user_Admin/configuration"
	"zipride/internal/domain/user_Admin/prize_pool/prizepoolmanagment"
	staffmanagment "zipride/internal/domain/user_Admin/staffManagment"
	"zipride/internal/domain/user_Admin/staffManagment/controllers"
	subscriptionuser "zipride/internal/domain/user_Admin/subscription_Plan/subscription_user"
	"zipride/internal/domain/user_Admin/subscription_Plan/subscriptionmanagment"
	services "zipride/internal/domain/user_Admin/userManagment/services"
	vehiclemanagement "zipride/internal/domain/user_Admin/vehicleManagement"
	"zipride/internal/middleware"

	"github.com/gin-gonic/gin"
)

// admin routes

func SuperAdminRoutes(c *gin.Engine) {

	admin := c.Group("/admin")

	admin.Use(middleware.JwtValidation())
	admin.Use(middleware.RoleCheck(constants.RoleSuperAdmin, constants.RoleAdmin))

	// verify route
	admin.GET("/verify", allpermission.Verify)

	admin.GET("/profile", staffmanagment.StaffProfile)

	//user management
	admin.GET("/allusers", middleware.RequirePermission(constants.PermissionViewUsers), services.GetAllUsers)
	admin.GET("/alluserslength", middleware.RequirePermission(constants.PermissionAccessAdminDash), services.AllusersLength)
	admin.GET("/latestuserslength", middleware.RequirePermission(constants.PermissionAccessAdminDash), services.LatestUsersLength)
	admin.POST("/createuser", middleware.RequirePermission(constants.PermissionAddUser), services.AddUser)
	admin.PUT("/user/:id", middleware.RequirePermission(constants.PermissionEditUser), services.UpdateUser)
	admin.DELETE("/user/:id", middleware.RequirePermission(constants.PermissionDeleteUser), services.DeleteUser)
	admin.PUT("/userblock/:id", middleware.RequirePermission(constants.PermissionBlockUser), services.BlockUser)
	admin.PUT("/userunblock/:id", middleware.RequirePermission(constants.PermissionUnBlockUser), services.UnBlockUser)

	//Staff management
	admin.GET("/allstaffs", middleware.RequirePermission(constants.PermissionViewStaffs), controllers.GETAllStaff)
	admin.POST("/createstaff", middleware.RequirePermission(constants.PermissionAddStaff), controllers.CreateStaff)
	admin.PUT("/staffupdate/:id", middleware.RequirePermission(constants.PermissionEditStaff), controllers.UpdateStaff)
	admin.DELETE("/staffdelete/:id", middleware.RequirePermission(constants.PermissionDeleteStaff), controllers.DeleteStaff)
	admin.PUT("/staffblock/:id", middleware.RequirePermission(constants.PermissionBlockStaff), controllers.BlockStaff)
	admin.PUT("/staffunblock/:id", middleware.RequirePermission(constants.PermissionUnBlockStaff), controllers.UnblockStaff)

	//route for all permissions
	admin.GET("/allpermissions", middleware.RequirePermission(constants.ViewAllPermissions), allpermission.Permissions)

	// Vehicle Fare Management (SuperAdmin / Manager)
	vehicleFare := admin.Group("/vehiclefare")
	{
		vehicleFare.POST("/", middleware.RequirePermission(constants.PermissionSystemSettings), vehiclemanagement.VehicleFareCreation)
		vehicleFare.GET("/", middleware.RequirePermission(constants.PermissionSystemSettings), vehiclemanagement.GetAllVehicleFares)
		vehicleFare.PUT("/:id", middleware.RequirePermission(constants.PermissionSystemSettings), vehiclemanagement.UpdateVehicleFare)
		vehicleFare.DELETE("/:id", middleware.RequirePermission(constants.PermissionSystemSettings), vehiclemanagement.DeleteVehicleFare)
	}

	// prize pool

	pricePool := admin.Group("/pricepool")

	{
		pricePool.GET("/", middleware.RequirePermission(constants.PermissionPrizePool), prizepoolmanagment.GetAllPrizePool)
		pricePool.POST("/", middleware.RequirePermission(constants.PermissionPrizePool), prizepoolmanagment.CreatePrizePool)
		pricePool.PUT("/:id", middleware.RequirePermission(constants.PermissionPrizePool), prizepoolmanagment.UpdatePrizePool)
		pricePool.PUT("/status/:id", middleware.RequirePermission(constants.PermissionPrizePool), prizepoolmanagment.UpdateStatus)
		pricePool.DELETE("/:id", middleware.RequirePermission(constants.PermissionPrizePool), prizepoolmanagment.DeletePrizePool)
	}

	// subscription plan

	subscription := admin.Group("/subscription")

	{
		subscription.GET("/", middleware.RequirePermission(constants.PermissionSubscription), subscriptionmanagment.GetSubScription)
		subscription.POST("/", middleware.RequirePermission(constants.PermissionSubscription), subscriptionmanagment.CreateSubscription)
		subscription.PUT("/:id", middleware.RequirePermission(constants.PermissionSubscription), subscriptionmanagment.UpdateSubScription)
		subscription.DELETE("/:id", middleware.RequirePermission(constants.PermissionSubscription), subscriptionmanagment.DeleteSubscription)

		// subscribed users details
		subscription.GET("/users", middleware.RequirePermission(constants.PermissionSubscription), subscriptionuser.SubScribedUser)
	}

	// configuration

	config := admin.Group("/config")

	{
		config.GET("/", middleware.RequirePermission(constants.PermissionSystemSettings), configuration.AllConfiguration)
		config.PUT("/", middleware.RequirePermission(constants.PermissionSystemSettings), configuration.UpdateConfig)
	}

}
