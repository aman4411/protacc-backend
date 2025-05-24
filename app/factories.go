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
}

// NewFactories initializes all dependencies and returns them
func NewFactories() (*Factories, error) {
	// Initialize database connection
	db.InitDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Pool)
	serviceRepo := repository.NewServiceRepository(db.Pool)

	// Initialize services
	mailService := service.NewMailService()
	userService := service.NewUserService(userRepo, mailService)
	serviceService := service.NewServiceService(serviceRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)
	serviceHandler := handler.NewServiceHandler(serviceService)

	return &Factories{
		UserHandler:    userHandler,
		ServiceHandler: serviceHandler,
	}, nil
}
