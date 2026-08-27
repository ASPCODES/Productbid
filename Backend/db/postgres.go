package db

import (
	"log"
	"Productbid/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func ConnectDB(dsn string) *gorm.DB {
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Database connected successfully")

	// Auto-migrate all models — creates tables if they don't exist
	err = database.AutoMigrate(
		&models.Category{},
		&models.Product{},
		&models.Bid{},
		models.Payment{},
	)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Migration completed successfully!!")

	return database
}

