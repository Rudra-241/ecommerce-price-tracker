package db

import (
	"ecommerce-price-tracker/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func InitWithDSN(dsn string) *gorm.DB {

	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Cannot open DB")
	}
	db.AutoMigrate(&models.User{}, &models.Product{})

	return db

}

func GetDB() *gorm.DB {
	if db == nil {
		panic("DB not initialized yet")
	}
	return db
}
