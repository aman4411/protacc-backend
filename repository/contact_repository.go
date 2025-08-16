package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository interface {
	CreateContact(contact *models.CreateContactRequest, ipAddress, userAgent string) (*models.ContactMessage, error)
	GetContacts(filters models.ContactFilters) ([]models.ContactMessage, int, error)
	GetContactByID(id int) (*models.ContactMessage, error)
	UpdateContactStatus(id int, updates models.UpdateContactRequest) error
	DeleteContact(id int) error
	GetContactStats() (*models.ContactStats, error)
}

type contactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(db *pgxpool.Pool) ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) CreateContact(contact *models.CreateContactRequest, ipAddress, userAgent string) (*models.ContactMessage, error) {
	query := `
		INSERT INTO contact_messages (
			name, email, phone, company, subject, message, 
			service_interest, preferred_contact, ip_address, user_agent
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, name, email, phone, company, subject, message, 
				  service_interest, preferred_contact, status, ip_address, 
				  user_agent, created_at, updated_at, responded_at, responded_by
	`

	var result models.ContactMessage
	err := r.db.QueryRow(
		context.Background(),
		query,
		contact.Name,
		contact.Email,
		contact.Phone,
		contact.Company,
		contact.Subject,
		contact.Message,
		contact.ServiceInterest,
		contact.PreferredContact,
		ipAddress,
		userAgent,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Email,
		&result.Phone,
		&result.Company,
		&result.Subject,
		&result.Message,
		&result.ServiceInterest,
		&result.PreferredContact,
		&result.Status,
		&result.IPAddress,
		&result.UserAgent,
		&result.CreatedAt,
		&result.UpdatedAt,
		&result.RespondedAt,
		&result.RespondedBy,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create contact message: %w", err)
	}

	return &result, nil
}

func (r *contactRepository) GetContacts(filters models.ContactFilters) ([]models.ContactMessage, int, error) {
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.ServiceInterest != "" {
		conditions = append(conditions, fmt.Sprintf("service_interest = $%d", argIndex))
		args = append(args, filters.ServiceInterest)
		argIndex++
	}

	if filters.DateFrom != "" {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", argIndex))
		args = append(args, filters.DateFrom)
		argIndex++
	}

	if filters.DateTo != "" {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", argIndex))
		args = append(args, filters.DateTo)
		argIndex++
	}

	if filters.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d OR subject ILIKE $%d OR message ILIKE $%d)", argIndex, argIndex, argIndex, argIndex))
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM contact_messages %s", whereClause)
	var totalCount int
	err := r.db.QueryRow(context.Background(), countQuery, args...).Scan(&totalCount)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get contact messages count: %w", err)
	}

	// Set pagination defaults
	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 {
		filters.Limit = 10
	}

	offset := (filters.Page - 1) * filters.Limit

	// Get contacts with pagination
	query := fmt.Sprintf(`
		SELECT id, name, email, phone, company, subject, message, 
			   service_interest, preferred_contact, status, ip_address, 
			   user_agent, created_at, updated_at, responded_at, responded_by
		FROM contact_messages %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, filters.Limit, offset)

	rows, err := r.db.Query(context.Background(), query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get contact messages: %w", err)
	}
	defer rows.Close()

	var contacts []models.ContactMessage
	for rows.Next() {
		var contact models.ContactMessage
		err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Email,
			&contact.Phone,
			&contact.Company,
			&contact.Subject,
			&contact.Message,
			&contact.ServiceInterest,
			&contact.PreferredContact,
			&contact.Status,
			&contact.IPAddress,
			&contact.UserAgent,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.RespondedAt,
			&contact.RespondedBy,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan contact message: %w", err)
		}
		contacts = append(contacts, contact)
	}

	return contacts, totalCount, nil
}

func (r *contactRepository) GetContactByID(id int) (*models.ContactMessage, error) {
	query := `
		SELECT id, name, email, phone, company, subject, message, 
			   service_interest, preferred_contact, status, ip_address, 
			   user_agent, created_at, updated_at, responded_at, responded_by
		FROM contact_messages 
		WHERE id = $1
	`

	var contact models.ContactMessage
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&contact.ID,
		&contact.Name,
		&contact.Email,
		&contact.Phone,
		&contact.Company,
		&contact.Subject,
		&contact.Message,
		&contact.ServiceInterest,
		&contact.PreferredContact,
		&contact.Status,
		&contact.IPAddress,
		&contact.UserAgent,
		&contact.CreatedAt,
		&contact.UpdatedAt,
		&contact.RespondedAt,
		&contact.RespondedBy,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("contact message not found")
		}
		return nil, fmt.Errorf("failed to get contact message: %w", err)
	}

	return &contact, nil
}

func (r *contactRepository) UpdateContactStatus(id int, updates models.UpdateContactRequest) error {
	query := `
		UPDATE contact_messages 
		SET status = $1, 
			responded_by = $2, 
			responded_at = CASE 
				WHEN $3 IN ('replied', 'resolved') AND responded_at IS NULL 
				THEN CURRENT_TIMESTAMP 
				ELSE responded_at 
			END
		WHERE id = $4
	`

	result, err := r.db.Exec(context.Background(), query, updates.Status, updates.RespondedBy, updates.Status, id)
	if err != nil {
		return fmt.Errorf("failed to update contact message: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("contact message not found")
	}

	return nil
}

func (r *contactRepository) DeleteContact(id int) error {
	query := "DELETE FROM contact_messages WHERE id = $1"

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("failed to delete contact message: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("contact message not found")
	}

	return nil
}

func (r *contactRepository) GetContactStats() (*models.ContactStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_messages,
			COUNT(*) FILTER (WHERE status = 'new') as new_messages,
			COUNT(*) FILTER (WHERE status = 'replied') as replied_messages,
			COUNT(*) FILTER (WHERE status = 'resolved') as resolved_messages,
			COUNT(*) FILTER (WHERE DATE(created_at) = CURRENT_DATE) as today_messages,
			COUNT(*) FILTER (WHERE created_at >= CURRENT_DATE - INTERVAL '7 days') as week_messages,
			COUNT(*) FILTER (WHERE created_at >= CURRENT_DATE - INTERVAL '30 days') as month_messages
		FROM contact_messages
	`

	var stats models.ContactStats
	err := r.db.QueryRow(context.Background(), query).Scan(
		&stats.TotalMessages,
		&stats.NewMessages,
		&stats.RepliedMessages,
		&stats.ResolvedMessages,
		&stats.TodayMessages,
		&stats.WeekMessages,
		&stats.MonthMessages,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get contact stats: %w", err)
	}

	return &stats, nil
}
