package controller

import (
	"aino_document/models"
	"aino_document/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/lib/pq"
)

func AddRole(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

	// Periksa apakah tokenString tidak kosong
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Langkah 1: Mendekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userUUID := c.Get("user_uuid").(string)
	_, errK := service.GetUserInfoFromToken(tokenOnly)
	if errK != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid token atau token tidak ditemukan!")
	}

	var addRole models.Role

	if err := c.Bind(&addRole); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	errVal := c.Validate(&addRole)

	if errVal == nil {
		addroleErr := service.AddRole(addRole, userUUID)
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

}

func GetAllRole(c echo.Context) error {
	role, err := service.GetAllRole()

	if err != nil {
		log.Print(err)
		response := models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal server. Mohon coba beberapa saat lagi",
			Status:  false,
		}
		return c.JSON(http.StatusInternalServerError, response)
	}
	return c.JSON(http.StatusOK, role)
}

func ShowRoleById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var getRole models.Roles

	getRole, err := service.ShowRoleById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Role tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getRole)
}

func GetRoleById(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))

	var getRole models.Role

	getRole, err := service.GetRoleById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			response := models.Response{
				Code:    404,
				Message: "Role tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getRole)
}

func UpdateRole(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

	// Periksa apakah tokenString tidak kosong
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	// Periksa apakah tokenString mengandung "Bearer "
	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	// Hapus "Bearer " dari tokenString
	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Langkah 1: Mendekripsi token JWE
	decrypted, err := DecryptJWE(tokenOnly, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}
	userUUID := c.Get("user_uuid").(string)
	_, errK := service.GetUserInfoFromToken(tokenOnly)
	if errK != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid token atau token tidak ditemukan!")
	}

	id, _ := strconv.Atoi(c.Param("id"))

	perviousContent, errGet := service.ShowRoleById(id)
	if errGet != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Gagal mengambil data role saat ini. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	var editRole models.Role
	if err := c.Bind(&editRole); err != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data invalid!",
			Status:  false,
		})
	}

	errValidate := c.Validate(&editRole)
	if errValidate != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    422,
			Message: "Data tidak boleh kosong!",
			Status:  false,
		})
	}
	if err == nil {
		_, errService := service.UpdateRole(editRole, id, userUUID)
		if errService != nil {
			if dbErr, ok := errService.(*pq.Error); ok {
				// Check for duplicate key violation (unique constraint violation)
				if dbErr.Code.Name() == "unique_violation" {
					log.Println(dbErr)
					return c.JSON(http.StatusBadRequest, &models.Response{
						Code:    400,
						Message: "Data sudah ada",
						Status:  false,
					})
				}
			}

			log.Println("Kesalahan selama pembaruan:", errService)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi",
				Status:  false,
			})
		}

		log.Println(perviousContent)
		return c.JSON(http.StatusOK, &models.Response{
			Code:    200,
			Message: "Role telah diperbarui",
			Status:  true,
		})
	} else {
		log.Println("Kesalahan sebelum pembaruan:", err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi",
			Status:  false,
		})
	}
}
