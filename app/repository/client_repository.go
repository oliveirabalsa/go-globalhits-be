package repository

import (
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"gorm.io/gorm"
)

type ClientRepository struct {
	DB *gorm.DB
}

func NewClientRepository(db *gorm.DB) *ClientRepository {
	return &ClientRepository{
		DB: db,
	}
}

func (r *ClientRepository) GetAll() ([]*model.Client, error) {
	var clients []*model.Client
	if err := r.DB.Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ClientRepository) Create(client *model.Client) (*model.Client, error) {
	if err := r.DB.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}
