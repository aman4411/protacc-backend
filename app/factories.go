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
	ContactHandler   *handler.ContactHandler
	ReviewHandler    *handler.ReviewHandler
	CouponHandler    *handler.CouponHandler
	DeadlineHandler  *handler.DeadlineHandler
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
	orderDocumentRepo := repository.NewOrderDocumentRepository(db.Pool)
	settingsRepo := repository.NewSettingsRepository(db.Pool)
	analyticsRepo := repository.NewAnalyticsRepository(db.Pool)
	leadRepo := repository.NewLeadRepository(db.Pool)
	contactRepo := repository.NewContactRepository(db.Pool)
	reviewRepo := repository.NewReviewRepository(db.Pool)
	couponRepo := repository.NewCouponRepository(db.Pool)
	deadlineRepo := repository.NewDeadlineRepository(db.Pool)

	// Initialize services
	mailService := service.NewMailService()
	userService := service.NewUserService(userRepo, mailService)
	serviceService := service.NewServiceService(serviceRepo)
	cartService := service.NewCartService(cartRepo)
	couponService := service.NewCouponService(couponRepo)
	deadlineService := service.NewDeadlineService(deadlineRepo)
	orderService := service.NewOrderService(orderRepo, serviceRepo, cartRepo, couponService)
	orderDocumentService := service.NewOrderDocumentService(orderDocumentRepo, orderRepo)
	paymentService := service.NewPaymentService(orderRepo)
	settingsService := service.NewSettingsService(settingsRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	leadService := service.NewLeadService(leadRepo)
	contactService := service.NewContactService(contactRepo)
	reviewService := service.NewReviewService(reviewRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, orderService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	cartHandler := handler.NewCartHandler(cartService)
	orderHandler := handler.NewOrderHandler(orderService, orderDocumentService)
	paymentHandler := handler.NewPaymentHandler(paymentService, orderService)
	settingsHandler := handler.NewSettingsHandler(settingsService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	leadHandler := handler.NewLeadHandler(leadService)
	contactHandler := handler.NewContactHandler(contactService)
	reviewHandler := handler.NewReviewHandler(reviewService)
	couponHandler := handler.NewCouponHandler(couponService)
	deadlineHandler := handler.NewDeadlineHandler(deadlineService)

	return &Factories{
		UserHandler:      userHandler,
		ServiceHandler:   serviceHandler,
		CartHandler:      cartHandler,
		OrderHandler:     orderHandler,
		PaymentHandler:   paymentHandler,
		SettingsHandler:  settingsHandler,
		AnalyticsHandler: analyticsHandler,
		LeadHandler:      leadHandler,
		ContactHandler:   contactHandler,
		ReviewHandler:    reviewHandler,
		CouponHandler:    couponHandler,
		DeadlineHandler:  deadlineHandler,
	}, nil
}
