package service

import (
	"aino_document/models"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("role not found")

func AddRole(addRole models.Role, userUUID string) error {
	username, errP := GetUsernameByID(userUUID)
	if errP != nil {
		return errP
	}

	var existingRoleID int
	err := db.QueryRow("SELECT role_id FROM role_ms WHERE (role_title = $1 OR role_code = $2) AND deleted_at IS NOT NULL", addRole.Title, addRole.Code).Scan(&existingRoleID)

	// Jika data ditemukan
	if err == nil {
		// Duplikat ditemukan dan sudah dihapus lembut, kembalikan data yang dihapus lembut berdasarkan division_id
		err = RestoreSoftDeletedRole(existingRoleID, userUUID, addRole)
		if err != nil {
			log.Printf("error restore : %s", err)
			return err
		}
	}
	// Pengecekan duplikat berdasarkan role_title
	// var existingRoleTitle string
	// err := db.QueryRow("SELECT role_title FROM role_ms WHERE role_title = $1 AND deleted_at IS NOT NULL", addRole.Title).Scan(&existingRoleTitle)
	// if err == nil {
	// 	// Duplikat ditemukan dan sudah dihapus lembut, kembalikan data yang dihapus lembut
	// 	err = RestoreSoftDeletedRole(addRole.Title, userUUID)
	// 	if err != nil {
	// 		log.Print(err)
	// 		return err
	// 	}
	// }

	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	roleid := currentTimestamp + int64(uniqueID)

	uuid := uuid.New()
	uuidString := uuid.String()

	_, err = db.NamedExec("INSERT INTO role_ms(role_id, role_uuid, role_code, role_title, created_by)VALUES(:role_id, :role_uuid, :role_code, :role_title, :created_by)", map[string]interface{}{
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

	rows, err := db.Queryx("SELECT role_uuid, role_order, role_code, role_title, role_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM role_ms WHERE deleted_at IS NULL")
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

func ShowRoleById(id string) (models.Roles, error) {
	var roleUUID models.Roles

	err := db.Get(&roleUUID, "SELECT  role_uuid, role_order, role_code, role_title, role_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at FROM role_ms WHERE role_uuid = $1 AND deleted_at IS NULL", id)
	if err != nil {
		return models.Roles{}, err
	}
	return roleUUID, nil
}

func GetRoleById(id string) (models.Role, error) {
	var roleId models.Role

	err := db.Get(&roleId, "SELECT * FROM role_ms WHERE role_order = $1", id)
	if err != nil {
		return models.Role{}, err
	}
	return roleId, nil

}

func UpdateRole(updateRole models.Role, id string, userUUID string) (models.Role, error) {
	username, errUser := GetUsernameByID(userUUID)
	if errUser != nil {
		log.Print(errUser)
		return models.Role{}, errUser

	}

	currentTime := time.Now()

	_, err := db.NamedExec("UPDATE role_ms SET role_code = :role_code, role_title = :role_title, updated_by = :updated_by, updated_at = :updated_at WHERE role_uuid = :id", map[string]interface{}{
		"role_code":  updateRole.Code,
		"role_title": updateRole.Title,
		"updated_by": username,
		"updated_at": currentTime,
		"id":         id,
	})

	if err != nil {
		log.Print(err)
		return models.Role{}, err
	}
	return updateRole, nil
}

func RestoreSoftDeletedRole(roleID int, userUUID string, addRole models.Role) error {
	username, _ := GetUsernameByID(userUUID)
	// Cari role yang dihapus lembut dengan role_title tertentu
	log.Printf("Restoring role with ID: %d", roleID)
	// Lakukan UPDATE untuk mengembalikan division yang dihapus lembut
	_, err := db.Exec("UPDATE role_ms SET created_at = NOW(), created_by = $2, updated_at = NULL, updated_by = '',  deleted_at = NULL, deleted_by = '', role_code = $3, role_title = $4 WHERE role_id = $1", roleID, username, addRole.Code, addRole.Title)
	if err != nil {
		log.Printf("Error during role restore: %s", err)
		return err
	}

	return nil
}

func DeleteRole(id string, userUUID string) error {
	username, errUser := GetUsernameByID(userUUID)
	if errUser != nil {
		log.Print(errUser)
		return errUser
	}

	currentTime := time.Now()
	result, err := db.NamedExec("UPDATE role_ms SET deleted_by = :deleted_by, deleted_at = :deleted_at WHERE role_uuid = :id AND deleted_at IS NULL", map[string]interface{}{
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
		return ErrNotFound // Mengembalikan error jika tidak ada rekaman yang cocok
	}

	return nil
}
