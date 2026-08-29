package routes

import (
	"Productbid/handlers"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


func RegisterLeaderboardRoutes(app *fiber.App, db *gorm.DB) {
	leaderboardHandler := handlers.NewLeaderboardHandler(db)

	app.Get("/api/leaderboard/all", leaderboardHandler.GetLeaderboard)
}

