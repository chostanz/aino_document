package service

import (
	"aino_document/models"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func AddRole(addRole models.Role, userID int) error {
	username, errP := GetUsernameByID(userID)
	if errP != nil {
		return errP
	}
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
		"created_by": username,
	})

	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

func GetAllRole() ([]models.Roles, error) {
	allRole := []models.Roles{}

	rows, err := db.Queryx("SELECT role_order, role_code, role_title, role_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM role_ms")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		place := models.Roles{}
		err := rows.StructScan(&place)
		if err != nil {
			return nil, err
		}
		allRole = append(allRole, place)

	}
	return allRole, nil
}

func ShowRoleById(id int) (models.Roles, error) {
	var roleId models.Roles
	idStr := strconv.Itoa(id)

	err := db.Get(&roleId, "SELECT role_order, role_code, role_title, role_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM role_ms WHERE role_order = $1", idStr)
	if err != nil {
		return models.Roles{}, err
	}
	return roleId, nil
}

func GetRoleById(id int) (models.Role, error) {
	var roleId models.Role
	idStr := strconv.Itoa(id)

	err := db.Get(&roleId, "SELECT * FROM role_ms WHERE role_order = $1", idStr)
	if err != nil {
		return models.Role{}, err
	}
	return roleId, nil

}
