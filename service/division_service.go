package service

import (
	"aino_document/database"
	"aino_document/models"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB = database.Connection()

func GetAllDivision() ([]models.Divisions, error) {
	divisi := []models.Divisions{}

	rows, errSelect := db.Queryx("select division_order, division_code, division_title, division_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at from division_ms")
	if errSelect != nil {
		return nil, errSelect
	}

	for rows.Next() {
		place := models.Divisions{}
		rows.StructScan(&place)
		divisi = append(divisi, place)
	}

	return divisi, nil
}

func ShowDivisionById(id int) (models.Divisions, error) {
	var divisiId models.Divisions
	idStr := strconv.Itoa(id)

	err := db.Get(&divisiId, "SELECT division_order, division_code, division_title, division_show, created_by, created_at, updated_by, updated_at, deleted_by, deleted_at from division_ms WHERE division_order = $1", idStr)
	if err != nil {
		return models.Divisions{}, err
	}
	return divisiId, nil

}

func GetDivisionById(id int) (models.Division, error) {
	var divisiId models.Division
	idStr := strconv.Itoa(id)

	err := db.Get(&divisiId, "SELECT * FROM division_ms WHERE division_order = $1", idStr)
	if err != nil {
		return models.Division{}, err
	}
	return divisiId, nil

}

func AddDivision(addDivision models.Division, userUUID string) error {
	username, errP := GetUsernameByID(userUUID)
	if errP != nil {
		return errP
	}
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	divisionid := currentTimestamp + int64(uniqueID)

	uuid := uuid.New()
	uuidString := uuid.String()
	_, err := db.NamedExec("INSERT INTO division_ms (division_id, division_uuid, division_code, division_title, created_by) VALUES (:division_id, :division_uuid, :division_code, :division_title, :created_by)", map[string]interface{}{
		"division_id":    divisionid,
		"division_uuid":  uuidString,
		"division_code":  addDivision.Code,
		"division_title": addDivision.Title,
		"created_by":     username,
	})
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
