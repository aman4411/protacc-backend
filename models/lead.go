package models

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/lib/pq"
)

// CustomDate handles date-only JSON parsing (YYYY-MM-DD format)
type CustomDate struct {
	time.Time
}

// UnmarshalJSON implements json.Unmarshaler interface for CustomDate
func (cd *CustomDate) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	dateStr := strings.Trim(string(data), "\"")

	// Parse date in YYYY-MM-DD format
	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return err
	}

	cd.Time = t
	return nil
}

// MarshalJSON implements json.Marshaler interface for CustomDate
func (cd CustomDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(cd.Time.Format("2006-01-02"))
}

type LeadStatus string
type LeadPriority string
type ContactMethod string

const (
	LeadStatusNew        LeadStatus = "new"
	LeadStatusContacted  LeadStatus = "contacted"
	LeadStatusInProgress LeadStatus = "in_progress"
	LeadStatusQualified  LeadStatus = "qualified"
	LeadStatusConverted  LeadStatus = "converted"
	LeadStatusRejected   LeadStatus = "rejected"
	LeadStatusClosed     LeadStatus = "closed"
)

const (
	PriorityLow    LeadPriority = "low"
	PriorityMedium LeadPriority = "medium"
	PriorityHigh   LeadPriority = "high"
	PriorityUrgent LeadPriority = "urgent"
)

const (
	ContactEmail ContactMethod = "email"
	ContactPhone ContactMethod = "phone"
	ContactBoth  ContactMethod = "both"
)

type BusinessLead struct {
	ID                     int            `json:"id" db:"id"`
	FirstName              string         `json:"first_name" db:"first_name"`
	LastName               string         `json:"last_name" db:"last_name"`
	Email                  string         `json:"email" db:"email"`
	Phone                  string         `json:"phone" db:"phone"`
	CompanyName            *string        `json:"company_name" db:"company_name"`
	BusinessType           *string        `json:"business_type" db:"business_type"`
	ServicesInterested     pq.StringArray `json:"services_interested" db:"services_interested"`
	BudgetRange            *string        `json:"budget_range" db:"budget_range"`
	PreferredContactMethod ContactMethod  `json:"preferred_contact_method" db:"preferred_contact_method"`
	Message                *string        `json:"message" db:"message"`
	Status                 LeadStatus     `json:"status" db:"status"`
	Priority               LeadPriority   `json:"priority" db:"priority"`
	AssignedTo             *string        `json:"assigned_to" db:"assigned_to"`
	Source                 string         `json:"source" db:"source"`
	FollowUpDate           *time.Time     `json:"follow_up_date" db:"follow_up_date"`
	Notes                  *string        `json:"notes" db:"notes"`
	CreatedAt              time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at" db:"updated_at"`

	// Populated via joins
	AssignedUser *User `json:"assigned_user,omitempty" db:"-"`
}

type CreateLeadRequest struct {
	FirstName              string   `json:"first_name" validate:"required,min=2,max=100"`
	LastName               string   `json:"last_name" validate:"required,min=2,max=100"`
	Email                  string   `json:"email" validate:"required,email"`
	Phone                  string   `json:"phone" validate:"required,min=10,max=20"`
	CompanyName            string   `json:"company_name,omitempty"`
	BusinessType           string   `json:"business_type,omitempty"`
	ServicesInterested     []string `json:"services_interested"`
	BudgetRange            string   `json:"budget_range,omitempty"`
	PreferredContactMethod string   `json:"preferred_contact_method" validate:"omitempty,oneof=email phone both"`
	Message                string   `json:"message,omitempty"`
}

type UpdateLeadRequest struct {
	Status       string `json:"status,omitempty" validate:"omitempty,oneof=new contacted in_progress qualified converted rejected closed"`
	Priority     string `json:"priority,omitempty" validate:"omitempty,oneof=low medium high urgent"`
	AssignedTo   string `json:"assigned_to,omitempty"`
	FollowUpDate string `json:"follow_up_date,omitempty"`
	Notes        string `json:"notes,omitempty"`
}

type LeadFilters struct {
	Status     string `json:"status,omitempty"`
	Priority   string `json:"priority,omitempty"`
	AssignedTo string `json:"assigned_to,omitempty"`
	Search     string `json:"search,omitempty"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
}

type LeadStats struct {
	TotalLeads      int     `json:"total_leads"`
	NewLeads        int     `json:"new_leads"`
	InProgressLeads int     `json:"in_progress_leads"`
	ConvertedLeads  int     `json:"converted_leads"`
	ConversionRate  float64 `json:"conversion_rate"`
}
