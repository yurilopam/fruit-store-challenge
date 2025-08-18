package database

import (
	"log"

	"github.com/yurilopam/fruit-store-challenge/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Connect() {
	dsn := "host=localhost user=fruit_store_admin_user password=password dbname=fruit_store port=5432 sslmode=disable"
	DB, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Panic(err)
	}
	DB.AutoMigrate(&models.Fruit{}, &models.User{})
}
