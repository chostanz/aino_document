package service

import (
	"aino_document/models"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func AddApplicationRole(addAppRole models.ApplicationRole, userUUID string) error {
	username, errP := GetUsernameByID(userUUID)
	if errP != nil {
		return errP
	}
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	app_role_id := currentTimestamp + int64(uniqueID)

	_, err := db.NamedExec("INSERT INTO application_role_ms(application_role_id, application_id, role_id, created_by) VALUES(:application_role_id, :application_id, :role_id, :created_by)", map[string]interface{}{
		"application_role_id": app_role_id,
		"application_id":      addAppRole.Application_id,
		"role_id":             addAppRole.Role_id,
		"created_by":          username,
	})

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func GetAppRole(id int) (models.ApplicationRole, error) {
	var appRoleId models.ApplicationRole
	idStr := strconv.Itoa(id)

	err := db.Get(&appRoleId, "SELECT application_role_id, application_id, role_id,created_by, created_at, updated_by, updated_at, deleted_by, deleted_at from application_role_ms WHERE application_role_id = $1", idStr)
	if err != nil {
		return models.ApplicationRole{}, err
	}
	return appRoleId, nil

}
