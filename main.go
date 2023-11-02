package main

import (
	"aino_document/routes"
	"aino_document/utils"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := routes.Route()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())
	customValidator := &utils.CustomValidator{Validator: validator.New()}
	e.Validator = customValidator
	e.Logger.Fatal(e.Start(":8080"))

}
