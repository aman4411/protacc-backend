package app

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, f *Factories) {
	// Setup CORS
	app.Use(cors.New())

	// API routes group
	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/signup", f.UserHandler.Signup)
	auth.Post("/verify-email", f.UserHandler.VerifyEmail)

	// Health check route
	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
}
