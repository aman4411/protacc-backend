package models

import (
	"time"
)

// ContactMessage represents a contact form submission
type ContactMessage struct {
	ID               int        `json:"id" db:"id"`
	Name             string     `json:"name" db:"name" validate:"required,min=2,max=100"`
	Email            string     `json:"email" db:"email" validate:"required,email"`
	Phone            string     `json:"phone" db:"phone" validate:"required,min=10,max=15"`
	Company          *string    `json:"company" db:"company"`
	Subject          string     `json:"subject" db:"subject" validate:"required,min=5,max=200"`
	Message          string     `json:"message" db:"message" validate:"required,min=10,max=1000"`
	ServiceInterest  *string    `json:"service_interest" db:"service_interest"`
	PreferredContact string     `json:"preferred_contact" db:"preferred_contact" validate:"required,oneof=email phone whatsapp"`
	Status           string     `json:"status" db:"status" validate:"oneof=new replied resolved"`
	IPAddress        *string    `json:"ip_address" db:"ip_address"`
	UserAgent        *string    `json:"user_agent" db:"user_agent"`
	CreatedAt        time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at" db:"updated_at"`
	RespondedAt      *time.Time `json:"responded_at" db:"responded_at"`
	RespondedBy      *string    `json:"responded_by" db:"responded_by"`
}

// CreateContactRequest represents the request payload for creating a contact message
type CreateContactRequest struct {
	Name             string  `json:"name" validate:"required,min=2,max=100"`
	Email            string  `json:"email" validate:"required,email"`
	Phone            string  `json:"phone" validate:"required,min=10,max=15"`
	Company          *string `json:"company"`
	Subject          string  `json:"subject" validate:"required,min=5,max=200"`
	Message          string  `json:"message" validate:"required,min=10,max=1000"`
	ServiceInterest  *string `json:"service_interest"`
	PreferredContact string  `json:"preferred_contact" validate:"required,oneof=email phone whatsapp"`
}

// UpdateContactRequest represents the request payload for updating a contact message status
type UpdateContactRequest struct {
	Status      string  `json:"status" validate:"required,oneof=new replied resolved"`
	RespondedBy *string `json:"responded_by"`
}

// ContactFilters represents filters for querying contact messages
type ContactFilters struct {
	Status          string `json:"status"`
	ServiceInterest string `json:"service_interest"`
	DateFrom        string `json:"date_from"`
	DateTo          string `json:"date_to"`
	Search          string `json:"search"`
	Page            int    `json:"page"`
	Limit           int    `json:"limit"`
}

// ContactStats represents contact message statistics
type ContactStats struct {
	TotalMessages    int `json:"total_messages" db:"total_messages"`
	NewMessages      int `json:"new_messages" db:"new_messages"`
	RepliedMessages  int `json:"replied_messages" db:"replied_messages"`
	ResolvedMessages int `json:"resolved_messages" db:"resolved_messages"`
	TodayMessages    int `json:"today_messages" db:"today_messages"`
	WeekMessages     int `json:"week_messages" db:"week_messages"`
	MonthMessages    int `json:"month_messages" db:"month_messages"`
}
