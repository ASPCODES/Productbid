package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")
	return db
}

