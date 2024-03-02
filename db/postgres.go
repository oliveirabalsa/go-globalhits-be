package db

import (
	"github.com/oliveirabalsa/go-globalhitss-be/app/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewPostgresDB() (*gorm.DB, error) {
	dsn := "host=localhost user=globalhitss password=globalhitss dbname=globalhitss port=5432 sslmode=disable TimeZone=UTC"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&model.Client{})
	return db, nil
}
