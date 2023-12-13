package service

import (
	"aino_document/models"
	"database/sql"
	"fmt"
	"log"
	"time"
)

func GetUserApplicationRole() ([]models.Users, error) {
	getUserAppRole := []models.Users{}

	rows, err := db.Queryx("SELECT u.user_uuid, uar.user_application_role_uuid, u.user_name, u.user_email, r.role_title, a.application_title, d.division_title, uar.created_by, uar.created_at, uar.updated_by, uar.updated_at FROM user_ms u INNER JOIN user_application_role_ms uar ON u.user_id = uar.user_id INNER JOIN application_role_ms ar ON uar.application_role_id = ar.application_role_id INNER JOIN application_ms a ON ar.application_id = a.application_id INNER JOIN role_ms r ON ar.role_id = r.role_id INNER JOIN division_ms d ON uar.division_id = d.division_id WHERE uar.deleted_at IS NULL")
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}

	for rows.Next() {
		place := models.Users{}
		err := rows.StructScan(&place)
		if err != nil {
			log.Println("Error scanning row to struct:", err)
			continue
		}
		getUserAppRole = append(getUserAppRole, place)
	}

	return getUserAppRole, nil
}

func GetSpecUseApplicationRole(id string) (models.Users, error) {
	var users models.Users
	err := db.Get(&users, "SELECT u.user_uuid, uar.user_application_role_uuid, u.user_name, u.user_email, r.role_title, a.application_title, d.division_title,  uar.created_by, uar.created_at, uar.updated_by, uar.updated_at FROM user_ms u INNER JOIN user_application_role_ms uar ON u.user_id = uar.user_id INNER JOIN application_role_ms ar ON uar.application_role_id = ar.application_role_id INNER JOIN application_ms a ON ar.application_id = a.application_id INNER JOIN role_ms r ON ar.role_id = r.role_id INNER JOIN division_ms d ON uar.division_id = d.division_id WHERE uar.user_application_role_uuid = $1 AND uar.deleted_at IS NULL", id)
	if err != nil {
		return models.Users{}, err
	}

	return users, nil

}

func GetUsernameByUserAppRoleUUID(userAppRoleUUID string) (string, error) {
	var user_uuid string

	err := db.Get(&user_uuid, "SELECT u.user_uuid FROM user_ms u JOIN user_application_role_ms uarm ON u.user_id = uarm.user_id WHERE uarm.user_application_role_uuid = $1", userAppRoleUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found for user_application_role_uuid:")
			return "", err
		}

		log.Println("Error getting user uuid by user_application_role_uuid:", err)
		return "", err
	}

	return user_uuid, nil
}

func GetUsernameByIDUser(user_uuid string) (string, error) {
	userUUID, errUser := GetUsernameByUserAppRoleUUID(user_uuid)
	if errUser != nil {
		log.Print(errUser)
		return "", errUser
	}
	var username string
	err := db.QueryRow("SELECT user_name from user_ms WHERE user_uuid = $1", userUUID).Scan(&username)
	if err != nil {
		return "", err
	}
	return username, nil
}

