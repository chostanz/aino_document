package service

import (
	"aino_document/database"
	"aino_document/models"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB = database.Connection()

// func GetAllDivision([]models.Division, error) {
// 	divisi := []models.Division{}

// 	rows, errSelect := db.Queryx("select * from tb_division")
// 	if errSelect != nil {

// 	}

// 	for rows.Next() {
// 		place := models.Division{}
// 		rows.StructScan(&place)
// 		divisi = append(divisi, place)
// 	}

// 	return
// }

func AddDivision(addDivision models.Division) error {
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
		"created_by":     addDivision.Created_by,
	})
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
