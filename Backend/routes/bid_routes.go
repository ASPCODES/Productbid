package routes

import(
	"Productbid/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func RegisterBidRoutes(app *fiber.App, db *gorm.DB) {
	bidHandler := handlers.NewBidHandler(db)

	app.Post("/api/bids/preview", bidHandler.PreviewRank)
	app.Post("/api/bids/place", bidHandler.PlaceBid)
}

