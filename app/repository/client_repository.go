package repository

import (
	"math"

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

func (r *ClientRepository) GetAll(page int, pageSize int) ([]*model.Client, int, error) {
	var clients []*model.Client
	var totalRecords int64

	err := WithTransaction(r.DB, func(tx *gorm.DB) error {
		if err := tx.Model(&model.Client{}).Where("active = ?", true).Count(&totalRecords).Error; err != nil {
			return err
		}

		if err := tx.Where("active = ?", true).Offset((page - 1) * pageSize).Limit(pageSize).Find(&clients).Error; err != nil {
			return err
		}

		return nil
	})

	totalPages := int(math.Ceil(float64(totalRecords) / float64(pageSize)))

	return clients, totalPages, err
}

func (r *ClientRepository) GetByID(clientID uuid.UUID) (*model.Client, error) {
	var client model.Client
	err := WithTransaction(r.DB, func(tx *gorm.DB) error {
		return tx.Where("id = ? and active = true", clientID).First(&client).Error
	})

	return &client, err
}

func (r *ClientRepository) Create(client *model.Client) (*model.Client, error) {
	err := WithTransaction(r.DB, func(tx *gorm.DB) error {
		return tx.Create(client).Error
	})

	return client, err
}

func (r *ClientRepository) Update(client *model.Client) (*model.Client, error) {
	err := WithTransaction(r.DB, func(tx *gorm.DB) error {
		return tx.Model(&client).Updates(client).Error
	})

	return client, err
}

func (r *ClientRepository) DeleteClient(clientId uuid.UUID) error {
	err := WithTransaction(r.DB, func(tx *gorm.DB) error {
		var client model.Client
		if err := tx.Where("id = ? AND active = ?", clientId, true).First(&client).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Client{}).Where("id = ?", clientId).Update("active", false).Error; err != nil {
			return err
		}
		return nil
	})

	return err
}

func WithTransaction(db *gorm.DB, fn func(tx *gorm.DB) error) error {
	tx := db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
