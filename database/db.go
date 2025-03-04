package database

import (
	"log"

	"github.com/ahnafasif/MarauderBoard/configs"
	"github.com/ahnafasif/MarauderBoard/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	configs.LoadEnv()
	dsn := configs.DbUrl

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	DB = db

	err = DB.AutoMigrate(models.GetModels()...)
	if err != nil {
		log.Fatal("Error migrating models: ", err)
	}

	log.Println("Connected to database")
}
