package models

type Users struct {
	User_uuid         string `json:"user_uuid" db:"user_uuid"`
	User_name         string `json:"user_name" db:"user_name"`
	Role_title        string `json:"role_title" db:"role_title"`
	Application_title string `json:"application_title" db:"application_title"`
	Division_title    string `json:"division_title" db:"division_title"`
}
