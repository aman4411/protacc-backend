package app

import (
	"github.com/aman4411/protacc-backend/db"
	"github.com/aman4411/protacc-backend/handler"
	"github.com/aman4411/protacc-backend/repository"
	"github.com/aman4411/protacc-backend/service"
)

// Factories holds all the initialized dependencies
type Factories struct {
	UserHandler *handler.UserHandler
}

// NewFactories initializes all dependencies and returns them
func NewFactories() (*Factories, error) {
	// Initialize database connection
	db.InitDB()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.Pool)

	// Initialize services
	mailService := service.NewMailService()
	userService := service.NewUserService(userRepo, mailService)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService)

	return &Factories{
		UserHandler: userHandler,
	}, nil
}
