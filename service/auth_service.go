package service

import (
	"aino_document/models"
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ValidationError struct {
	Message string
	Field   string
	Tag     string
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

func RegisterUser(userRegister models.Register, userID int) error {
	if len(userRegister.Password) < 8 {
		return &ValidationError{
			Message: "Password should be of 8 characters long",
			Field:   "password",
			Tag:     "strong_password",
		}
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRegister.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	hashedPasswordStr := base64.StdEncoding.EncodeToString(hashedPassword)
	fmt.Println(hashedPassword)
	fmt.Println(hashedPasswordStr)

	username, errP := GetUsernameByID(userID)
	if errP != nil {
		return errP
	}
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	// uuid := uuid.New()
	// user_id := uuid.ID()

	user_id := currentTimestamp + int64(uniqueID)
	uuid := uuid.New()
	uuidString := uuid.String()
	_, errInsert := db.NamedExec("INSERT INTO user_ms (user_id, user_uuid, user_name, user_email, user_password, created_by) VALUES (:user_id, :user_uuid, :user_name, :user_email, :user_password, :created_by)", map[string]interface{}{
		"user_id":       user_id,
		"user_uuid":     uuidString,
		"user_name":     userRegister.Username,
		"user_email":    userRegister.Email,
		"user_password": hashedPasswordStr,
		"created_by":    username,
	})

	if errInsert != nil {
		return errInsert
	}

	err = db.Get(&user_id, "SELECT user_id FROM user_ms WHERE user_name = $1", userRegister.Username)
	if err != nil {
		return err
	}
	// Cetak nilai user_id dari middleware
	fmt.Println("Middleware UserID:", userID)

	// Cetak nilai user_id yang ingin dimasukkan ke dalam database
	fmt.Println("Database UserID:", user_id)

	var applicationRoles []models.UserAppRole
	// Mengambil application_role_id dari database
	rows, err := db.Query("SELECT application_role_id FROM application_role_ms WHERE application_id = $1", userRegister.ApplicationRole.Application_id)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var applicationRole models.UserAppRole
		err := rows.Scan(&applicationRole.Application_role_id)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return err
		}
		applicationRoles = append(applicationRoles, applicationRole)
	}

	// Mengambil division_id dari database
	rowsDivision, err := db.Query("SELECT division_id FROM division_ms WHERE division_code = $1", userRegister.ApplicationRole.Division_code)
	if err != nil {
		log.Println("gabisa ambil division_id", err)
		return err
	}
	defer rowsDivision.Close()

	for rowsDivision.Next() {
		err := rowsDivision.Scan(&applicationRoles[0].Division_id) // Asumsikan Anda ingin mengisi division_id ke dalam slice pertama saja
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return err
		}
	}

	for _, applicationRole := range applicationRoles {
		_, err := db.Exec("INSERT INTO user_application_role_ms(user_id, application_role_id, division_id) VALUES ($1, $2, $3)", user_id, applicationRole.Application_role_id, applicationRole.Division_id)
		if err != nil {
			log.Println("Error inserting data into user_application_role_ms:", err)
			return err
		}
	}
	return nil
}

func Login(userLogin models.Login) (int, bool, error) {
	var isAuthentication bool
	var user_id int
	// var application_role_id int
	// var division_id int

	rows, err := db.Query("SELECT CASE WHEN COUNT(*) > 0 THEN 'true' ELSE 'false' END FROM user_ms WHERE user_email = $1 AND user_password = $2", userLogin.Email, userLogin.Password)
	if err != nil {
		return 0, false, err
	}

	defer rows.Close()

	rows, err = db.Query("SELECT user_id, user_password from user_ms where user_email = $1", userLogin.Email)
	if err != nil {
		fmt.Println("Error querying users:", err)
		return 0, false, err
	}

	defer rows.Close()

	var dbPasswordBase64 string
	if rows.Next() {
		err = rows.Scan(&user_id, &dbPasswordBase64)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return 0, false, err
		}
		dbPassword, errBycript := base64.StdEncoding.DecodeString(dbPasswordBase64)

		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return 0, false, errBycript
		}
		errBycript = bcrypt.CompareHashAndPassword(dbPassword, []byte(userLogin.Password))
		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return 0, false, errBycript
		}
		isAuthentication = true
	}

	if isAuthentication {
		rows, err = db.Query("SELECT application_role_id, division_id FROM user_application_role_ms WHERE user_id = $1", user_id)
		if err != nil {
			fmt.Println("Error querying user roles:", err)
			return 0, false, err
		}
		defer rows.Close()

		// if rows.Next() {
		// 	err = rows.Scan(&application_role_id, &division_id)
		// 	if err != nil {
		// 		fmt.Println("Error scanning role row:", err)
		// 		return 0, false, 0, 0, err
		// 	}
		// }
		return user_id, isAuthentication, nil
	}
	return 0, false, nil // Jika tidak ada authentikasi yang berhasil

}
