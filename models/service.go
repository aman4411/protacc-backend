package models

import (
	"time"

	"github.com/gosimple/slug"
)

type ServiceStatus string

const (
	ServiceStatusActive   ServiceStatus = "active"
	ServiceStatusInactive ServiceStatus = "inactive"
)

type ServiceCategory struct {
	ID          int           `json:"id" db:"id"`
	Name        string        `json:"name" db:"name"`
	Slug        string        `json:"slug" db:"slug"`
	Description string        `json:"description" db:"description"`
	Icon        string        `json:"icon" db:"icon"`
	Status      ServiceStatus `json:"status" db:"status"`
	CreatedAt   time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" db:"updated_at"`
}

type Service struct {
	ID                    int              `json:"id" db:"id"`
	CategoryID            int              `json:"category_id" db:"category_id"`
	Name                  string           `json:"name" db:"name"`
	Slug                  string           `json:"slug" db:"slug"`
	Description           string           `json:"description" db:"description"`
	ShortDescription      string           `json:"short_description" db:"short_description"`
	Features              []string         `json:"features" db:"features"`
	Requirements          []string         `json:"requirements" db:"requirements"`
	Price                 float64          `json:"price" db:"price"`
	BookingAmount         float64          `json:"booking_amount" db:"booking_amount"`
	EstimatedDeliveryDays int              `json:"estimated_delivery_days" db:"estimated_delivery_days"`
	Icon                  string           `json:"icon" db:"icon"`
	Status                ServiceStatus    `json:"status" db:"status"`
	CreatedAt             time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at" db:"updated_at"`
	Category              *ServiceCategory `json:"category,omitempty" db:"-"`
}

func (s *Service) GenerateSlug() {
	s.Slug = slug.Make(s.Name)
}

type CartItem struct {
	ID        int       `json:"id" db:"id"`
	UserID    string    `json:"user_id" db:"user_id"`
	ServiceID int       `json:"service_id" db:"service_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Service   *Service  `json:"service,omitempty" db:"-"`
}

type OrderStatus string

const (
	OrderStatusPendingPayment    OrderStatus = "pending_payment"
	OrderStatusPaymentReceived   OrderStatus = "payment_received"
	OrderStatusProcessing        OrderStatus = "processing"
	OrderStatusDocumentsRequired OrderStatus = "documents_required"
	OrderStatusDocumentsReceived OrderStatus = "documents_received"
	OrderStatusInProgress        OrderStatus = "in_progress"
	OrderStatusCompleted         OrderStatus = "completed"
	OrderStatusCancelled         OrderStatus = "cancelled"
)

type OrderItem struct {
	ID            int       `json:"id" db:"id"`
	OrderID       int       `json:"order_id" db:"order_id"`
	ServiceID     int       `json:"service_id" db:"service_id"`
	Quantity      int       `json:"quantity" db:"quantity"`
	Price         float64   `json:"price" db:"price"`
	BookingAmount float64   `json:"booking_amount" db:"booking_amount"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	Service       *Service  `json:"service,omitempty" db:"-"`
}

type Order struct {
	ID                int         `json:"id" db:"id"`
	UserID            string      `json:"user_id" db:"user_id"`
	OrderNumber       string      `json:"order_number" db:"order_number"`
	TotalAmount       float64     `json:"total_amount" db:"total_amount"`
	BookingAmount     float64     `json:"booking_amount" db:"booking_amount"`
	RemainingAmount   float64     `json:"remaining_amount" db:"remaining_amount"`
	Status            OrderStatus `json:"status" db:"status"`
	PaymentStatus     bool        `json:"payment_status" db:"payment_status"`
	RazorpayOrderID   *string     `json:"razorpay_order_id,omitempty" db:"razorpay_order_id"`
	RazorpayPaymentID *string     `json:"razorpay_payment_id,omitempty" db:"razorpay_payment_id"`
	PaymentMethod     *string     `json:"payment_method,omitempty" db:"payment_method"`
	PaymentGateway    *string     `json:"payment_gateway,omitempty" db:"payment_gateway"`
	Notes             *string     `json:"notes,omitempty" db:"notes"`
	CreatedAt         time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at" db:"updated_at"`
	Items             []OrderItem `json:"items,omitempty" db:"-"`
	User              *User       `json:"user,omitempty" db:"-"`
}

type OrderStatusHistory struct {
	ID        int         `json:"id" db:"id"`
	OrderID   int         `json:"order_id" db:"order_id"`
	Status    OrderStatus `json:"status" db:"status"`
	Notes     *string     `json:"notes,omitempty" db:"notes"`
	CreatedBy *string     `json:"created_by,omitempty" db:"created_by"`
	CreatedAt time.Time   `json:"created_at" db:"created_at"`
	User      *User       `json:"user,omitempty" db:"-"`
}
