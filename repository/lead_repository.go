package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LeadRepository struct {
	db *pgxpool.Pool
}

func NewLeadRepository(db *pgxpool.Pool) *LeadRepository {
	return &LeadRepository{db: db}
}

func (r *LeadRepository) CreateLead(ctx context.Context, lead *models.CreateLeadRequest) (*models.BusinessLead, error) {
	query := `
		INSERT INTO business_leads (
			first_name, last_name, email, phone, company_name, business_type,
			services_interested, budget_range, preferred_contact_method, message, source
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 'website')
		RETURNING id, first_name, last_name, email, phone, company_name, business_type,
				  services_interested, budget_range, preferred_contact_method, message,
				  status, priority, assigned_to, source, follow_up_date, notes, created_at, updated_at`

	var newLead models.BusinessLead
	var companyName, businessType, budgetRange, message *string
	var preferredContactMethod string

	if lead.CompanyName != "" {
		companyName = &lead.CompanyName
	}
	if lead.BusinessType != "" {
		businessType = &lead.BusinessType
	}
	if lead.BudgetRange != "" {
		budgetRange = &lead.BudgetRange
	}
	if lead.Message != "" {
		message = &lead.Message
	}
	if lead.PreferredContactMethod != "" {
		preferredContactMethod = lead.PreferredContactMethod
	} else {
		preferredContactMethod = "email"
	}

	err := r.db.QueryRow(ctx, query,
		lead.FirstName, lead.LastName, lead.Email, lead.Phone,
		companyName, businessType, lead.ServicesInterested, budgetRange,
		preferredContactMethod, message,
	).Scan(
		&newLead.ID, &newLead.FirstName, &newLead.LastName, &newLead.Email,
		&newLead.Phone, &newLead.CompanyName, &newLead.BusinessType,
		&newLead.ServicesInterested, &newLead.BudgetRange, &newLead.PreferredContactMethod,
		&newLead.Message, &newLead.Status, &newLead.Priority, &newLead.AssignedTo,
		&newLead.Source, &newLead.FollowUpDate, &newLead.Notes, &newLead.CreatedAt, &newLead.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &newLead, nil
}

func (r *LeadRepository) GetLeads(ctx context.Context, filters models.LeadFilters) ([]models.BusinessLead, int, error) {
	var whereConditions []string
	var args []interface{}
	argIndex := 1

	// Build WHERE conditions
	if filters.Status != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("bl.status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if filters.Priority != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("bl.priority = $%d", argIndex))
		args = append(args, filters.Priority)
		argIndex++
	}

	if filters.AssignedTo != "" {
		whereConditions = append(whereConditions, fmt.Sprintf("bl.assigned_to = $%d", argIndex))
		args = append(args, filters.AssignedTo)
		argIndex++
	}

	if filters.Search != "" {
		searchCondition := fmt.Sprintf(`(
			bl.first_name ILIKE $%d OR bl.last_name ILIKE $%d OR 
			bl.email ILIKE $%d OR bl.phone ILIKE $%d OR 
			bl.company_name ILIKE $%d
		)`, argIndex, argIndex, argIndex, argIndex, argIndex)
		whereConditions = append(whereConditions, searchCondition)
		args = append(args, "%"+filters.Search+"%")
		argIndex++
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = "WHERE " + strings.Join(whereConditions, " AND ")
	}

	// Count total records
	countQuery := fmt.Sprintf(`
		SELECT COUNT(*) FROM business_leads bl 
		LEFT JOIN users u ON bl.assigned_to = u.id 
		%s`, whereClause)

	var total int
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (filters.Page - 1) * filters.Limit

	// Main query with pagination
	query := fmt.Sprintf(`
		SELECT bl.id, bl.first_name, bl.last_name, bl.email, bl.phone,
			   bl.company_name, bl.business_type, bl.services_interested, bl.budget_range,
			   bl.preferred_contact_method, bl.message, bl.status, bl.priority,
			   bl.assigned_to, bl.source, bl.follow_up_date, bl.notes, 
			   bl.created_at, bl.updated_at,
			   u.first_name as assigned_first_name, u.last_name as assigned_last_name,
			   u.email as assigned_email
		FROM business_leads bl
		LEFT JOIN users u ON bl.assigned_to = u.id
		%s
		ORDER BY bl.created_at DESC
		LIMIT $%d OFFSET $%d`, whereClause, argIndex, argIndex+1)

	args = append(args, filters.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var leads []models.BusinessLead
	for rows.Next() {
		var lead models.BusinessLead
		var assignedFirstName, assignedLastName, assignedEmail *string

		err := rows.Scan(
			&lead.ID, &lead.FirstName, &lead.LastName, &lead.Email, &lead.Phone,
			&lead.CompanyName, &lead.BusinessType, &lead.ServicesInterested, &lead.BudgetRange,
			&lead.PreferredContactMethod, &lead.Message, &lead.Status, &lead.Priority,
			&lead.AssignedTo, &lead.Source, &lead.FollowUpDate, &lead.Notes,
			&lead.CreatedAt, &lead.UpdatedAt,
			&assignedFirstName, &assignedLastName, &assignedEmail,
		)
		if err != nil {
			return nil, 0, err
		}

		// Populate assigned user if exists
		if assignedFirstName != nil && assignedLastName != nil && assignedEmail != nil {
			lead.AssignedUser = &models.User{
				ID:        *lead.AssignedTo,
				FirstName: *assignedFirstName,
				LastName:  *assignedLastName,
				Email:     *assignedEmail,
			}
		}

		leads = append(leads, lead)
	}

	return leads, total, nil
}

func (r *LeadRepository) GetLeadByID(ctx context.Context, id int) (*models.BusinessLead, error) {
	query := `
		SELECT bl.id, bl.first_name, bl.last_name, bl.email, bl.phone,
			   bl.company_name, bl.business_type, bl.services_interested, bl.budget_range,
			   bl.preferred_contact_method, bl.message, bl.status, bl.priority,
			   bl.assigned_to, bl.source, bl.follow_up_date, bl.notes, 
			   bl.created_at, bl.updated_at,
			   u.first_name as assigned_first_name, u.last_name as assigned_last_name,
			   u.email as assigned_email
		FROM business_leads bl
		LEFT JOIN users u ON bl.assigned_to = u.id
		WHERE bl.id = $1`

	var lead models.BusinessLead
	var assignedFirstName, assignedLastName, assignedEmail *string

	err := r.db.QueryRow(ctx, query, id).Scan(
		&lead.ID, &lead.FirstName, &lead.LastName, &lead.Email, &lead.Phone,
		&lead.CompanyName, &lead.BusinessType, &lead.ServicesInterested, &lead.BudgetRange,
		&lead.PreferredContactMethod, &lead.Message, &lead.Status, &lead.Priority,
		&lead.AssignedTo, &lead.Source, &lead.FollowUpDate, &lead.Notes,
		&lead.CreatedAt, &lead.UpdatedAt,
		&assignedFirstName, &assignedLastName, &assignedEmail,
	)

	if err != nil {
		return nil, err
	}

	// Populate assigned user if exists
	if assignedFirstName != nil && assignedLastName != nil && assignedEmail != nil {
		lead.AssignedUser = &models.User{
			ID:        *lead.AssignedTo,
			FirstName: *assignedFirstName,
			LastName:  *assignedLastName,
			Email:     *assignedEmail,
		}
	}

	return &lead, nil
}

func (r *LeadRepository) UpdateLead(ctx context.Context, id int, updates *models.UpdateLeadRequest) error {
	var setParts []string
	var args []interface{}
	argIndex := 1

	if updates.Status != "" {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIndex))
		args = append(args, updates.Status)
		argIndex++
	}

	if updates.Priority != "" {
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argIndex))
		args = append(args, updates.Priority)
		argIndex++
	}

	if updates.AssignedTo != "" {
		if updates.AssignedTo == "null" {
			setParts = append(setParts, "assigned_to = NULL")
		} else {
			setParts = append(setParts, fmt.Sprintf("assigned_to = $%d", argIndex))
			args = append(args, updates.AssignedTo)
			argIndex++
		}
	}

	if updates.FollowUpDate != "" {
		// Parse the date string
		followUpDate, err := time.Parse("2006-01-02", updates.FollowUpDate)
		if err != nil {
			return fmt.Errorf("invalid follow_up_date format: %v", err)
		}
		setParts = append(setParts, fmt.Sprintf("follow_up_date = $%d", argIndex))
		args = append(args, followUpDate)
		argIndex++
	}

	if updates.Notes != "" {
		setParts = append(setParts, fmt.Sprintf("notes = $%d", argIndex))
		args = append(args, updates.Notes)
		argIndex++
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	query := fmt.Sprintf(`
		UPDATE business_leads 
		SET %s, updated_at = NOW() 
		WHERE id = $%d`, strings.Join(setParts, ", "), argIndex)
	args = append(args, id)

	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}

func (r *LeadRepository) DeleteLead(ctx context.Context, id int) error {
	query := `DELETE FROM business_leads WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("lead not found")
	}

	return nil
}

func (r *LeadRepository) GetLeadStats(ctx context.Context) (*models.LeadStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_leads,
			COUNT(CASE WHEN status = 'new' THEN 1 END) as new_leads,
			COUNT(CASE WHEN status IN ('in_progress', 'contacted', 'qualified') THEN 1 END) as in_progress_leads,
			COUNT(CASE WHEN status = 'converted' THEN 1 END) as converted_leads
		FROM business_leads`

	var stats models.LeadStats
	err := r.db.QueryRow(ctx, query).Scan(
		&stats.TotalLeads, &stats.NewLeads, &stats.InProgressLeads, &stats.ConvertedLeads,
	)
	if err != nil {
		return nil, err
	}

	// Calculate conversion rate
	if stats.TotalLeads > 0 {
		stats.ConversionRate = float64(stats.ConvertedLeads) / float64(stats.TotalLeads) * 100
	}

	return &stats, nil
}
