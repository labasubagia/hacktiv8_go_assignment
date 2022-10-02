package database

import (
	"assignment2/models"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	Host     = "localhost"
	Port     = "5432"
	User     = "root"
	Password = "root"
	Database = "orders_by"
)

func NewPostgres() (*gorm.DB, error) {
	config := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", Host, User, Password, Database, Port)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Debug().AutoMigrate(models.Order{}, models.Item{})
	return db, nil
}
