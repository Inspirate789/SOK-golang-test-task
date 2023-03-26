package models

type User struct {
	ID               int    `json:"id" db:"id"`
	Name             string `json:"name" db:"name"`
	RegistrationDate int    `json:"registration_date" db:"registration_date"`
	Balance          int    `json:"balance" db:"balance"`
}
