package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Login struct {
	Username         string `json:"user_name" db:"user_name" validate:"required"`
	Password         string `json:"user_password" db:"user_password" validate:"required"`
	Application_Role int    `json:"application_role_id" db:"application_role_id"`
	Division         int    `json:"division_id" db:"division_id"`
}

type Register struct {
	Id         int            `json:"user_id" db:"user_id"`
	UUID       uuid.UUID      `json:"user_uuid" db:"user_uuid"`
	Username   string         `json:"user_name" db:"user_name" validate:"required"`
	Password   string         `json:"user_password" db:"user_password" validate:"required"`
	Email      string         `json:"user_email" db:"user_email" validate:"required,email"`
	Created_by string         `json:"created_by" db:"created_by"`
	Created_at time.Time      `json:"created_at" db:"created_at"`
	Updated_by sql.NullString `json:"updated_by" db:"updated_by"`
	Updated_at sql.NullTime   `json:"updated_at" db:"updated_by"`
	Deleted_by sql.NullString `json:"deleted_by" db:"deleted_by"`
	Deleted_at sql.NullTime   `json:"deleted_at" db:"deleted_at"`
}

type UserAppRole struct {
	Application_role_id int `json:"application_role_id" db:"application_role_id"`
	Division_id         int `json:"division_id" db:"division_id"`
}
