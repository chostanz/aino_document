package models

type Profile struct {
	UUID     string `json:"user_uuid" db:"user_uuid"`
	Username string `json:"user_name" db:"user_name"`
	RoleCode string `json:"role_code" db:"role_code"`
}
