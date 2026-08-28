package routes

import (
	"Productbid/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func RegisterProductRoutes(app *fiber.App, db *gorm.DB) {
	productHandler := handlers.NewProductHandler(db)

	app.Post("/api/products/resolve", productHandler.ResolveProduct)
}
