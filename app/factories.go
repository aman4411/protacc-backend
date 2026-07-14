package app

import (
	"github.com/aman4411/protacc-backend/cache"
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
	PostHandler      *handler.PostHandler
}

// NewFactories initializes all dependencies and returns them
func NewFactories() (*Factories, error) {
	// Initialize database connection
	db.InitDB()

	// Shared in-process cache for public, read-heavy endpoints. On the single
	// Render instance this serves repeated catalog reads from memory, so the DB
	// sees roughly one query per key per TTL window regardless of traffic.
	appCache := cache.New()

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
	postRepo := repository.NewPostRepository(db.Pool)

	// Initialize services
	mailService := service.NewMailService()
	userService := service.NewUserService(userRepo, mailService)
	serviceService := service.NewServiceService(serviceRepo, appCache)
	cartService := service.NewCartService(cartRepo)
	couponService := service.NewCouponService(couponRepo, appCache)
	deadlineService := service.NewDeadlineService(deadlineRepo, appCache)
	postService := service.NewPostService(postRepo, appCache)
	orderService := service.NewOrderService(orderRepo, serviceRepo, cartRepo, couponService, mailService)
	orderDocumentService := service.NewOrderDocumentService(orderDocumentRepo, orderRepo)
	paymentService := service.NewPaymentService(orderRepo, mailService)
	settingsService := service.NewSettingsService(settingsRepo, appCache)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	leadService := service.NewLeadService(leadRepo, mailService)
	contactService := service.NewContactService(contactRepo, mailService)
	reviewService := service.NewReviewService(reviewRepo, appCache)

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
	postHandler := handler.NewPostHandler(postService)

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
		PostHandler:      postHandler,
	}, nil
}
