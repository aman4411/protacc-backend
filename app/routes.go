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

	// Service routes
	services := api.Group("/services")
	services.Get("/", f.ServiceHandler.GetServices)
	services.Get("/categories", f.ServiceHandler.GetServiceCategories)
	services.Get("/search", f.ServiceHandler.SearchServices)
	services.Get("/slug/:slug", f.ServiceHandler.GetServiceBySlug)
	services.Get("/:id", f.ServiceHandler.GetServiceByID)

	// Protected service routes
	protectedServices := api.Group("/services", middleware.Protected())
	protectedServices.Post("/categories", f.ServiceHandler.CreateServiceCategory)
	protectedServices.Post("/", f.ServiceHandler.CreateService)

	// Cart routes
	cart := api.Group("/cart", middleware.Protected())
	cart.Get("/", f.ServiceHandler.GetCartItems)
	cart.Post("/:serviceId", f.ServiceHandler.AddToCart)
	cart.Delete("/:serviceId", f.ServiceHandler.RemoveFromCart)

	// Order routes
	orders := api.Group("/orders", middleware.Protected())
	orders.Get("/", f.ServiceHandler.GetOrders)
	orders.Post("/", f.ServiceHandler.CreateOrderFromCart)            // New route for cart-based orders
	orders.Post("/services/:serviceId", f.ServiceHandler.CreateOrder) // Single service orders
	orders.Get("/:orderId/history", f.ServiceHandler.GetOrderStatusHistory)

	// Admin routes
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("admin"))
	admin.Get("/users", f.UserHandler.GetUsers)
	admin.Put("/users/:userId/role", f.UserHandler.UpdateUserRole)
	admin.Get("/orders", f.ServiceHandler.GetOrders)
	admin.Put("/orders/:orderId/status", f.ServiceHandler.UpdateOrderStatus)

	// Service management
	admin.Get("/services", f.ServiceHandler.GetServices)
	admin.Post("/services", f.ServiceHandler.CreateService)
	admin.Put("/services/:id", f.ServiceHandler.UpdateService)
	admin.Delete("/services/:id", f.ServiceHandler.DeleteService)

	// Category management
	admin.Get("/categories", f.ServiceHandler.GetServiceCategories)
	admin.Post("/categories", f.ServiceHandler.CreateServiceCategory)
	admin.Put("/categories/:id", f.ServiceHandler.UpdateServiceCategory)
	admin.Delete("/categories/:id", f.ServiceHandler.DeleteServiceCategory)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Static files
	app.Static("/images", "./static/images")
}
