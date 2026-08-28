package middlewares

import(
	"github.com/gofiber/fiber/v2/middleware/cors"
)


func SetupCORS() cors.Config {
	return cors.Config{
		AllowOrigins: "*",  // productbid.space
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowCredentials: false,
	}
}

