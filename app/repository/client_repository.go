package repository

import (
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"gorm.io/gorm"
)

type ClientRepository struct {
	DB *gorm.DB
}

func (r *ClientRepository) Create(client *model.Client) (*model.Client, error) {
	err := r.DB.Create(client).Error
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Outros métodos para Read, Update e Delete
