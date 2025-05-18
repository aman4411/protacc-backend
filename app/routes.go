package app

import (
	"github.com/aman4411/protacc-backend/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// SetupRoutes configures all the routes for the application
func SetupRoutes(app *fiber.App, f *Factories) {
	// Setup CORS with specific configuration
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "https://protacc.netlify.app, http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders:    "Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
		AllowCredentials: true,
	}))

	// API routes group
	api := app.Group("/api/v1")

	// Public routes
	auth := api.Group("/auth")
	auth.Post("/signup", f.UserHandler.Signup)
	auth.Post("/login", f.UserHandler.Login)
	auth.Post("/verify-email", f.UserHandler.VerifyEmail)

	// Protected routes
	user := api.Group("/user", middleware.Protected())
	user.Get("/profile", f.UserHandler.GetProfile)

	// Admin routes
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("admin"))
	admin.Get("/users", f.UserHandler.GetUsers)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})
}
