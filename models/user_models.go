package models

import (
	"database/sql"
	"time"
)

type Users struct {
	User_uuid         string         `json:"user_uuid" db:"user_uuid"`
	User_App_uuid     string         `json:"user_application_role_uuid" db:"user_application_role_uuid"`
	User_name         string         `json:"user_name" db:"user_name"`
	User_email        string         `json:"user_email" db:"user_email"`
	Role_title        string         `json:"role_title" db:"role_title"`
	Application_title string         `json:"application_title" db:"application_title"`
	Division_title    string         `json:"division_title" db:"division_title"`
	Created_by        sql.NullString `json:"created_by" db:"created_by"`
	Created_at        time.Time      `json:"created_at" db:"created_at"`
	Updated_by        sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at        sql.NullTime   `json:"updated_at" db:"updated_at"`
	Deleted_by        sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at        sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}
