package service

import (
	"aino_document/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
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

func GetUsernameByID(userID int) (string, error) {
	var username string
	err := db.QueryRow("SELECT user_name from user_ms WHERE user_id = $1", userID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func GetUserInfoFromToken(tokenStr string) (int, error) {
	secretKey := "secretJwToken" // Ganti dengan kunci yang benar

	decrypted, err := DecryptJWE(tokenStr, secretKey)
	if err != nil {
		fmt.Println("Gagal mendekripsi token:", err)
		return 0, err
	}

	fmt.Println("Token yang sudah dideskripsi:", decrypted)

	var claims JwtCustomClaims
	errJ := json.Unmarshal([]byte(decrypted), &claims)
	if errJ != nil {
		fmt.Println("Gagal mengurai klaim:", errJ)
		return 0, errJ
	}

	userID := claims.UserId // Mengakses UserID langsung
	return userID, nil
}

func AddApplication(addApplication models.Application, userID int) error {
	username, errP := GetUsernameByID(userID)
	if errP != nil {
		return errP
	}
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	app_id := currentTimestamp + int64(uniqueID)

	_, err := db.NamedExec("INSERT INTO application_ms (application_id, application_code, application_title, application_description, created_by) VALUES (:application_id, :application_code, :application_title, :application_description, :created_by)", map[string]interface{}{
		"application_id":          app_id,
		"application_code":        addApplication.Code,
		"application_title":       addApplication.Title,
		"application_description": addApplication.Description,
		"created_by":              username,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetAllApp() ([]models.Applications, error) {
	application := []models.Applications{}

	rows, err := db.Queryx("SELECT application_order, application_code, application_title, application_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM application_ms")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		place := models.Applications{}
		rows.StructScan(&place)
		application = append(application, place)
	}
	return application, nil
}

func ShowApplicationById(id int) (models.Applications, error) {
	var appid models.Applications
	idStr := strconv.Itoa(id)

	err := db.Get(&appid, "SELECT application_order, application_code, application_title, application_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM application_ms WHERE application_order = $1", idStr)
	if err != nil {
		return models.Applications{}, err
	}
	return appid, nil

}

func GetAppById(id int) (models.Application, error) {
	var appid models.Application
	idStr := strconv.Itoa(id)

	err := db.Get(&appid, "SELECT * FROM application_ms WHERE application_order = $1", idStr)
	if err != nil {
		return models.Application{}, err
	}
	return appid, nil

}
