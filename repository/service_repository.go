package repository

import (
	"context"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Service Category Methods
func (r *ServiceRepository) CreateServiceCategory(ctx context.Context, category *models.ServiceCategory) error {
	query := `
		INSERT INTO service_categories (name, slug, description, icon, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
		category.Status,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)
}

func (r *ServiceRepository) GetServiceCategories(ctx context.Context) ([]models.ServiceCategory, error) {
	query := `
		SELECT id, name, slug, description, icon, status, created_at, updated_at
		FROM service_categories
		WHERE status = 'active'
		ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.ServiceCategory
	for rows.Next() {
		var category models.ServiceCategory
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Icon,
			&category.Status,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, category)
	}

	return categories, nil
}

// Service Methods
func (r *ServiceRepository) CreateService(ctx context.Context, service *models.Service) error {
	query := `
		INSERT INTO services (
			category_id, name, slug, description, short_description,
			features, requirements, price, booking_amount,
			estimated_delivery_days, icon, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		service.CategoryID,
		service.Name,
		service.Slug,
		service.Description,
		service.ShortDescription,
		service.Features,
		service.Requirements,
		service.Price,
		service.BookingAmount,
		service.EstimatedDeliveryDays,
		service.Icon,
		service.Status,
	).Scan(&service.ID, &service.CreatedAt, &service.UpdatedAt)
}

func (r *ServiceRepository) GetServices(ctx context.Context, categoryID *int) ([]models.Service, error) {
	query := `
		SELECT s.id, s.category_id, s.name, s.slug, s.description,
			s.short_description, s.features, s.requirements, s.price,
			s.booking_amount, s.estimated_delivery_days, s.icon, s.status,
			s.created_at, s.updated_at,
			c.id, c.name, c.slug, c.description, c.icon, c.status
		FROM services s
		JOIN service_categories c ON s.category_id = c.id
		WHERE s.status = 'active'
		AND ($1::int IS NULL OR s.category_id = $1)
		ORDER BY s.name`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var service models.Service
		service.Category = &models.ServiceCategory{}

		err := rows.Scan(
			&service.ID,
			&service.CategoryID,
			&service.Name,
			&service.Slug,
			&service.Description,
			&service.ShortDescription,
			&service.Features,
			&service.Requirements,
			&service.Price,
			&service.BookingAmount,
			&service.EstimatedDeliveryDays,
			&service.Icon,
			&service.Status,
			&service.CreatedAt,
			&service.UpdatedAt,
			&service.Category.ID,
			&service.Category.Name,
			&service.Category.Slug,
			&service.Category.Description,
			&service.Category.Icon,
			&service.Category.Status,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, service)
	}

	return services, nil
}

func (r *ServiceRepository) GetServiceByID(ctx context.Context, id int) (*models.Service, error) {
	query := `
		SELECT s.id, s.category_id, s.name, s.slug, s.description,
			s.short_description, s.features, s.requirements, s.price,
			s.booking_amount, s.estimated_delivery_days, s.icon, s.status,
			s.created_at, s.updated_at,
			c.id, c.name, c.slug, c.description, c.icon, c.status
		FROM services s
		JOIN service_categories c ON s.category_id = c.id
		WHERE s.id = $1 AND s.status = 'active'`

	service := &models.Service{Category: &models.ServiceCategory{}}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.CategoryID,
		&service.Name,
		&service.Slug,
		&service.Description,
		&service.ShortDescription,
		&service.Features,
		&service.Requirements,
		&service.Price,
		&service.BookingAmount,
		&service.EstimatedDeliveryDays,
		&service.Icon,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
		&service.Category.ID,
		&service.Category.Name,
		&service.Category.Slug,
		&service.Category.Description,
		&service.Category.Icon,
		&service.Category.Status,
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *ServiceRepository) GetServiceBySlug(ctx context.Context, slug string) (*models.Service, error) {
	query := `
		SELECT s.id, s.category_id, s.name, s.slug, s.description,
			s.short_description, s.features, s.requirements, s.price,
			s.booking_amount, s.estimated_delivery_days, s.icon, s.status,
			s.created_at, s.updated_at,
			c.id, c.name, c.slug, c.description, c.icon, c.status
		FROM services s
		JOIN service_categories c ON s.category_id = c.id
		WHERE s.slug = $1 AND s.status = 'active'`

	service := &models.Service{Category: &models.ServiceCategory{}}
	err := r.db.QueryRow(ctx, query, slug).Scan(
		&service.ID,
		&service.CategoryID,
		&service.Name,
		&service.Slug,
		&service.Description,
		&service.ShortDescription,
		&service.Features,
		&service.Requirements,
		&service.Price,
		&service.BookingAmount,
		&service.EstimatedDeliveryDays,
		&service.Icon,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
		&service.Category.ID,
		&service.Category.Name,
		&service.Category.Slug,
		&service.Category.Description,
		&service.Category.Icon,
		&service.Category.Status,
	)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// Cart Methods
func (r *ServiceRepository) AddToCart(ctx context.Context, userID string, serviceID int) error {
	query := `
		INSERT INTO cart_items (user_id, service_id)
		VALUES ($1, $2)
		ON CONFLICT (user_id, service_id) DO NOTHING`

	_, err := r.db.Exec(ctx, query, userID, serviceID)
	return err
}

func (r *ServiceRepository) GetCartItems(ctx context.Context, userID string) ([]models.CartItem, error) {
	query := `
		SELECT ci.id, ci.user_id, ci.service_id, ci.created_at,
			s.id, s.name, s.price, s.booking_amount
		FROM cart_items ci
		JOIN services s ON ci.service_id = s.id
		WHERE ci.user_id = $1`

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
			&item.Service.ID,
			&item.Service.Name,
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

func (r *ServiceRepository) RemoveFromCart(ctx context.Context, userID string, serviceID int) error {
	query := `DELETE FROM cart_items WHERE user_id = $1 AND service_id = $2`
	_, err := r.db.Exec(ctx, query, userID, serviceID)
	return err
}

// Order Methods
func (r *ServiceRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (
			user_id, service_id, order_number, total_amount,
			booking_amount, remaining_amount, status, payment_status, notes
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at`

	return r.db.QueryRow(ctx, query,
		order.UserID,
		order.ServiceID,
		order.OrderNumber,
		order.TotalAmount,
		order.BookingAmount,
		order.RemainingAmount,
		order.Status,
		order.PaymentStatus,
		order.Notes,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
}

func (r *ServiceRepository) GetOrders(ctx context.Context, userID *string) ([]models.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.service_id, o.order_number,
			o.total_amount, o.booking_amount, o.remaining_amount,
			o.status, o.payment_status, o.notes, o.created_at, o.updated_at,
			s.name, s.short_description,
			u.first_name, u.last_name, u.email
		FROM orders o
		JOIN services s ON o.service_id = s.id
		JOIN users u ON o.user_id = u.id
		WHERE ($1::uuid IS NULL OR o.user_id = $1)
		ORDER BY o.created_at DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		order.Service = &models.Service{}
		order.User = &models.User{}

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.ServiceID,
			&order.OrderNumber,
			&order.TotalAmount,
			&order.BookingAmount,
			&order.RemainingAmount,
			&order.Status,
			&order.PaymentStatus,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Service.Name,
			&order.Service.ShortDescription,
			&order.User.FirstName,
			&order.User.LastName,
			&order.User.Email,
		)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}

