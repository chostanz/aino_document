package service

import (
	"aino_document/models"
	"encoding/base64"
	"errors"
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

func RegisterUser(userRegister models.Register, userUUID string) error {
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

	username, errP := GetUsernameByID(userUUID)
	if errP != nil {
		return errP
	}
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
		"created_by":    username,
	})

	if errInsert != nil {
		return errInsert
	}

	err = db.Get(&user_id, "SELECT user_id FROM user_ms WHERE user_name = $1", userRegister.Username)
	if err != nil {
		return err
	}

	// Mendapatkan role_id yang baru saja diinsert
	var roleID int64
	err = db.Get(&roleID, "SELECT role_id FROM role_ms WHERE role_uuid = $1 AND deleted_at IS NULL", userRegister.ApplicationRole.Role_UUID)
	if err != nil {
		log.Println("Error getting role_id:", err)
		return err
	}
	var applicationID int64
	err = db.Get(&applicationID, "SELECT application_id FROM application_ms WHERE application_uuid = $1 AND deleted_at IS NULL", userRegister.ApplicationRole.Application_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return err
	}

	// Get division_id
	var divisionID int64
	err = db.Get(&divisionID, "SELECT division_id FROM division_ms WHERE division_uuid = $1 AND deleted_at IS NULL", userRegister.ApplicationRole.Division_UUID)
	if err != nil {
		log.Println("Error fetching division_id:", err)
		return err
	}

	AppRoleId := currentTimestamp + int64(uniqueID)
	// Insert data ke application_role_ms
	_, err = db.Exec("INSERT INTO application_role_ms(application_role_id, application_id, role_id, created_by) VALUES ($1, $2, $3, $4)",
		AppRoleId, applicationID, roleID, username)
	if err != nil {
		log.Println("Error inserting data into application_role_ms:", err)
		return err
	}
	log.Println("Data inserted into application_role_ms successfully")

	// Get application_role_id
	var applicationRoleID int64
	err = db.Get(&applicationRoleID, "SELECT application_role_id FROM application_role_ms WHERE application_id = $1 AND role_id = $2",
		applicationID, roleID)
	if err != nil {
		log.Println("Error getting application_role_id:", err)
		return err
	}
	log.Println("Application Role ID:", applicationRoleID)

	// Insert user_application_role_ms data
	_, err = db.Exec("INSERT INTO user_application_role_ms(user_application_role_uuid, user_id, application_role_id, division_id, created_by) VALUES ($1, $2, $3, $4, $5)", uuidString, user_id, applicationRoleID, divisionID, username)
	if err != nil {
		log.Println("Error inserting data into user_application_role_ms:", err)
		return err
	}

	return nil
}

func Login(userLogin models.Login) (string, string, string, int, bool, error) {
	var isAuthentication bool
	var user_id int
	var user_uuid string
	var role_code string
	// var application_role_id int
	var division_code string

	rows, err := db.Query("SELECT CASE WHEN COUNT(*) > 0 THEN 'true' ELSE 'false' END FROM user_ms WHERE user_email = $1 AND user_password = $2", userLogin.Email, userLogin.Password)
	if err != nil {
		return "", "", "", 0, false, err
	}

	defer rows.Close()

	rows, err = db.Query("SELECT user_uuid, user_password from user_ms where user_email = $1", userLogin.Email)
	if err != nil {
		fmt.Println("Error querying users:", err)
		return "", "", "", 0, false, err
	}

	defer rows.Close()

	var dbPasswordBase64 string
	if rows.Next() {
		err = rows.Scan(&user_uuid, &dbPasswordBase64)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			return "", "", "", 0, false, err
		}
		dbPassword, errBycript := base64.StdEncoding.DecodeString(dbPasswordBase64)

		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return "", "", "", 0, false, errBycript
		}
		errBycript = bcrypt.CompareHashAndPassword(dbPassword, []byte(userLogin.Password))
		if errBycript != nil {
			fmt.Println("Password comparison failed:", errBycript)
			return "", "", "", 0, false, errBycript
		}
		isAuthentication = true
	}

	if isAuthentication {
		// Query untuk mendapatkan division_code
		rows, err := db.Query("SELECT d.division_code FROM division_ms d JOIN user_application_role_ms uar ON d.division_id = uar.division_id JOIN user_ms u ON uar.user_id = u.user_id WHERE u.user_uuid = $1", user_uuid)
		if err != nil {
			fmt.Println("Error querying division code:", err)
			return "", "", "", 0, false, err
		}
		defer rows.Close()

		// Periksa hasil query division_code
		if rows.Next() {
			err = rows.Scan(&division_code)
			if err != nil {
				fmt.Println("Error scanning division_code row:", err)
				return "", "", "", 0, false, err
			}
		}

		// Query untuk mendapatkan role_code
		rows, err = db.Query("SELECT r.role_code FROM role_ms r JOIN application_role_ms ar ON r.role_id = ar.role_id JOIN user_application_role_ms uar ON ar.application_role_id = uar.application_role_id JOIN user_ms u ON u.user_id = uar.user_id WHERE u.user_uuid = $1", user_uuid)
		if err != nil {
			fmt.Println("Error querying role code:", err)
			return "", "", "", 0, false, err
		}
		defer rows.Close()

		// Periksa hasil query role_code
		if rows.Next() {
			err = rows.Scan(&role_code)
			if err != nil {
				fmt.Println("Error scanning role_code row:", err)
				return "", "", "", 0, false, err
			}
		}

		// if rows.Next() {
		// 	err = rows.Scan(&application_role_id, &division_id)
		// 	if err != nil {
		// 		fmt.Println("Error scanning role row:", err)
		// 		return 0, false, 0, 0, err
		// 	}
		// }
		return user_uuid, role_code, division_code, user_id, isAuthentication, nil
	}
	return "", "", "", 0, false, nil // Jika tidak ada authentikasi yang berhasil

}

