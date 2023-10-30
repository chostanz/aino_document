package service

import (
	"aino_document/models"
	"encoding/base64"
	"fmt"
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

func RegisterUser(userRegister models.Register) error {
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

	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	user_id := currentTimestamp + int64(uniqueID)
	uuid := uuid.New()
	uuidString := uuid.String()
	_, errInsert := db.NamedExec("INSERT INTO user_ms (user_id, user_uuid, user_name, user_email, user_password, created_by) VALUES (:user_id, :user_uuid, :user_name, :user_email, :user_password, :created_by)", map[string]interface{}{
		"user_id":       user_id,
		"user_uuid":     uuidString,
		"user_name":     userRegister.Username,
		"user_email":    userRegister.Email,
		"user_password": hashedPasswordStr,
		"created_by":    userRegister.Created_by,
	})

	if errInsert != nil {
		return errInsert
	}

	err = db.Get(&user_id, "SELECT user_id FROM user_ms WHERE user_name = $1", userRegister.Username)
	if err != nil {
		return err
	}

	var applicationRole models.UserAppRole
	_, err = db.Exec("INSERT INTO user_application_role_ms(user_id, application_role_id, division_id) VALUES ($1, $2, $3)", user_id, applicationRole.Application_role_id, applicationRole.Division_id)
	if err != nil {
		return err
	}

	return nil
}

func Login(userLogin models.Login) (int, bool, int, int, error) {
	var isAuthentication bool
	var user_id int
	var application_role_id int
	var division_id int

	rows, err := db.Query("SELECT CASE WHEN COUNT(*) > 0 THEN 'true' ELSE 'false' END FROM user_ms WHERE user_name = $1 AND user_password = $2", userLogin.Username, userLogin.Password)
	if err != nil {
		return 0, false, 0, 0, err
	}

	defer rows.Close()

	rows, err = db.Query("SELECT user_id, user_password from user_ms where user_name = $1", userLogin.Username)
	if err != nil {
		fmt.Println("Error querying users:", err)
		return 0, false, 0, 0, err
	}

	defer rows.Close()

	var dbPasswordBase64 string
	if rows.Next() {
		err = rows.Scan(&user_id, &dbPasswordBase64)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return 0, false, 0, 0, err
		}
		dbPassword, errBycript := base64.StdEncoding.DecodeString(dbPasswordBase64)

		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return 0, false, 0, 0, errBycript
		}
		errBycript = bcrypt.CompareHashAndPassword(dbPassword, []byte(userLogin.Password))
		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return 0, false, 0, 0, errBycript
		}
		isAuthentication = true
	}

	if isAuthentication {
		rows, err = db.Query("SELECT application_role_id, division_id FROM user_application_role_ms WHERE user_id = $1", user_id)
		if err != nil {
			fmt.Println("Error querying user roles:", err)
			return 0, false, 0, 0, err
		}
		defer rows.Close()

		if rows.Next() {
			err = rows.Scan(&application_role_id, &division_id)
			if err != nil {
				fmt.Println("Error scanning role row:", err)
				return 0, false, 0, 0, err
			}
		}
		return user_id, isAuthentication, application_role_id, division_id, nil
	}
	return 0, false, 0, 0, nil // Jika tidak ada authentikasi yang berhasil

}
