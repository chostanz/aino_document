package routes

import (
	"aino_document/controller"

	"github.com/labstack/echo/v4"
)

func Route() *echo.Echo {
	r := echo.New()

	r.POST("/user/add", controller.RegisterUser)
	r.POST("/application/add", controller.AddApplication)
	r.POST("/division/add", controller.AddDivision)
	r.POST("/application/role/add", controller.AddAppRole)
	r.POST("/role/add", controller.AddRole)

	r.GET("user/application/role", controller.GetUserAppRole)

	r.POST("/login", controller.Login)
	r.POST("logout", controller.Logout)
	return r

}
