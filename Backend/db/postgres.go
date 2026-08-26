package db

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func NewConnection(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Database isn't Connected!!", err)
	}

	log.Println("Database Connected!!")
	return db
}

