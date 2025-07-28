package repository

import (
	"context"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartRepository struct {
	db *pgxpool.Pool
}

func NewCartRepository(db *pgxpool.Pool) *CartRepository {
	return &CartRepository{db: db}
}

func (r *CartRepository) AddToCart(ctx context.Context, userID string, serviceID int) error {
	query := `
		INSERT INTO cart_items (user_id, service_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, service_id) DO NOTHING`

	_, err := r.db.Exec(ctx, query, userID, serviceID)
	return err
}

func (r *CartRepository) GetCartItems(ctx context.Context, userID string) ([]models.CartItem, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.service_id, ci.created_at,
			s.name, s.short_description, s.price, s.booking_amount
		FROM cart_items ci
		JOIN services s ON ci.service_id = s.id
		WHERE ci.user_id = $1
		ORDER BY ci.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		item.Service = &models.Service{}

		err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ServiceID,
			&item.CreatedAt,
			&item.Service.Name,
			&item.Service.ShortDescription,
			&item.Service.Price,
			&item.Service.BookingAmount,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *CartRepository) RemoveFromCart(ctx context.Context, userID string, serviceID int) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND service_id = $2`
	_, err := r.db.Exec(ctx, query, userID, serviceID)
	return err
}

func (r *CartRepository) ClearCart(ctx context.Context, userID string) error {
	query := `DELETE FROM cart_items WHERE user_id = $1`
	_, err := r.db.Exec(ctx, query, userID)
	return err
}
