package repository

import (
	"github.com/google/uuid"
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
	if err := r.DB.Where("active = ?", true).Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *ClientRepository) GetByID(clientID uuid.UUID) (*model.Client, error) {
	var client model.Client
	if err := r.DB.Where("id = ? and active = true", clientID).First(&client).Error; err != nil {
		return nil, err
	}
	return &client, nil
}

func (r *ClientRepository) Create(client *model.Client) (*model.Client, error) {
	if err := r.DB.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (r *ClientRepository) Update(client *model.Client) (*model.Client, error) {
	if err := r.DB.Model(&client).Updates(client).Error; err != nil {
		return nil, err
	}

	return client, nil

}

func (r *ClientRepository) DeleteClient(clientId uuid.UUID) error {
	var client model.Client
	if err := r.DB.Where("id = ? AND active = ?", clientId, true).First(&client).Error; err != nil {
		return err
	}
	if err := r.DB.Model(&model.Client{}).Where("id = ?", clientId).Update("active", false).Error; err != nil {
		return err
	}
	return nil
}
