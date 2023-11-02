package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Id         int            `json:"role_id" db:"role_id"`
	UUID       uuid.UUID      `json:"role_uuid" db:"role_uuid"`
	Code       string         `json:"role_code" db:"role_id" validate:"required"`
	Title      string         `json:"role_title" db:"role_title" validate:"required"`
	Order      int            `json:"role_order" db:"role_order"`
	Show       bool           `json:"role_show" db:"role_show"`
	Created_by string         `json:"created_by" db:"created_by"`
	Created_at time.Time      `json:"created_at" db:"created_at"`
	Updated_by sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at sql.NullTime   `json:"updated_at" db:"updated_by"`
	Deleted_by sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}
