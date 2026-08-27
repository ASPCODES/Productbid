package main

import (
	"log"
	"Productbid/config"
	"Productbid/db"


	"github.com/gofiber/fiber/v2"
	
)


func main() {
	// Load config from .env
	cfg := config.LoadConfig()

	// Connect to database
	db := db.ConnectDB(cfg.DatabaseURL)

	_ = db

	// Create a fiber app
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ProductBid backend is running")
	})

	log.Fatal(app.Listen(":" + cfg.Port))
}
