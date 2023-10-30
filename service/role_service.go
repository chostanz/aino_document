package service

import (
	"aino_document/models"
	"log"
	"time"

	"github.com/google/uuid"
)

func AddRole(addRole models.Role) error {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	roleid := currentTimestamp + int64(uniqueID)

	uuid := uuid.New()
	uuidString := uuid.String()

	_, err := db.NamedExec("INSERT INTO role_ms(role_id, role_uuid, role_code, role_title, created_by)VALUES(:role_id, :role_uuid, :role_code, :role_title, :created_by)", map[string]interface{}{
		"role_id":    roleid,
		"role_uuid":  uuidString,
		"role_code":  addRole.Code,
		"role_title": addRole.Title,
		"created_by": addRole.Created_by,
	})

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
