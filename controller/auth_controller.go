package controller

import (
	"aino_document/models"
	"aino_document/service"
	"aino_document/utils"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/badoux/checkmail"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	_ "github.com/dgrijalva/jwt-go"
)

type TokenCheck struct {
	Token string `json:"token"`
}

type JwtCustomClaims struct {
	UserId             int `json:"user_id"`
	AppRoleId          int `json:"application_role_id"`
	DivisionId         int `json:"division_id"`
	jwt.StandardClaims     // Embed the StandardClaims struct

}

func RegisterUser(c echo.Context) error {
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	var userRegister models.Register

	c.Bind(&userRegister)
	errEmail := checkmail.ValidateFormat(userRegister.Email)
	if errEmail != nil {
		return c.JSON(http.StatusUnprocessableEntity, &models.Response{
			Code:    422,
			Message: "Format email tidak valid",
			Status:  false,
		})
	}

	err := c.Validate(&userRegister)

	if err == nil {
		registerErr := service.RegisterUser(userRegister)
		if registerErr != nil {

			if validationErr, ok := registerErr.(*service.ValidationError); ok {
				if validationErr.Tag == "strong_password" {
					return c.JSON(http.StatusUnprocessableEntity, &models.Response{
						Code:    422,
						Message: "Password harus memiliki setidaknya 8 karakter",
						Status:  false,
					})
				}
			}
			log.Print(registerErr)
			return c.JSON(http.StatusBadRequest, &models.Response{
				Code:    400,
				Message: "Username atau email telah digunakan!",
				Status:  false,
			})
		}
		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil register",
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

func Login(c echo.Context) error {
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

	var loginbody models.Login
	c.Bind(&loginbody)

	err := c.Validate(&loginbody)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan pada internal server. Coba beberapa saat lagi!",
			Status:  false,
		})
	}

	user_id, isAuthentication, application_role_id, division_id, _ := service.Login(loginbody)

	fmt.Println("isAuthentication:", isAuthentication)

	if !isAuthentication {
		fmt.Println("Authentication failed")

		return c.JSON(http.StatusUnauthorized, &models.Response{
			Code:    401,
			Message: "Akun tidak ada atau password salah",
			Status:  false,
		})
	}
	claims := &JwtCustomClaims{
		UserId:     user_id,
		AppRoleId:  application_role_id,
		DivisionId: division_id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(), // Tambahkan waktu kadaluwarsa (15 menit)
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte("rahasia"))
	fmt.Println("token:", t)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.ResponseLogin{
			Code:    500,
			Message: "Gagal membuat token. Mohon coba beberapa saat lagi!",
			Status:  false,
		})
	}

	return c.JSON(http.StatusOK, &models.ResponseLogin{
		Code:    200,
		Message: "Berhasil login",
		Status:  true,
		Token:   t,
	})

}

func Logout(c echo.Context) error {
	// token := c.Request().Header.Get("Authorization") // Anda harus menyesuaikan ini sesuai dengan bagaimana token disimpan

	// if token == "" {
	// 	return c.JSON(http.StatusUnauthorized, map[string]string{
	// 		"message": "Token not found",
	// 	})
	// }

	claims := &JwtCustomClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix(),
		},
	}

	// Buat token yang sudah kadaluwarsa
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	_, err := expiredToken.SignedString([]byte("rahasia"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Failed to create an invalid token. Please try again later.",
			Status:  false,
		})
	}

	// c.Response().Header().Set("Authorization", tokenString)

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Logout successful",
		Status:  true,
	})
}
