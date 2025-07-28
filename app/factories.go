package app

import (
	"github.com/aman4411/protacc-backend/db"
	"github.com/aman4411/protacc-backend/handler"
	"github.com/aman4411/protacc-backend/repository"
	"github.com/aman4411/protacc-backend/service"
)

// Factories holds all the initialized dependencies
type Factories struct {
	UserHandler    *handler.UserHandler
	ServiceHandler *handler.ServiceHandler
	CartHandler    *handler.CartHandler
	OrderHandler   *handler.OrderHandler
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

	// Initialize services
	mailService := service.NewMailService()
	userService := service.NewUserService(userRepo, mailService)
	serviceService := service.NewServiceService(serviceRepo)
	cartService := service.NewCartService(cartRepo)
	orderService := service.NewOrderService(orderRepo, serviceRepo, cartRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService)

	return &Factories{
		UserHandler:    userHandler,
		ServiceHandler: serviceHandler,
		CartHandler:    cartHandler,
		OrderHandler:   orderHandler,
	}, nil
}
