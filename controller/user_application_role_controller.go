package controller

import (
	"aino_document/models"
	"aino_document/service"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/badoux/checkmail"
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

func ShowAppRoleById(c echo.Context) error {
	id := c.Param("id")

	var getAppRole models.Users

	getAppRole, err := service.GetSpecUseApplicationRole(id)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Print(err)
			response := models.Response{
				Code:    404,
				Message: "User tidak ditemukan!",
				Status:  false,
			}
			return c.JSON(http.StatusNotFound, response)
		} else {
			log.Print(err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi!",
				Status:  false,
			})
		}
	}

	return c.JSON(http.StatusOK, getAppRole)
}

func UpdateUserAppRole(c echo.Context) error {

	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak ditemukan!",
			"status":  false,
		})
	}

	if !strings.HasPrefix(tokenString, "Bearer ") {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"code":    401,
			"message": "Token tidak valid!",
			"status":  false,
		})
	}

	tokenOnly := strings.TrimPrefix(tokenString, "Bearer ")

	// Mendekripsi token JWE
	decrypted, errDec := DecryptJWE(tokenOnly, secretKey)
	if errDec != nil {
		fmt.Println("Gagal mendekripsi token:", errDec)
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
	//userUUID := c.Get("user_uuid").(string)
	_, errK := service.GetUserInfoFromToken(tokenOnly)
	if errK != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid token atau token tidak ditemukan!")
	}

	var updateUserAppRole models.UpdateUser

	if errBind := c.Bind(&updateUserAppRole); errBind != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	errEmail := checkmail.ValidateFormat(updateUserAppRole.Email)
	if errEmail != nil {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Format email tidak valid",
			Status:  false,
		})
	}

	err := c.Validate(&updateUserAppRole)

	if err == nil {
		// var existingUserID int
		// err := db.QueryRow("SELECT user_id FROM user_ms WHERE (user_name = $1 OR user_email = $2) AND deleted_at IS NULL", updateUserAppRole.Username, updateUserAppRole.Email).Scan(&existingUserID)

		// if err == nil {
		// 	log.Print(err)
		// 	return c.JSON(http.StatusBadRequest, &models.Response{
		// 		Code:    400,
		// 		Message: "Username atau email telah digunakan!",
		// 		Status:  false,
		// 	})
		// }
		userApplicationRoleUUID := c.Param("user_application_role_uuid")
		log.Println("user_application_role_uuid from URL:", userApplicationRoleUUID)

		existingUser, err := service.GetUserByUsernameAndEmail(userApplicationRoleUUID)
		if err != nil {
			log.Printf("Error getting existing user data: %v", err)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server.",
				Status:  false,
			})
		}

		// Data username atau email berubah, lakukan validasi unik
		if updateUserAppRole.Username == existingUser.UserName || updateUserAppRole.Email == existingUser.UserEmail {
			isUnique, err := service.IsUniqueUsernameOrEmail(userApplicationRoleUUID, updateUserAppRole.Username, updateUserAppRole.Email)
			if err != nil {
				log.Println("Error checking uniqueness:", err)
				return c.JSON(http.StatusInternalServerError, &models.Response{
					Code:    500,
					Message: "Terjadi kesalahan internal pada server.",
					Status:  false,
				})
			}

			if !isUnique {
				log.Println("Username atau email telah digunakan oleh data lain.")
				return c.JSON(http.StatusBadRequest, &models.Response{
					Code:    400,
					Message: "Username atau email telah digunakan oleh data lain!",
					Status:  false,
				})
			}
		}
		// if updateUserAppRole.Username == existingUser.UserName || updateUserAppRole.Email == existingUser.UserEmail {
		// 	// Data username atau email berubah, lakukan validasi unik
		// 	var existingUserID int
		// 	err := db.QueryRow("SELECT user_id FROM user_ms WHERE (user_name = $1 OR user_email = $2) AND deleted_at IS NULL", updateUserAppRole.Username, updateUserAppRole.Email).Scan(&existingUserID)

		// 	if err == nil {
		// 		// User dengan username atau email yang sama sudah ada
		// 		log.Print(err)
		// 		return c.JSON(http.StatusBadRequest, &models.Response{
		// 			Code:    400,
		// 			Message: "Username atau email telah digunakan!",
		// 			Status:  false,
		// 		})
		// 	} else if err != sql.ErrNoRows {
		// 		// Terjadi kesalahan lain saat menjalankan query
		// 		log.Print(err)
		// 		return c.JSON(http.StatusInternalServerError, &models.Response{
		// 			Code:    500,
		// 			Message: "Terjadi kesalahan internal pada server.",
		// 			Status:  false,
		// 		})
		// 	}
		// }

		//beda
		// // Bandingkan nilai username dan user_email dengan yang baru diterima
		// if existingUser.UserName != updateUserAppRole.Username || existingUser.UserEmail != updateUserAppRole.Email {
		// 	// Jika nilai berubah, tolak pembaruan
		// 	return c.JSON(http.StatusBadRequest, &models.Response{
		// 		Code:    400,
		// 		Message: "Username atau email telah digunakan!",
		// 		Status:  false,
		// 	})
		// }

		_, addroleErr := service.UpdateUserAppRole(updateUserAppRole, userApplicationRoleUUID)
		if addroleErr != nil {
			log.Printf("Error updating user application role: %v", addroleErr)
			return c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Terjadi kesalahan internal pada server.",
				Status:  false,
			})
		}

		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil mengupdate akun!",
			Status:  true,
		})
	} else {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi",
			Status:  false,
		})
	}
}

func DeleteUserAppRole(c echo.Context) error {
	tokenString := c.Request().Header.Get("Authorization")
	secretKey := "secretJwToken"

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

	//dekripsi token JWE
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

	id := c.Param("id")

	if err == nil {
		errService := service.DeleteUserAppRole(id, userUUID)
		if errService == service.ErrNotFound {
			log.Printf("error kenapa %s", errService)
			return c.JSON(http.StatusNotFound, &models.Response{
				Code:    404,
				Message: "Gagal menghapus user. User tidak ditemukan!",
				Status:  false,
			})
		}

		return c.JSON(http.StatusOK, &models.Response{
			Code:    200,
			Message: "User berhasil dihapus!",
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
