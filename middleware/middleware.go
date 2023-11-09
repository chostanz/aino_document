package middleware

import (
	"aino_document/controller"
	"aino_document/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserId int `json:"user_id"`
	// AppRoleId          int `json:"application_role_id"`
	// DivisionId         int `json:"division_id"`
	jwt.StandardClaims // Embed the StandardClaims struct

}

func DecryptJWE(jweToken string, secretKey string) (string, error) {
	// Dekripsi token JWE
	decrypted, _, err := jose.Decode(jweToken, secretKey)
	if err != nil {
		return "", err
	}
	return decrypted, nil
}

func ExtractClaims(jwtToken string) (JwtCustomClaims, error) {
	claims := &JwtCustomClaims{}
	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

	token, err := jwt.ParseWithClaims(jwtToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return JwtCustomClaims{}, err
	}

	return *claims, nil
}
func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		fmt.Println("Token yang sudah dideskripsi:", decrypted)

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

		// Sekarang Anda memiliki data dalam struct JwtCustomClaims
		// Anda bisa mengakses UserId atau klaim lain sesuai kebutuhan
		// fmt.Println("UserID:", claims.UserId)

		userID := claims.UserId // Mengakses UserID langsung
		fmt.Println("UserID:", userID)
		if userID != 1698632900865066 {
			return c.JSON(http.StatusForbidden, &models.Response{
				Code:    403,
				Message: "Akses ditolak!",
				Status:  false,
			})
		}

		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}

func AdmiknMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		fmt.Println("Token yang sudah dideskripsi:", decrypted)

		// Langkah 2: Mengakses klaim dalam token JWT yang telah didekripsi
		userID, err := ExtractClaims(decrypted)
		if err != nil {
			fmt.Println("Gagal mengakses klaim:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		fmt.Println("UserID:", userID)

		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}

func AdminMiuddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		// Dekripsi token JWE
		decrypted, _, err := jose.Decode(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}
		// Parse token JWT
		claims := &JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(tokenOnly, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			fmt.Println(err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token invalid!",
				"status":  false,
			})
		}

		// Token JWE valid, ID pengguna ada dalam claims
		c.Set("user_id", claims.UserId)

		fmt.Println("token yg sudah dideskripsi", decrypted)
		// Token JWE valid, Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}

func AdmimnMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		// Decode token JWE
		decodedToken, _, err := jose.Decode(tokenOnly, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		// Ubah klaim token JWE menjadi klaim token JWT
		jweClaims := map[string]interface{}{}
		err = json.Unmarshal([]byte(decodedToken), &jweClaims)
		if err != nil {
			fmt.Println("Gagal mengurai klaim JWE:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		claims := &JwtCustomClaims{
			UserId: int(jweClaims["user_id"].(float64)), // Ubah sesuai kebutuhan
		}

		// Tanda tangani token dengan kunci rahasia
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err = token.SignedString([]byte(secretKey))
		if err != nil {
			fmt.Println("Gagal membuat token JWT:", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}

		fmt.Println("Token JWT:", tokenString)

		// Sekarang Anda memiliki token JWT yang dapat digunakan dalam aplikasi Anda

		// Anda dapat melanjutkan dengan pengolahan berikutnya
		return next(c)
	}
}

func AdminnMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		secretKey := "secretJwToken" // Ganti dengan kunci
		decodedToken, _, err := jose.Decode(tokenString, secretKey)
		if err != nil {
			fmt.Println("Gagal mendekripsi token:", err)
		}

		tokenSplit := strings.Split(decodedToken, " ")

		if len(tokenSplit) != 2 || tokenSplit[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Token tidak valid!",
				"status":  false,
			})
		}
		tokenOnly := tokenSplit[1]
		fmt.Println("Token:", tokenOnly)
		if tokenString == "" {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Code:    401,
				Message: "Token tidak ditemukan!",
				Status:  false,
			})
		}
		// token, err := jwt.Parse(tokenOnly, func(token *jwt.Token) (interface{}, error) {
		// 	return []byte("rahasia"), nil
		// })
		claims := &JwtCustomClaims{}
		token, err := jwt.ParseWithClaims(tokenOnly, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("secretJwToken"), nil // Ganti dengan kunci yang benar
		})

		if err != nil || !token.Valid {
			log.Print(err)
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Code:    401,
				Message: "Token invalid!",
				Status:  false,
			})
		}

		// claims, ok := token.Claims.(jwt.MapClaims)
		// if !ok {
		// 	fmt.Println("Tipe asli dari token.Claims:", reflect.TypeOf(token.Claims))
		// 	return c.JSON(http.StatusUnauthorized, &models.Response{
		// 		Code:    401,
		// 		Message: "Token claims tidak valid!",
		// 		Status:  false,
		// 	})
		// }

		fmt.Println(claims)

		// Check apakah token ada di invalidTokens
		if _, exists := controller.InvalidTokens[token.Raw]; exists {
			return c.JSON(http.StatusUnauthorized, &models.Response{
				Code:    401,
				Message: "Sesi berakhir! Silahkan login kembali",
				Status:  false,
			})
		}

		c.Set("users", token)

		// roleID := int(claims["id_role"].(float64))
		// if roleID != 1 {
		// 	return c.JSON(http.StatusForbidden, &models.Response{
		// 		Code:    403,
		// 		Message: "Akses ditolak!",
		// 		Status:  false,
		// 	})
		// }

		return next(c)
	}
}