func UpdateUserProfile(userUpdate models.UpdateUser, id string, userUUID string) error {
	var userID int64
	err := db.Get(&userID, "SELECT user_id FROM user_ms WHERE user_uuid = $1", userUUID)
	if err != nil {
		log.Println("Error getting user ID:", err)
		return err
	}

	// Update nama pengguna dan email di tabel user_ms
	_, err = db.NamedExec("UPDATE user_ms SET user_name = :user_name, user_email = :user_email WHERE user_uuid = :user_uuid",
		map[string]interface{}{
			"user_uuid":  userUUID,
			"user_name":  userUpdate.Username,
			"user_email": userUpdate.Email,
		})
	if err != nil {
		log.Println("Error updating user profile in user_ms:", err)
		return err
	}

	// Dapatkan role_id, application_id, dan division_id yang sudah ada
	var roleID, applicationID, applicationUUID, divisionID int64
	err = db.Get(&roleID, "SELECT role_id FROM role_ms WHERE role_uuid = $1", userUpdate.ApplicationRole.Role_UUID)
	if err != nil {
		log.Println("Error getting role_id:", err)
		return err
	}

	err = db.Get(&applicationID, "SELECT application_id FROM application_ms WHERE application_uuid = $1", userUpdate.ApplicationRole.Application_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return err
	}

	err = db.Get(&divisionID, "SELECT division_id FROM division_ms WHERE division_uuid = $1", userUpdate.ApplicationRole.Division_UUID)
	if err != nil {
		log.Println("Error fetching division_id:", err)
		return err
	}

	err = db.Get(&applicationUUID, "SELECT application_role_uuid FROM application_ms WHERE application_role_uuid = $1", userUpdate.ApplicationRole.Application_role_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return err
	}

	// Update data di tabel application_role_ms
	_, err = db.Exec("UPDATE application_role_ms SET application_id = $1, role_id = $2 WHERE application_uuid = $3",
		applicationID, roleID, applicationUUID)
	if err != nil {
		log.Println("Error updating data in application_role_ms:", err)
		return err
	}

	// Update data di tabel user_application_role_ms
	_, err = db.Exec("UPDATE user_application_role_ms SET division_id = $1 WHERE user_id = $2",
		divisionID, userID)
	if err != nil {
		log.Println("Error updating data in user_application_role_ms:", err)
		return err
	}

	return nil
}

func ChangePassword(changePassword models.ChangePasswordRequest, userUUID string) error {
	var dbPassword string
	err := db.Get(&dbPassword, "SELECT user_password FROM user_ms WHERE user_uuid = $1", userUUID)
	if err != nil {
		return err
	}

	decodedPassword, err := base64.StdEncoding.DecodeString(dbPassword)
	if err != nil {
		return err
	}

	if len(changePassword.NewPassword) < 8 {
		return &ValidationError{
			Message: "Password should be of 8 characters long",
			Field:   "password",
			Tag:     "strong_password",
		}
	}
	errBycript := bcrypt.CompareHashAndPassword(decodedPassword, []byte(changePassword.OldPassword))
	if errBycript != nil {
		fmt.Println("Error comparing old passwords:", errBycript)
		return errBycript
	}

	if changePassword.OldPassword == changePassword.NewPassword {
		return errors.New("new password must be different from old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(changePassword.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("Error generating hashed password:", err)
		return err
	}

	hashedPasswordStr := base64.StdEncoding.EncodeToString(hashedPassword)
	_, err = db.NamedExec("UPDATE user_ms SET user_password = :user_password WHERE user_uuid = :user_uuid", map[string]interface{}{
		"user_password": hashedPasswordStr,
		"user_uuid":     userUUID,
	})

	if err != nil {
		fmt.Println("Error updating password in database:", err)
		return err
	}
	return nil
}
