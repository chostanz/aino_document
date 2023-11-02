package controller

import (
	"aino_document/models"
	"aino_document/service"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func AddApplication(c echo.Context) error {
	var addApplication models.Application

	if err := c.Bind(&addApplication); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	err := c.Validate(&addApplication)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	if err := service.AddApplication(addApplication); err != nil {
		log.Print(err)
		if dbErr, ok := err.(*pq.Error); ok {
			if dbErr.Code.Name() == "unique_violation" {
				return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Gagal menambahkan application. Application sudah ada!",
					Status:  false,
				})
			}
		}
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi",
			Status:  false,
		})
	}
	return c.JSON(http.StatusCreated, &models.Response{
		Code:    201,
		Message: "Berhasil menambahkan application!",
		Status:  true,
	})
}