func (r *ServiceRepository) UpdateOrderStatus(ctx context.Context, orderID int, status models.OrderStatus, notes string, updatedBy string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update order status
	updateQuery := `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2`

	_, err = tx.Exec(ctx, updateQuery, status, orderID)
	if err != nil {
		return err
	}

	// Add status history
	historyQuery := `
		INSERT INTO order_status_history (order_id, status, notes, created_by)
		VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(ctx, historyQuery, orderID, status, notes, updatedBy)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *ServiceRepository) GetOrderStatusHistory(ctx context.Context, orderID int) ([]models.OrderStatusHistory, error) {
	query := `
		SELECT h.id, h.order_id, h.status, h.notes, h.created_by, h.created_at,
			u.first_name, u.last_name
		FROM order_status_history h
		JOIN users u ON h.created_by = u.id
		WHERE h.order_id = $1
		ORDER BY h.created_at DESC`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderStatusHistory
	for rows.Next() {
		var h models.OrderStatusHistory
		h.User = &models.User{}

		err := rows.Scan(
			&h.ID,
			&h.OrderID,
			&h.Status,
			&h.Notes,
			&h.CreatedBy,
			&h.CreatedAt,
			&h.User.FirstName,
			&h.User.LastName,
		)
		if err != nil {
			return nil, err
		}
		history = append(history, h)
	}

	return history, nil
}
