package app

import (
	"github.com/aman4411/protacc-backend/db"
	"github.com/aman4411/protacc-backend/handler"
	"github.com/aman4411/protacc-backend/repository"
	"github.com/aman4411/protacc-backend/service"
)

// Factories holds all the initialized dependencies
type Factories struct {
	UserHandler      *handler.UserHandler
	ServiceHandler   *handler.ServiceHandler
	CartHandler      *handler.CartHandler
	OrderHandler     *handler.OrderHandler
	PaymentHandler   *handler.PaymentHandler
	SettingsHandler  *handler.SettingsHandler
	AnalyticsHandler *handler.AnalyticsHandler
	LeadHandler      *handler.LeadHandler
}

// NewFactories initializes all dependencies and returns them
func NewFactories() (*Factories, error) {
	// Initialize database connection
	db.InitDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Pool)
	serviceRepo := repository.NewServiceRepository(db.Pool)
	cartRepo := repository.NewCartRepository(db.Pool)
	orderRepo := repository.NewOrderRepository(db.Pool)
	settingsRepo := repository.NewSettingsRepository(db.Pool)
	analyticsRepo := repository.NewAnalyticsRepository(db.Pool)
	leadRepo := repository.NewLeadRepository(db.Pool)

	// Initialize services
	mailService := service.NewMailService()
	userService := service.NewUserService(userRepo, mailService)
	serviceService := service.NewServiceService(serviceRepo)
	cartService := service.NewCartService(cartRepo)
	orderService := service.NewOrderService(orderRepo, serviceRepo, cartRepo)
	paymentService := service.NewPaymentService(orderRepo)
	settingsService := service.NewSettingsService(settingsRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	leadService := service.NewLeadService(leadRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, orderService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)
	paymentHandler := handler.NewPaymentHandler(paymentService, orderService)
	settingsHandler := handler.NewSettingsHandler(settingsService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	leadHandler := handler.NewLeadHandler(leadService)

	return &Factories{
		UserHandler:      userHandler,
		ServiceHandler:   serviceHandler,
		CartHandler:      cartHandler,
		OrderHandler:     orderHandler,
		PaymentHandler:   paymentHandler,
		SettingsHandler:  settingsHandler,
		AnalyticsHandler: analyticsHandler,
		LeadHandler:      leadHandler,
	}, nil
}
