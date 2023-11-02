package controller

import (
	"aino_document/models"
	"aino_document/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUserAppRole(c echo.Context) error {
	userAppRole, err := service.GetUserApplicationRole()

	if err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		})
	}
	return c.JSON(http.StatusOK, userAppRole)
}
