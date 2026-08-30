package routes

import (
	"Productbid/handlers"
	"Productbid/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterBidRoutes(
	app *fiber.App,
	db *gorm.DB,
	bidService *services.BidService,
	paymentService *services.PaymentService,
	frontendURL string,
) {
	bidHandler := handlers.NewBidHandler(db, bidService, paymentService, frontendURL)

	app.Post("/api/bids/preview", bidHandler.PreviewRank)
	app.Post("/api/bids/initiate", bidHandler.InitiateBid)
}

