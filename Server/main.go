package main

import (
	"Productbid/config"
	"Productbid/db"
	"Productbid/middlewares"
	"Productbid/routes"
	"Productbid/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func main() {
	// Load config from .env
	cfg := config.LoadConfig()

	// Connect to database
	database := db.ConnectDB(cfg.DatabaseURL)

	// Initialize Services
	paymentService := services.NewPaymentService(
		cfg.DodoAPIKey,
		cfg.DodoWebhookSecret,
		cfg.DodoMode,
	)
	bidService := services.NewBidService(database)

	// Create a fiber app
	app := fiber.New()

	// Adding middlewares
	app.Use(cors.New(middlewares.SetupCORS()))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("ProductBid backend is running")
	})


	// Register all Routes
	routes.RegisterCategoryRoutes(app, database)
	routes.RegisterProductRoutes(app, database)
	routes.RegisterBidRoutes(app, database, bidService, paymentService, cfg.FrontendURL, cfg.DodoProductID)
	routes.RegisterLeaderboardRoutes(app, database)
	routes.RegisterWebhookRoutes(app, database, paymentService, bidService)

	log.Println("ProductBid backend listening on port :" + cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}
