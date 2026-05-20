package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderDocumentRepository struct {
	db *pgxpool.Pool
}

func NewOrderDocumentRepository(db *pgxpool.Pool) *OrderDocumentRepository {
	return &OrderDocumentRepository{db: db}
}

func (r *OrderDocumentRepository) Create(ctx context.Context, doc *models.OrderDocument) error {
	query := `
		INSERT INTO order_documents (order_id, uploaded_by, document_type, title, drive_url, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	var notes *string
	if doc.Notes != nil && *doc.Notes != "" {
		notes = doc.Notes
	}

	return r.db.QueryRow(ctx, query,
		doc.OrderID,
		doc.UploadedBy,
		doc.DocumentType,
		doc.Title,
		doc.DriveURL,
		notes,
	).Scan(&doc.ID, &doc.CreatedAt, &doc.UpdatedAt)
}

func (r *OrderDocumentRepository) GetByOrderID(ctx context.Context, orderID int) ([]models.OrderDocument, error) {
	query := `
		SELECT od.id, od.order_id, od.uploaded_by, od.document_type, od.title, od.drive_url, od.notes,
			od.created_at, od.updated_at,
			u.id, u.first_name, u.last_name, u.email, u.phone, u.is_email_verified, u.role, u.created_at, u.updated_at
		FROM order_documents od
		JOIN users u ON od.uploaded_by = u.id
		WHERE od.order_id = $1
		ORDER BY od.created_at DESC`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("error fetching order documents: %v", err)
	}
	defer rows.Close()

	var docs []models.OrderDocument
	for rows.Next() {
		var doc models.OrderDocument
		var uploader models.User
		var notes *string

		err := rows.Scan(
			&doc.ID,
			&doc.OrderID,
			&doc.UploadedBy,
			&doc.DocumentType,
			&doc.Title,
			&doc.DriveURL,
			&notes,
			&doc.CreatedAt,
			&doc.UpdatedAt,
			&uploader.ID,
			&uploader.FirstName,
			&uploader.LastName,
			&uploader.Email,
			&uploader.Phone,
			&uploader.IsEmailVerified,
			&uploader.Role,
			&uploader.CreatedAt,
			&uploader.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning order document: %v", err)
		}

		doc.Notes = notes
		doc.Uploader = &models.UserResponse{
			ID:              uploader.ID,
			FirstName:       uploader.FirstName,
			LastName:        uploader.LastName,
			Email:           uploader.Email,
			Phone:           uploader.Phone,
			IsEmailVerified: uploader.IsEmailVerified,
			Role:            uploader.Role,
			CreatedAt:       uploader.CreatedAt,
			UpdatedAt:       uploader.UpdatedAt,
		}
		docs = append(docs, doc)
	}

	return docs, nil
}

func (r *OrderDocumentRepository) GetOrderOwnerID(ctx context.Context, orderID int) (string, error) {
	var userID string
	err := r.db.QueryRow(ctx, `SELECT user_id::text FROM orders WHERE id = $1`, orderID).Scan(&userID)
	if err == pgx.ErrNoRows {
		return "", errors.New("order not found")
	}
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (r *OrderDocumentRepository) GetOrderStatus(ctx context.Context, orderID int) (models.OrderStatus, error) {
	var status models.OrderStatus
	err := r.db.QueryRow(ctx, `SELECT status FROM orders WHERE id = $1`, orderID).Scan(&status)
	if err == pgx.ErrNoRows {
		return "", errors.New("order not found")
	}
	return status, err
}
