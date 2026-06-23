package models

import "time"

// Deadline is an admin-managed tax/compliance due date shown on the homepage.
type Deadline struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	DueDate     time.Time `json:"due_date"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
