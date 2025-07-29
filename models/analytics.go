package models

import "time"

// Revenue Analytics
type RevenueAnalytics struct {
	TotalRevenue   float64              `json:"total_revenue"`
	RevenueGrowth  float64              `json:"revenue_growth"`
	DailyRevenue   []DailyRevenueData   `json:"daily_revenue"`
	MonthlyRevenue []MonthlyRevenueData `json:"monthly_revenue"`
	PaymentMethods []PaymentMethodData  `json:"payment_methods"`
}

type DailyRevenueData struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Orders int     `json:"orders"`
}

type MonthlyRevenueData struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
	Orders int     `json:"orders"`
}

type PaymentMethodData struct {
	Method string  `json:"method"`
	Amount float64 `json:"amount"`
	Count  int     `json:"count"`
}

// Order Analytics
type OrderAnalytics struct {
	TotalOrders          int               `json:"total_orders"`
	OrderGrowth          float64           `json:"order_growth"`
	AverageOrderValue    float64           `json:"average_order_value"`
	OrderStatusBreakdown []OrderStatusData `json:"order_status_breakdown"`
	DailyOrders          []DailyOrderData  `json:"daily_orders"`
	CompletionRate       float64           `json:"completion_rate"`
}

type OrderStatusData struct {
	Status     string  `json:"status"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

type DailyOrderData struct {
	Date    string  `json:"date"`
	Orders  int     `json:"orders"`
	Revenue float64 `json:"revenue"`
}

// User Analytics
type UserAnalytics struct {
	TotalUsers        int             `json:"total_users"`
	UserGrowth        float64         `json:"user_growth"`
	DailySignups      []DailyUserData `json:"daily_signups"`
	UserRoleBreakdown []UserRoleData  `json:"user_role_breakdown"`
	VerificationRate  float64         `json:"verification_rate"`
}

type DailyUserData struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type UserRoleData struct {
	Role       string  `json:"role"`
	Count      int     `json:"count"`
	Percentage float64 `json:"percentage"`
}

// Service Analytics
type ServiceAnalytics struct {
	TotalServices       int                       `json:"total_services"`
	PopularServices     []PopularServiceData      `json:"popular_services"`
	CategoryPerformance []CategoryPerformanceData `json:"category_performance"`
	ServiceConversion   float64                   `json:"service_conversion"`
}

type PopularServiceData struct {
	ServiceID   int     `json:"service_id"`
	ServiceName string  `json:"service_name"`
	OrderCount  int     `json:"order_count"`
	Revenue     float64 `json:"revenue"`
}

type CategoryPerformanceData struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	ServiceCount int     `json:"service_count"`
	OrderCount   int     `json:"order_count"`
	Revenue      float64 `json:"revenue"`
}

// Overall Metrics
type OverallMetrics struct {
	TotalRevenue          float64 `json:"total_revenue"`
	TotalOrders           int     `json:"total_orders"`
	TotalUsers            int     `json:"total_users"`
	AverageOrderValue     float64 `json:"average_order_value"`
	ConversionRate        float64 `json:"conversion_rate"`
	CustomerLifetimeValue float64 `json:"customer_lifetime_value"`
	RevenueGrowth         float64 `json:"revenue_growth"`
	OrderGrowth           float64 `json:"order_growth"`
	UserGrowth            float64 `json:"user_growth"`
}

// Recent Activity
type RecentActivity struct {
	Activities []ActivityData `json:"activities"`
}

type ActivityData struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
	UserName    string    `json:"user_name"`
	Amount      *float64  `json:"amount,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
