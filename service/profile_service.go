package service

import (
	"aino_document/models"
	"log"
)

func MyProfile(uuid string) (models.Profile, error) {
	var myProfile models.Profile

	err := db.Get(&myProfile, `
	SELECT u.user_uuid, u.user_name, r.role_code 
	FROM user_ms u 
	INNER JOIN user_application_role_ms uar ON u.user_id = uar.user_id 
	INNER JOIN application_role_ms ar ON uar.application_role_id = ar.application_role_id
	INNER JOIN role_ms r ON ar.role_id = r.role_id
	WHERE u.deleted_at IS NULL AND u.user_uuid = $1
`, uuid)

	if err != nil {
		log.Print(err)
		return models.Profile{}, err
	}

	return myProfile, nil

}
