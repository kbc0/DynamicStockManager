package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	jwtware "github.com/gofiber/jwt/v2"
)

func RegisterMiddleware(app *fiber.App) {
	// Apply CORS settings for all routes, adjust as per your requirements
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	// JWT Middleware for protected routes
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"), 
		Filter: func(ctx *fiber.Ctx) bool {
			// List of routes that don't require authentication
			unprotectedPaths := []string{
				"/api/v1/login",
				"/api/v1/register",
			}

			// Skip JWT middleware for the above routes
			for _, path := range unprotectedPaths {
				if ctx.Path() == path {
					return true
				}
			}

			// By default require JWT authentication
			return false
		},
	}))
}
