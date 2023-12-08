package routes

import (
	"aino_document/controller"
	"aino_document/middleware"

	"github.com/labstack/echo/v4"
)

func Route() *echo.Echo {
	r := echo.New()

	r.POST("/user/add", controller.RegisterUser)
	r.POST("/application/add", controller.AddApplication)
	r.POST("/division/add", controller.AddDivision)
	r.POST("/application/role/add", controller.AddAppRole)
	r.POST("/role/add", controller.AddRole)

	superAdminGroup := r.Group("/superadmin")
	superAdminGroup.Use(middleware.AdminMiddleware)
	superAdminGroup.GET("/user/application/role", controller.GetUserAppRole)
	superAdminGroup.POST("/application/add", controller.AddApplication)
	superAdminGroup.POST("/application/role/add", controller.AddAppRole)
	superAdminGroup.POST("/division/add", controller.AddDivision)
	superAdminGroup.POST("/role/add", controller.AddRole)
	superAdminGroup.POST("/user/add", controller.RegisterUser)
	superAdminGroup.PUT("/division/update/:id", controller.UpdateDivision)
	superAdminGroup.PUT("/role/update/:id", controller.UpdateRole)
	superAdminGroup.PUT("/application/update/:id", controller.UpdateApp)
	superAdminGroup.PUT("/application/role/update/:id", controller.UpdateAppRole)
	superAdminGroup.PUT("/user/update/:id", controller.UpdateUser)

	superAdminGroup.PUT("/role/delete/:id", controller.DeleteRole)
	superAdminGroup.PUT("/division/delete/:id", controller.DeleteDivision)
	superAdminGroup.PUT("/application/delete/:id", controller.DeleteApp)
	superAdminGroup.PUT("/application/role/delete/:id", controller.DeleteAppRole)

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.PUT("/change/password", controller.ChangePassword)

	r.GET("/user/application/role", controller.GetUserAppRole)
	r.GET("/division/all", controller.GetAllDivision)
	r.GET("/role/all", controller.GetAllRole)
	r.GET("/application/all", controller.GetAllApp)
	r.GET("/application/:id", controller.ShowApplicationById)
	r.GET("/get/application/:id", controller.GetAppById)
	r.GET("/division/:id", controller.ShowDivisionById)
	r.GET("/get/division/:id", controller.GetDivisionById)
	r.GET("/role/:id", controller.ShowRoleById)
	r.GET("/get/role/:id", controller.GetRoleById)
	r.GET("/application/role/all", controller.GetAllAppRole)
	r.GET("/application/role/:id", controller.GetAppRole)
	r.POST("/login", controller.Login)
	r.POST("logout", controller.Logout)
	return r

}
