package routes

import (
	"Productbid/handlers"
	"Productbid/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func RegisterWebhookRoutes(
	app *fiber.App,
	db *gorm.DB,
	paymentService *services.PaymentService,
	bidService *services.BidService,
) {
	webhookHandler := handlers.NewWebhookHandler(db, paymentService, bidService)

	// Server-to-server webhook endpoint called by Dodo Payments
	app.Post("/api/webhooks/dodo", webhookHandler.HandleDodoWebhook)
}
