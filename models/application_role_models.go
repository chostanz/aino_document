package models

import (
	"database/sql"
	"time"
)

type ApplicationRole struct {
	Application_role_id int            `json:"application_role_id" db:"application_role_id"`
	Application_id      int            `json:"application_id" db:"application_id"`
	Role_id             int            `json:"role_id" db:"role_Id"`
	Created_by          string         `json:"created_by" db:"created_by"`
	Created_at          time.Time      `json:"created_at" db:"created_at"`
	Updated_by          sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at          sql.NullTime   `json:"updated_at" db:"updated_by"`
	Deleted_by          sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at          sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}
