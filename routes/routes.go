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
	adminGroup.GET("/user/application/role", controller.GetUserAppRole)
	r.GET("/user/application/role", controller.GetUserAppRole)
	r.POST("/login", controller.Login)
	r.POST("logout", controller.Logout)
	return r

}
