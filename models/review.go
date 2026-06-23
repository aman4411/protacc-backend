package models

import "time"

type ReviewStatus string

const (
	ReviewStatusPublished ReviewStatus = "published"
	ReviewStatusHidden    ReviewStatus = "hidden"
)

// Review is a customer rating + comment for a service.
type Review struct {
	ID        int          `json:"id"`
	ServiceID int          `json:"service_id"`
	UserID    string       `json:"user_id"`
	Rating    int          `json:"rating"`
	Comment   string       `json:"comment"`
	Status    ReviewStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	// Derived display fields (populated via joins, never stored).
	ReviewerName string `json:"reviewer_name,omitempty"`
	ServiceName  string `json:"service_name,omitempty"`
	ServiceSlug  string `json:"service_slug,omitempty"`
}

// ReviewSummary is the aggregate rating for a service.
type ReviewSummary struct {
	Average float64 `json:"average"`
	Count   int     `json:"count"`
}
