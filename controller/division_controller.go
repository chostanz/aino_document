package controller

import (
	"aino_document/models"
	"aino_document/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

func AddDivision(c echo.Context) error {
	var addDivision models.Division
	if err := c.Bind(&addDivision); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, "Error")
	}

	err := c.Validate(&addDivision)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	if err := service.AddDivision(addDivision); err != nil {
		log.Print(err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server",
			Status:  false,
		})
	}
	return c.JSON(http.StatusCreated, &models.Response{
		Code:    201,
		Message: "Berhasil menambahkan divisi!",
		Status:  true,
	})
}
