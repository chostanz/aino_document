package controller

import (
	"aino_document/models"
	"aino_document/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddAppRole(c echo.Context) error {
	var addAppRole models.ApplicationRole

	c.Bind(&addAppRole)

	if err := service.AddApplicationRole(addAppRole); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server",
			Status:  false,
		})
	}
	return c.JSON(http.StatusCreated, &models.Response{
		Code:    201,
		Message: "Berhasil menambahkan application role!",
		Status:  true,
	})
}
