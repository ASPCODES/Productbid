package middlewares

import(
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func SetupCORS() cors.Config {
	return cors.Config{
		AllowOrigins: "https://productbid.space,http://localhost:3000",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization, X-Requested-With, webhook-id, webhook-signature, webhook-timestamp, Webhook-Id, Webhook-Signature, Webhook-Timestamp",
		AllowCredentials: false,
	}
}

