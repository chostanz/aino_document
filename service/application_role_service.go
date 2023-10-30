package service

import (
	"aino_document/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddApplicationRole(addAppRole models.ApplicationRole) error {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	app_role_id := currentTimestamp + int64(uniqueID)

	_, err := db.NamedExec("INSERT INTO application_role_ms(application_role_id, application_id, role_id, created_by) VALUES(:application_role_id, :application_id, :role_id, :created_by)", map[string]interface{}{
		"application_role_id": app_role_id,
		"application_id":      addAppRole.Application_id,
		"role_id":             addAppRole.Role_id,
		"created_by":          addAppRole.Created_by,
	})

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
