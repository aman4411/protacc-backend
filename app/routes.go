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
		AllowOrigins:     "https://protacc.in, https://www.protacc.in, https://protacc.netlify.app, http://localhost:3000",
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
	auth.Post("/forgot-password", f.UserHandler.ForgotPassword)
	auth.Post("/reset-password", f.UserHandler.ResetPassword)
	auth.Get("/reset-password/validate", f.UserHandler.ValidateResetToken)

	// Protected routes
	user := api.Group("/user", middleware.Protected())
	user.Get("/profile", f.UserHandler.GetProfile)
	user.Put("/profile", f.UserHandler.UpdateProfile)
	user.Post("/change-password-request", f.UserHandler.RequestPasswordReset)

	// Service routes
	services := api.Group("/services")
	services.Get("/", f.ServiceHandler.GetServices)
	services.Get("/categories", f.ServiceHandler.GetServiceCategories)
	services.Get("/search", f.ServiceHandler.SearchServices)
	services.Get("/slug/:slug", f.ServiceHandler.GetServiceBySlug)
	services.Get("/:id/reviews", f.ReviewHandler.GetServiceReviews)
	services.Get("/:id", f.ServiceHandler.GetServiceByID)

	// Public reviews
	api.Get("/reviews/top", f.ReviewHandler.GetTopReviews)

	// Public: coupons flagged to show on the cart
	api.Get("/coupons/available", f.CouponHandler.ListAvailable)
	// Public: coupons flagged for the homepage campaign banner
	api.Get("/coupons/promotions", f.CouponHandler.ListHomepage)
	// Public: upcoming tax/compliance deadlines for the homepage
	api.Get("/deadlines/upcoming", f.DeadlineHandler.ListUpcoming)
	// Public: blog posts
	api.Get("/posts", f.PostHandler.ListPublished)
	api.Get("/posts/categories", f.PostHandler.ListCategories)
	api.Get("/posts/:slug/related", f.PostHandler.GetRelated)
	api.Get("/posts/:slug", f.PostHandler.GetBySlug)

	// Protected review actions
	reviews := api.Group("/reviews", middleware.Protected())
	reviews.Post("/", f.ReviewHandler.SubmitReview)
	reviews.Get("/eligibility", f.ReviewHandler.GetEligibility)

	// Protected service routes
	protectedServices := api.Group("/services", middleware.Protected())
	protectedServices.Post("/categories", f.ServiceHandler.CreateServiceCategory)
	protectedServices.Post("/", f.ServiceHandler.CreateService)

	// Cart routes
	cart := api.Group("/cart", middleware.Protected())
	cart.Get("/", f.CartHandler.GetCartItems)
	cart.Post("/:serviceId", f.CartHandler.AddToCart)
	cart.Delete("/:serviceId", f.CartHandler.RemoveFromCart)

	// Order routes
	orders := api.Group("/orders", middleware.Protected())
	orders.Get("/", f.OrderHandler.GetOrders)
	orders.Get("/number/:orderNumber", f.OrderHandler.GetOrderByNumber) // Get order by number
	orders.Post("/", f.OrderHandler.CreateOrderFromCart)                // New route for cart-based orders
	orders.Post("/coupon/preview", f.OrderHandler.PreviewCoupon)        // Validate a coupon against the cart
	orders.Post("/services/:serviceId", f.OrderHandler.CreateOrder)     // Single service orders
	orders.Get("/:orderId/history", f.OrderHandler.GetOrderStatusHistory)
	orders.Get("/:orderId/documents", f.OrderHandler.GetOrderDocuments)
	orders.Post("/:orderId/documents", f.OrderHandler.AddUserOrderDocument)

	// Payment routes
	payments := api.Group("/payments", middleware.Protected())
	payments.Post("/orders/:orderId/create", f.PaymentHandler.CreatePaymentOrder)
	payments.Post("/verify", f.PaymentHandler.VerifyPayment)
	payments.Get("/orders/:orderId/status", f.PaymentHandler.GetPaymentStatus)

	// Webhook route (no authentication required)
	api.Post("/payments/webhook", f.PaymentHandler.HandleWebhook)

	// Admin routes
	admin := api.Group("/admin", middleware.Protected(), middleware.RequireRole("admin"))
	admin.Get("/dashboard/stats", f.UserHandler.GetDashboardStats)
	// User management
	admin.Get("/users", f.UserHandler.GetUsers)
	admin.Put("/users/:userId/role", f.UserHandler.UpdateUserRole)
	admin.Get("/users/:userId/orders", f.UserHandler.GetUserOrders)

	// Order management
	admin.Get("/orders", f.OrderHandler.GetAdminOrders)
	admin.Put("/orders/:orderId/status", f.OrderHandler.UpdateOrderStatus)
	admin.Get("/orders/:orderId/documents", f.OrderHandler.GetOrderDocuments)
	admin.Post("/orders/:orderId/documents", f.OrderHandler.AddAdminOrderDocument)

	// Service management
	admin.Get("/services", f.ServiceHandler.GetServices)
	admin.Post("/services", f.ServiceHandler.CreateService)
	admin.Put("/services/:id", f.ServiceHandler.UpdateService)
	admin.Delete("/services/:id", f.ServiceHandler.DeleteService)

	// Blog management
	admin.Get("/posts", f.PostHandler.ListAll)
	admin.Post("/posts", f.PostHandler.Create)
	admin.Get("/posts/:id", f.PostHandler.GetByID)
	admin.Put("/posts/:id", f.PostHandler.Update)
	admin.Delete("/posts/:id", f.PostHandler.Delete)

	// Deadline management
	admin.Get("/deadlines", f.DeadlineHandler.List)
	admin.Post("/deadlines", f.DeadlineHandler.Create)
	admin.Put("/deadlines/:id", f.DeadlineHandler.Update)
	admin.Delete("/deadlines/:id", f.DeadlineHandler.Delete)

	// Coupon management
	admin.Get("/coupons", f.CouponHandler.ListCoupons)
	admin.Post("/coupons", f.CouponHandler.CreateCoupon)
	admin.Put("/coupons/:id", f.CouponHandler.UpdateCoupon)
	admin.Delete("/coupons/:id", f.CouponHandler.DeleteCoupon)

	// Review management
	admin.Get("/reviews", f.ReviewHandler.ListReviews)
	admin.Post("/reviews", f.ReviewHandler.AdminCreateReview)
	admin.Put("/reviews/:id/status", f.ReviewHandler.UpdateReviewStatus)
	admin.Delete("/reviews/:id", f.ReviewHandler.DeleteReview)

	// Category management
	admin.Get("/categories", f.ServiceHandler.GetServiceCategories)
	admin.Post("/categories", f.ServiceHandler.CreateServiceCategory)
	admin.Put("/categories/:id", f.ServiceHandler.UpdateServiceCategory)
	admin.Delete("/categories/:id", f.ServiceHandler.DeleteServiceCategory)

	// Priority management
	admin.Put("/categories/:id/priority", f.ServiceHandler.UpdateCategoryPriority)
	admin.Put("/services/:id/priority", f.ServiceHandler.UpdateServicePriority)

	// Analytics routes
	admin.Get("/analytics/revenue", f.AnalyticsHandler.GetRevenueAnalytics)
	admin.Get("/analytics/orders", f.AnalyticsHandler.GetOrderAnalytics)
	admin.Get("/analytics/users", f.AnalyticsHandler.GetUserAnalytics)
	admin.Get("/analytics/services", f.AnalyticsHandler.GetServiceAnalytics)
	admin.Get("/analytics/metrics", f.AnalyticsHandler.GetOverallMetrics)
	admin.Get("/analytics/activity", f.AnalyticsHandler.GetRecentActivity)

	// Lead management
	admin.Get("/leads", f.LeadHandler.GetLeads)
	admin.Get("/leads/stats", f.LeadHandler.GetLeadStats)
	admin.Get("/leads/:id", f.LeadHandler.GetLeadByID)
	admin.Put("/leads/:id", f.LeadHandler.UpdateLead)
	admin.Delete("/leads/:id", f.LeadHandler.DeleteLead)

	// Contact management
	admin.Get("/contacts", f.ContactHandler.GetContacts)
	admin.Get("/contacts/stats", f.ContactHandler.GetContactStats)
	admin.Get("/contacts/:id", f.ContactHandler.GetContactByID)
	admin.Put("/contacts/:id", f.ContactHandler.UpdateContactStatus)
	admin.Delete("/contacts/:id", f.ContactHandler.DeleteContact)

	// Settings management
	admin.Get("/settings", f.SettingsHandler.GetAllSettings)
	admin.Get("/settings/categories", f.SettingsHandler.GetSettingsByCategory)
	admin.Get("/settings/:category/:key", f.SettingsHandler.GetSetting)
	admin.Put("/settings/:category/:key", f.SettingsHandler.UpdateSetting)
	admin.Put("/settings/bulk", f.SettingsHandler.UpdateMultipleSettings)
	admin.Post("/settings", f.SettingsHandler.CreateSetting)
	admin.Delete("/settings/:category/:key", f.SettingsHandler.DeleteSetting)
	admin.Post("/settings/test-email", f.SettingsHandler.TestEmailSettings)
	admin.Post("/settings/reset-defaults", f.SettingsHandler.ResetToDefaults)

	// Public settings endpoint (for frontend)
	api.Get("/settings/public", f.SettingsHandler.GetPublicSettings)

	// Public lead creation
	api.Post("/leads", f.LeadHandler.CreateLead)

	// Public contact form submission
	api.Post("/contact", f.ContactHandler.CreateContact)

	app.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	})

	// Static files
	app.Static("/images", "./static/images")
}
