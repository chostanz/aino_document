package service

import (
	"aino_document/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddApplicationRole(addAppRole models.AddApplicationRole, userUUID string) error {
	username, errP := GetUsernameByID(userUUID)
	if errP != nil {
		return errP
	}
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()
	uuid := uuid.New()
	uuidString := uuid.String()
	app_role_id := currentTimestamp + int64(uniqueID)

	var roleID int64
	err := db.Get(&roleID, "SELECT role_id FROM role_ms WHERE role_uuid = $1", addAppRole.Role_UUID)
	if err != nil {
		log.Println("Error getting role_id:", err)
		return err
	}
	var applicationID int64
	err = db.Get(&applicationID, "SELECT application_id FROM application_ms WHERE application_uuid = $1", addAppRole.Application_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return err
	}
	_, errInsert := db.NamedExec("INSERT INTO application_role_ms(application_role_id, application_role_uuid, application_id, role_id, created_by) VALUES(:application_role_id, :application_role_uuid, :application_id, :role_id, :created_by)", map[string]interface{}{
		"application_role_id":   app_role_id,
		"application_role_uuid": uuidString,
		"application_id":        applicationID,
		"role_id":               roleID,
		"created_by":            username,
	})

	if errInsert != nil {
		log.Print(errInsert)
		return errInsert
	}
	return nil
}

func GetAppRole(id string) (models.ApplicationRole, error) {
	var appRoleId models.ApplicationRole
	// idStr := strconv.Itoa(id)

	err := db.Get(&appRoleId, "SELECT application_role_id, application_id, role_id, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at from application_role_ms WHERE application_role_uuid = $1", id)
	if err != nil {
		return models.ApplicationRole{}, err
	}
	return appRoleId, nil

}

func UpdateAppRole(updateAppRole models.AddApplicationRole, id string, userUUID string) (models.AddApplicationRole, error) {
	username, errUser := GetUsernameByID(userUUID)
	if errUser != nil {
		log.Print(errUser)
		return models.AddApplicationRole{}, errUser

	}

	var roleID int64
	err := db.Get(&roleID, "SELECT role_id FROM role_ms WHERE role_uuid = $1", updateAppRole.Role_UUID)
	if err != nil {
		log.Println("Error getting role_id:", err)
		return models.AddApplicationRole{}, err
	}
	var applicationID int64
	err = db.Get(&applicationID, "SELECT application_id FROM application_ms WHERE application_uuid = $1", updateAppRole.Application_UUID)
	if err != nil {
		log.Println("Error getting application_id:", err)
		return models.AddApplicationRole{}, err
	}

	currentTime := time.Now()

	_, errUpdate := db.NamedExec("UPDATE application_role_ms SET application_id = :application_id, role_id = :role_id, updated_by = :updated_by, updated_at = :updated_at WHERE application_role_uuid = :id", map[string]interface{}{
		"application_id": applicationID,
		"role_id":        roleID,
		"updated_by":     username,
		"updated_at":     currentTime,
		"id":             id,
	})

	if errUpdate != nil {
		log.Print(errUpdate)
		return models.AddApplicationRole{}, err
	}

	return updateAppRole, nil
}