func UpdateUserAppRole(updateUserAppRole models.UpdateUser, userApplicationRoleUUID string) (models.UpdateUser, error) {
	userUUID, err := GetUsernameByUserAppRoleUUID(userApplicationRoleUUID)
	if err != nil {
		return models.UpdateUser{}, err
	}

	var username string

	username, errUser := GetUsernameByIDUser(userUUID)
	if errUser != nil {
		if errUser == sql.ErrNoRows {
			log.Println("No rows found for user_uuid:", userUUID)
			return models.UpdateUser{}, errUser
		}
	}

	log.Println("user_application_role_uuid:", userApplicationRoleUUID)

	currentTime := time.Now()

	_, errInsert := db.NamedExec("UPDATE user_ms SET user_name = :user_name, user_email = :user_email, updated_by = :updated_by, updated_at = :updated_at  WHERE user_id = (SELECT user_id FROM user_application_role_ms WHERE user_application_role_uuid = :user_application_role_uuid)", map[string]interface{}{
		"user_name":                  updateUserAppRole.Username,
		"user_email":                 updateUserAppRole.Email,
		"updated_by":                 username,
		"updated_at":                 currentTime,
		"user_application_role_uuid": userApplicationRoleUUID,
	})

	if errInsert != nil {
		return models.UpdateUser{}, errInsert
	}

	var user_id int64
	err = db.Get(&user_id, "SELECT user_id FROM user_ms WHERE user_id = (SELECT user_id FROM user_application_role_ms WHERE user_application_role_uuid = $1)", userApplicationRoleUUID)
	if err != nil {
		return models.UpdateUser{}, err
	}

	// Mendapatkan role_id yang baru saja diinsert
	var roleID int64
	err = db.Get(&roleID, "SELECT role_id FROM role_ms WHERE role_uuid = $1 AND deleted_at IS NULL", updateUserAppRole.ApplicationRole.Role_UUID)
	if err != nil {
		log.Println("Error getting role_id:", err)
		return models.UpdateUser{}, err
	}
	var applicationID int64
	err = db.Get(&applicationID, "SELECT application_id FROM application_ms WHERE application_uuid = $1 AND deleted_at IS NULL", updateUserAppRole.ApplicationRole.Application_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return models.UpdateUser{}, err
	}

	// Get division_id
	var divisionID int64
	err = db.Get(&divisionID, "SELECT division_id FROM division_ms WHERE division_uuid = $1 AND deleted_at IS NULL", updateUserAppRole.ApplicationRole.Division_UUID)
	if err != nil {
		log.Println("Error fetching division_id:", err)
		return models.UpdateUser{}, err
	}

	var applicationRoleId string

	err = db.Get(&applicationRoleId, "SELECT application_role_id FROM user_application_role_ms WHERE user_application_role_uuid = $1", userApplicationRoleUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No rows found for user_application_role_uuid:", userApplicationRoleUUID)
			return models.UpdateUser{}, fmt.Errorf("no rows found for user_application_role_uuid: %s", userApplicationRoleUUID)
		}
		log.Println("Error getting application_role_id by user_application_role_uuid:", err)
		return models.UpdateUser{}, err
	}

	log.Println("Application Role ID:", applicationRoleId)

	// Update data di application_role_ms
	_, err = db.Exec("UPDATE application_role_ms SET application_id = $1, role_id = $2, created_by = $3 WHERE application_role_id = $4",
		applicationID, roleID, username, applicationRoleId)
	if err != nil {
		log.Println("Error updating data in application_role_ms:", err)
		return models.UpdateUser{}, err
	}
	log.Println("Data updated in application_role_ms successfully")

	// Get application_role_id
	var applicationRoleID int64
	err = db.Get(&applicationRoleID, "SELECT application_role_id FROM application_role_ms WHERE application_id = $1 AND role_id = $2",
		applicationID, roleID)
	if err != nil {
		log.Println("Error getting application_role_id:", err)
		return models.UpdateUser{}, err
	}
	log.Println("Application Role ID:", applicationRoleID)

	// Update data di user_application_role_ms
	_, err = db.Exec("UPDATE user_application_role_ms SET application_role_id = $1, division_id = $2, updated_by = $3, updated_at = $5 WHERE user_id = $4",
		applicationRoleID, divisionID, username, user_id, time.Now())
	if err != nil {
		log.Println("Error updating data in user_application_role_ms:", err)
		return models.UpdateUser{}, err
	}

	return models.UpdateUser{}, nil
}

func DeleteUserAppRole(id, userUUID string) error {
	username, errUser := GetUsernameByID(userUUID)
	if errUser != nil {
		log.Print(errUser)
		return errUser
	}

	currentTime := time.Now()
	result, err := db.NamedExec("UPDATE user_application_role_ms SET deleted_by = :deleted_by, deleted_at = :deleted_at WHERE user_application_role_uuid = :id AND deleted_at IS NULL", map[string]interface{}{
		"deleted_by": username,
		"deleted_at": currentTime,
		"id":         id,
	})
	if err != nil {
		log.Print(err)
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
