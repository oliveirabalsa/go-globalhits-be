package model

import (
	"fmt"
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/oliveirabalsa/go-globalhitss-be/app/utils"
)

type Client struct {
	ID          uuid.UUID `gorm:"primaryKey" json:"id" swaggerignore:"true"`
	Name        string    `json:"name" validate:"required" example:"Fulano"`
	LastName    string    `json:"last_name" validate:"required" example:"da Silva"`
	Contact     string    `json:"contact" example:"fulano.dasilva@example.com"`
	Address     string    `json:"address" example:"Rua x, bairro da rua, cidade/estado"`
	DateOfBirth string    `json:"date_of_birth" validate:"required,date_of_birth" example:"12/12/1912"`
	CPF         string    `json:"cpf" validate:"required,cpf" example:"123.456.789-00"`
	Active      bool      `gorm:"default:true" json:"active" swaggerignore:"true"`
	CreatedAt   time.Time `gorm:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"updated_at" json:"updatedAt"`
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

func (c *Client) EncryptSensitiveData() {
	if c.Name != "" {
		c.Name = encryption.Encrypt(c.Name)
	}
	if c.LastName != "" {
		c.LastName = encryption.Encrypt(c.LastName)
	}
	if c.Contact != "" {
		c.Contact = encryption.Encrypt(c.Contact)
	}
	if c.Address != "" {
		c.Address = encryption.Encrypt(c.Address)
	}
	if c.DateOfBirth != "" {
		c.DateOfBirth = encryption.Encrypt(c.DateOfBirth)
	}
	if c.CPF != "" {
		c.CPF = encryption.Encrypt(c.CPF)
	}
}

func (c *Client) DecryptSensitiveData() {
	if c.Name != "" {
		c.Name = encryption.Decrypt(c.Name)
	}
	if c.LastName != "" {
		c.LastName = encryption.Decrypt(c.LastName)
	}
	if c.Contact != "" {
		c.Contact = encryption.Decrypt(c.Contact)
	}
	if c.Address != "" {
		c.Address = encryption.Decrypt(c.Address)
	}
	if c.DateOfBirth != "" {
		c.DateOfBirth = encryption.Decrypt(c.DateOfBirth)
	}
	if c.CPF != "" {
		c.CPF = encryption.Decrypt(c.CPF)
	}
}
