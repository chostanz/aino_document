package controller

import (
	"aino_document/models"
	"aino_document/service"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func AddRole(c echo.Context) error {
	var addRole models.Role

	if err := c.Bind(&addRole); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	errVal := c.Validate(&addRole)
	// if errVal != nil {
	// 	return c.JSON(http.StatusUnprocessableEntity, &models.Response{
	// 		Code:    422,
	// 		Message: "Data tidak boleh kosong!",
	// 		Status:  false,
	// 	})
	// }

	if errVal == nil {
		addroleErr := service.AddRole(addRole)
		if addroleErr != nil {
			if dbErr, ok := addroleErr.(*pq.Error); ok {
				if dbErr.Code.Name() == "unique_violation" {
					return c.JSON(http.StatusBadRequest, &models.Response{
						Code:    400,
						Message: "Gagal menambahkan role. Role sudah ada!",
						Status:  false,
					})
				}
			}
			return c.JSON(http.StatusInternalServerError, "Error")
		}
		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil menambahkan role!",
			Status:  true,
		})
	} else {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}

	// if err := service.AddRole(addRole); err != nil {
	// 	log.Print(err)

	// 	if dbErr, ok := err.(*pq.Error); ok {
	// 		if dbErr.Code.Name() == "unique_violation" {
	// 			return c.JSON(http.StatusBadRequest, &models.Response{
	// 				Code:    400,
	// 				Message: "Gagal menambahkan role. Role sudah ada!",
	// 				Status:  false,
	// 			})
	// 		}
	// 	}
	// 	return c.JSON(http.StatusInternalServerError, &models.Response{
	// 		Code:    500,
	// 		Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi",
	// 		Status:  false,
	// 	})
	// }

	// return c.JSON(http.StatusCreated, &models.Response{
	// 	Code:    201,
	// 	Message: "Berhasil menambahkan role!",
	// 	Status:  true,
	// })
}
