package routes

import(
	"Productbid/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func RegisterCategoryRoutes(app *fiber.App, db *gorm.DB) {
	categoryHandler := handlers.NewCategoryHandler(db)

	app.Get("/api/categories", categoryHandler.GetCategories)
}

