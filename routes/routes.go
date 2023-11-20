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

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AdminMiddleware)
	// adminGroup.GET("/user/application/role", controller.GetUserAppRole)
	adminGroup.POST("/application/add", controller.AddApplication)
	adminGroup.POST("/application/role/add", controller.AddAppRole)
	adminGroup.POST("/division/add", controller.AddDivision)
	adminGroup.POST("/role/add", controller.AddRole)
	adminGroup.POST("/user/add", controller.RegisterUser)

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
	r.GET("/application/role/:id", controller.GetAppRole)
	r.POST("/login", controller.Login)
	r.POST("logout", controller.Logout)
	return r

}
