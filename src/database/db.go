package database

import (
	"fmt"

	"github.com/NYARAS/go-ambassador/src/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	var err error

	dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		"ambassador",
		"ambassador",
		"pleasechangeme",
		"ambassador",
		"5432",
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Could not connect with the database!")
	}
}

func AutoMigrate() {
	DB.AutoMigrate(models.User{}, models.Product{})
}
