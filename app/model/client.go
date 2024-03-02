package model

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type Client struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id"`
	Name        string    `json:"name" validate:"required"`
	LastName    string    `json:"last_name" validate:"required"`
	Contact     string    `json:"contact"`
	Address     string    `json:"address"`
	DateOfBirth string    `json:"date_of_birth" validate:"required,date_of_birth"`
	CPF         string    `json:"cpf" validate:"required,cpf"`
	Active      bool      `gorm:"default:true" json:"active"`
}

func (Client) TableName() string {
	return "clients"
}

func (c *Client) Validate() error {
	validate := validator.New()

	validate.RegisterValidation("cpf", func(fl validator.FieldLevel) bool {
		cpf := fl.Field().String()
		regex := regexp.MustCompile(`^\d{3}\.\d{3}\.\d{3}-\d{2}$`)
		return regex.MatchString(cpf)
	})

	validate.RegisterValidation("date_of_birth", func(fl validator.FieldLevel) bool {
		dateStr := fl.Field().String()
		_, err := time.Parse("02/01/2006", dateStr)
		return err == nil
	})

	if err := validate.Struct(c); err != nil {
		var validationErrors string
		errors := err.(validator.ValidationErrors)
		for i, e := range errors {
			if i > 0 {
				validationErrors += ", "
			}
			validationErrors += fmt.Sprintf("%s", e.Tag())
		}

		errorMessage := fmt.Sprintf("Validation error(s): %s", validationErrors)

		return fmt.Errorf(errorMessage)
	}

	return nil
}
