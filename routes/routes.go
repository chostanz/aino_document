package routes

import (
	"aino_document/controller"
	"aino_document/middleware"

	"github.com/labstack/echo/v4"
)

func Route() *echo.Echo {
	r := echo.New()

	//test
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

	//update
	superAdminGroup.PUT("/division/update/:id", controller.UpdateDivision)
	superAdminGroup.PUT("/role/update/:id", controller.UpdateRole)
	superAdminGroup.PUT("/application/update/:id", controller.UpdateApp)
	superAdminGroup.PUT("/application/role/update/:id", controller.UpdateAppRole)
	superAdminGroup.PUT("/user/update/:id", controller.UpdateUser)

	//delete
	superAdminGroup.PUT("/role/delete/:id", controller.DeleteRole)
	superAdminGroup.PUT("/division/delete/:id", controller.DeleteDivision)
	superAdminGroup.PUT("/application/delete/:id", controller.DeleteApp)
	superAdminGroup.PUT("/application/role/delete/:id", controller.DeleteAppRole)

	authGroup := r.Group("/auth")
	authGroup.Use(middleware.AuthMiddleware)
	authGroup.PUT("/change/password", controller.ChangePassword)

	//get all
	r.GET("/division/all", controller.GetAllDivision)
	r.GET("/role/all", controller.GetAllRole)
	r.GET("/application/all", controller.GetAllApp)

	//get spec
	r.GET("/application/:id", controller.ShowApplicationById)
	r.GET("/get/application/:id", controller.GetAppById)
	r.GET("/division/:id", controller.ShowDivisionById)
	r.GET("/get/division/:id", controller.GetDivisionById)
	r.GET("/role/:id", controller.ShowRoleById)
	r.GET("/get/role/:id", controller.GetRoleById)

	//app role opsional
	r.GET("/application/role/all", controller.GetAllAppRole)
	r.GET("/application/role/:id", controller.GetAppRole)
	r.POST("/login", controller.Login)
	r.POST("logout", controller.Logout)

	//user
	superAdminGroup.PUT("/user/update/:user_application_role_uuid", controller.UpdateUserAppRole)
	r.GET("/user/:id", controller.ShowAppRoleById)
	superAdminGroup.GET("/user/all", controller.GetUserAppRole)
	superAdminGroup.PUT("/user/delete/:id", controller.DeleteUserAppRole)
	return r

}
