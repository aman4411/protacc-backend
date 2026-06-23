package repository

import (
	"context"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DeadlineRepository struct {
	db *pgxpool.Pool
}

func NewDeadlineRepository(db *pgxpool.Pool) *DeadlineRepository {
	return &DeadlineRepository{db: db}
}

const deadlineCols = `id, title, description, category, due_date, is_active, created_at, updated_at`

func scanDeadlines(rows interface {
	Next() bool
	Scan(...any) error
}) ([]models.Deadline, error) {
	out := []models.Deadline{}
	for rows.Next() {
		var d models.Deadline
		if err := rows.Scan(&d.ID, &d.Title, &d.Description, &d.Category, &d.DueDate, &d.IsActive, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, nil
}

// List returns all deadlines (admin), soonest first.
func (r *DeadlineRepository) List(ctx context.Context) ([]models.Deadline, error) {
	rows, err := r.db.Query(ctx, `SELECT `+deadlineCols+` FROM deadlines ORDER BY due_date ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDeadlines(rows)
}

// ListUpcoming returns active deadlines on/after today (public).
func (r *DeadlineRepository) ListUpcoming(ctx context.Context, limit int) ([]models.Deadline, error) {
	rows, err := r.db.Query(ctx, `
		SELECT `+deadlineCols+` FROM deadlines
		WHERE is_active = true AND due_date >= CURRENT_DATE
		ORDER BY due_date ASC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanDeadlines(rows)
}

func (r *DeadlineRepository) Create(ctx context.Context, d *models.Deadline) (*models.Deadline, error) {
	err := r.db.QueryRow(ctx, `
		INSERT INTO deadlines (title, description, category, due_date, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`,
		d.Title, d.Description, d.Category, d.DueDate, d.IsActive,
	).Scan(&d.ID, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *DeadlineRepository) Update(ctx context.Context, d *models.Deadline) (*models.Deadline, error) {
	err := r.db.QueryRow(ctx, `
		UPDATE deadlines SET title=$1, description=$2, category=$3, due_date=$4, is_active=$5, updated_at=NOW()
		WHERE id=$6
		RETURNING created_at, updated_at`,
		d.Title, d.Description, d.Category, d.DueDate, d.IsActive, d.ID,
	).Scan(&d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return d, nil
}

func (r *DeadlineRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM deadlines WHERE id = $1`, id)
	return err
}
