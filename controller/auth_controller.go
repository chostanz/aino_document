package controller

import (
	"aino_document/models"
	"aino_document/service"
	"aino_document/utils"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/lib/pq"

	_ "github.com/dgrijalva/jwt-go"
)

type TokenCheck struct {
	Token string `json:"token"`
}

type JwtCustomClaims struct {
	UserID   int    `json:"user_id"`
	UserUUID string `json:"user_uuid"`
	RoleCode string `json:"role_code"`
	// AppRoleId          int `json:"application_role_id"`
	DivisionCode       string `json:"division_code"`
	jwt.StandardClaims        // Embed the StandardClaims struct

}

// Simpan token yang tidak valid dalam bentuk set
var InvalidTokens = make(map[string]struct{})

func RegisterUser(c echo.Context) error {
	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}

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
	userUUID := c.Get("user_uuid").(string)
	_, errK := service.GetUserInfoFromToken(tokenOnly)
	if errK != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid token atau token tidak ditemukan!")
	}

	var userRegister models.Register

	if errBind := c.Bind(&userRegister); errBind != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

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
		registerErr := service.RegisterUser(userRegister, userUUID)
		if registerErr != nil {
			if validationErr, ok := registerErr.(*service.ValidationError); ok {
				if validationErr.Tag == "strong_password" {
					return c.JSON(http.StatusUnprocessableEntity, &models.Response{
						Code:    422,
						Message: "Password harus memiliki setidaknya 8 karakter",
						Status:  false,
					})
				}
			} else if dbErr, ok := registerErr.(*pq.Error); ok {
				// Check for duplicate key violation (unique constraint violation)
				if dbErr.Code.Name() == "unique_violation" {
					log.Println(dbErr)
					return c.JSON(http.StatusBadRequest, &models.Response{
						Code:    400,
						Message: "Username atau email telah digunakan!",
						Status:  false,
					})
				}
			}
		}
		// log.Print(registerErr)
		return c.JSON(http.StatusCreated, &models.Response{
			Code:    201,
			Message: "Berhasil membuat akun!",
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
	if errBind := c.Bind(&loginbody); errBind != nil {
		return c.JSON(http.StatusBadRequest, &models.Response{
			Code:    400,
			Message: "Data tidak valid!",
			Status:  false,
		})
	}

	err := c.Validate(&loginbody)

	if err != nil {
		log.Println(err)
		return c.JSON(http.StatusInternalServerError, &models.Response{
			Code:    500,
			Message: "Terjadi kesalahan pada internal server. Coba beberapa saat lagi!",
			Status:  false,
		})
	}

	user_uuid, role_code, division_code, user_id, isAuthentication, _ := service.Login(loginbody) //bagian siniii dikasih role_id ama yg laen

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
		UserUUID: user_uuid,
		UserID:   user_id,
		RoleCode: role_code,
		// AppRoleId:  application_role_id,
		DivisionCode: division_code,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 12).Unix(), // Tambahkan waktu kadaluwarsa (15 menit)
		},
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// secretKey := []byte("secretKeyJWT")
	// t, err := token.SignedString(secretKey)
	// fmt.Println("token:", t)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, &models.ResponseLogin{
	// 		Code:    500,
	// 		Message: "Gagal membuat token. Mohon coba beberapa saat lagi!",
	// 		Status:  false,
	// 	})
	// }

	// Mengonversi claims ke format JSON
	claimData, err := json.Marshal(claims)
	if err != nil {
		fmt.Println("Gagal mengonversi klaim:", err)
	}

	// Buat token JWT dengan enkripsi JWE
	secretKey := "secretJwToken" // Ganti dengan kunci
	jweToken, err := jose.Encrypt(string(claimData), jose.PBES2_HS256_A128KW, jose.A128GCM, secretKey)
	if err != nil {
		fmt.Println("Gagal membuat token:", err)
	}

	fmt.Println("Token JWE:", jweToken)

	// Decode token JWE
	// decodedToken, _, err := jose.Decode(jweToken, secretKey)
	// if err != nil {
	// 	fmt.Println("Gagal mendekripsi token:", err)
	// }

	// // Tampilkan token yang telah dideskripsi
	// fmt.Println("Token yang telah dideskripsi:", decodedToken)
	return c.JSON(http.StatusOK, &models.ResponseLogin{
		Code:    200,
		Message: "Berhasil login",
		Status:  true,
		Token:   jweToken,
	})

}

func Logout(c echo.Context) error {

	defer func() {
		if r := recover(); r != nil {
			log.Println("Panic occurred:", r)
			c.JSON(http.StatusInternalServerError, &models.Response{
				Code:    500,
				Message: "Internal Server Error",
				Status:  false,
			})
		}
	}()
	tokenString := c.Request().Header.Get("Authorization")

	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Code:    401,
			Message: "Token tidak ditemukan",
			Status:  false,
		})
	}
	_, exists := InvalidTokens[tokenString]
	if exists {
		return c.JSON(http.StatusUnauthorized, &models.Response{
			Code:    401,
			Message: "Token tidak valid atau Anda telah logout",
			Status:  false,
		})
	}

	InvalidTokens[tokenString] = struct{}{}

	return c.JSON(http.StatusOK, &models.Response{
		Code:    200,
		Message: "Berhasil Logout!",
		Status:  true,
	})
}

func UpdateUser(c echo.Context) error {
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
	userUUID := c.Get("user_uuid").(string)
	_, errK := service.GetUserInfoFromToken(tokenOnly)
	if errK != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid token atau token tidak ditemukan!")
	}

	id := userUUID

	var updateUser models.UpdateUser
	if err := c.Bind(&updateUser); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"code":    400,
			"message": "Data invalid!",
			"status":  false,
		})
	}

	// Validasi data pembaruan
	if err := c.Validate(&updateUser); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    422,
			"message": "Data tidak boleh kosong!",
			"status":  false,
		})
	}

	// Panggil fungsi UpdateUserProfile dari service
	if err := service.UpdateUserProfile(updateUser, id, userUUID); err != nil {
		log.Println("Error updating user profile:", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"code":    500,
			"message": "Terjadi kesalahan internal pada server. Mohon coba beberapa saat lagi",
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": "Profil pengguna telah diperbarui!",
		"status":  true,
	})
}
