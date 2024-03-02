package model

import "github.com/google/uuid"

type Client struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name"`
	LastName    string    `json:"last_name"`
	Contact     string    `json:"contact"`
	Address     string    `json:"address"`
	DateOfBirth string    `json:"date_of_birth"`
	CPF         string    `json:"cpf"`
}
