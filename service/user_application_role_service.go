package service

import "aino_document/models"

func GetUserApplicationRole() ([]models.Users, error) {
	getUserAppRole := []models.Users{}

	rows, err := db.Queryx("SELECT u.user_uuid, u.user_name, r.role_title, a.application_title, d.division_title FROM user_ms u INNER JOIN user_application_role_ms uar ON u.user_id = uar.user_id INNER JOIN application_role_ms ar ON uar.application_role_id = ar.application_role_id INNER JOIN application_ms a ON ar.application_id = a.application_id INNER JOIN role_ms r ON ar.role_id = r.role_id INNER JOIN division_ms d ON uar.division_id = d.division_id")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		place := models.Users{}
		rows.StructScan(&place)
		getUserAppRole = append(getUserAppRole, place)
	}
	return getUserAppRole, nil
}
