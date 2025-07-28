package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *models.Order) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert order
	orderQuery := `
		INSERT INTO orders (user_id, order_number, total_amount, booking_amount, remaining_amount, status, payment_status, notes)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, orderQuery,
		order.UserID,
		order.OrderNumber,
		order.TotalAmount,
		order.BookingAmount,
		order.RemainingAmount,
		order.Status,
		order.PaymentStatus,
		order.Notes,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return err
	}

	// Insert order items
	for i := range order.Items {
		item := &order.Items[i]
		itemQuery := `
			INSERT INTO order_items (order_id, service_id, quantity, price, booking_amount)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

		err = tx.QueryRow(ctx, itemQuery,
			order.ID,
			item.ServiceID,
			item.Quantity,
			item.Price,
			item.BookingAmount,
		).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return err
		}
		item.OrderID = order.ID
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) CreateOrderFromCart(ctx context.Context, userID string) (*models.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Get cart items with service details
	cartQuery := `
		SELECT ci.service_id, s.price, s.booking_amount, s.name
		FROM cart_items ci
		JOIN services s ON ci.service_id = s.id
		WHERE ci.user_id = $1 AND s.status = 'active'`

	rows, err := tx.Query(ctx, cartQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []models.OrderItem
	var totalAmount, totalBookingAmount float64

	for rows.Next() {
		var serviceID int
		var price, bookingAmount float64
		var serviceName string

		err := rows.Scan(&serviceID, &price, &bookingAmount, &serviceName)
		if err != nil {
			return nil, err
		}

		orderItem := models.OrderItem{
			ServiceID:     serviceID,
			Quantity:      1,
			Price:         price,
			BookingAmount: bookingAmount,
		}

		orderItems = append(orderItems, orderItem)
		totalAmount += price
		totalBookingAmount += bookingAmount
	}

	if len(orderItems) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Generate order number
	orderNumber := fmt.Sprintf("ORD%s", time.Now().Format("20060102150405"))

	// Create order
	order := &models.Order{
		UserID:          userID,
		OrderNumber:     orderNumber,
		TotalAmount:     totalAmount,
		BookingAmount:   totalBookingAmount,
		RemainingAmount: totalAmount - totalBookingAmount,
		Status:          models.OrderStatusPendingPayment,
		PaymentStatus:   false,
		Items:           orderItems,
	}

	// Insert order
	orderQuery := `
		INSERT INTO orders (user_id, order_number, total_amount, booking_amount, remaining_amount, status, payment_status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at`

	err = tx.QueryRow(ctx, orderQuery,
		order.UserID,
		order.OrderNumber,
		order.TotalAmount,
		order.BookingAmount,
		order.RemainingAmount,
		order.Status,
		order.PaymentStatus,
	).Scan(&order.ID, &order.CreatedAt, &order.UpdatedAt)
	if err != nil {
		return nil, err
	}

	// Insert order items
	for i := range order.Items {
		item := &order.Items[i]
		itemQuery := `
			INSERT INTO order_items (order_id, service_id, quantity, price, booking_amount)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, created_at`

		err = tx.QueryRow(ctx, itemQuery,
			order.ID,
			item.ServiceID,
			item.Quantity,
			item.Price,
			item.BookingAmount,
		).Scan(&item.ID, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		item.OrderID = order.ID
	}

	// Clear cart
	clearCartQuery := `DELETE FROM cart_items WHERE user_id = $1`
	_, err = tx.Exec(ctx, clearCartQuery, userID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) GetOrders(ctx context.Context, userID *string) ([]models.Order, error) {
	query := `
		SELECT o.id, o.user_id, o.order_number,
			o.total_amount, o.booking_amount, o.remaining_amount,
			o.status, o.payment_status, o.notes, o.created_at, o.updated_at,
			u.first_name, u.last_name, u.email
		FROM orders o
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
		order.User = &models.User{}

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.TotalAmount,
			&order.BookingAmount,
			&order.RemainingAmount,
			&order.Status,
			&order.PaymentStatus,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.User.FirstName,
			&order.User.LastName,
			&order.User.Email,
		)
		if err != nil {
			return nil, err
		}

		// Load order items
		items, err := r.GetOrderItems(ctx, order.ID)
		if err != nil {
			return nil, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, nil
}

func (r *OrderRepository) GetOrdersWithFilters(ctx context.Context, userID *string, page, limit int, status, search string) ([]models.Order, int, error) {
	offset := (page - 1) * limit

	// Build dynamic query
	baseQuery := `
		FROM orders o
		JOIN users u ON o.user_id = u.id
		WHERE 1=1`

	countQuery := "SELECT COUNT(*) " + baseQuery
	selectQuery := `
		SELECT o.id, o.user_id, o.order_number,
			o.total_amount, o.booking_amount, o.remaining_amount,
			o.status, o.payment_status, o.notes, o.created_at, o.updated_at,
			u.first_name, u.last_name, u.email ` + baseQuery

	args := []interface{}{}
	argIndex := 1

	// Add user filter
	if userID != nil {
		baseQuery += fmt.Sprintf(" AND o.user_id = $%d", argIndex)
		args = append(args, *userID)
		argIndex++
	}

	// Add status filter
	if status != "" {
		baseQuery += fmt.Sprintf(" AND o.status = $%d", argIndex)
		args = append(args, status)
		argIndex++
	}

	// Add search filter (order number, user name, email)
	if search != "" {
		baseQuery += fmt.Sprintf(" AND (LOWER(o.order_number) LIKE LOWER($%d) OR LOWER(u.first_name) LIKE LOWER($%d) OR LOWER(u.last_name) LIKE LOWER($%d) OR LOWER(u.email) LIKE LOWER($%d))", argIndex, argIndex, argIndex, argIndex)
		searchTerm := "%" + search + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	// Update queries with filters
	countQuery = "SELECT COUNT(*) " + baseQuery
	selectQuery = `
		SELECT o.id, o.user_id, o.order_number,
			o.total_amount, o.booking_amount, o.remaining_amount,
			o.status, o.payment_status, o.notes, o.created_at, o.updated_at,
			u.first_name, u.last_name, u.email ` + baseQuery +
		fmt.Sprintf(" ORDER BY o.created_at DESC LIMIT $%d OFFSET $%d", argIndex, argIndex+1)

	args = append(args, limit, offset)

	// Get total count
	var total int
	err := r.db.QueryRow(ctx, countQuery, args[:len(args)-2]...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get orders
	rows, err := r.db.Query(ctx, selectQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		order.User = &models.User{}

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.TotalAmount,
			&order.BookingAmount,
			&order.RemainingAmount,
			&order.Status,
			&order.PaymentStatus,
			&order.Notes,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.User.FirstName,
			&order.User.LastName,
			&order.User.Email,
		)
		if err != nil {
			return nil, 0, err
		}

		// Load order items
		items, err := r.GetOrderItems(ctx, order.ID)
		if err != nil {
			return nil, 0, err
		}
		order.Items = items

		orders = append(orders, order)
	}

	return orders, total, nil
}

func (r *OrderRepository) GetOrderItems(ctx context.Context, orderID int) ([]models.OrderItem, error) {
	query := `
		SELECT oi.id, oi.order_id, oi.service_id, oi.quantity, oi.price, oi.booking_amount, oi.created_at,
			s.name, s.short_description
		FROM order_items oi
		JOIN services s ON oi.service_id = s.id
		WHERE oi.order_id = $1
		ORDER BY oi.created_at`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []models.OrderItem
	for rows.Next() {
		var item models.OrderItem
		item.Service = &models.Service{}

		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ServiceID,
			&item.Quantity,
			&item.Price,
			&item.BookingAmount,
			&item.CreatedAt,
			&item.Service.Name,
			&item.Service.ShortDescription,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	return items, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, orderID int, status models.OrderStatus, notes string, updatedBy string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update order status
	updateQuery := `
		UPDATE orders 
		SET status = $1, notes = $2, updated_at = NOW()
		WHERE id = $3`

	_, err = tx.Exec(ctx, updateQuery, status, notes, orderID)
	if err != nil {
		return err
	}

	// Insert status history
	historyQuery := `
		INSERT INTO order_status_history (order_id, status, notes, created_by)
		VALUES ($1, $2, $3, $4)`

	_, err = tx.Exec(ctx, historyQuery, orderID, status, notes, updatedBy)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *OrderRepository) GetOrderStatusHistory(ctx context.Context, orderID int) ([]models.OrderStatusHistory, error) {
	query := `
		SELECT id, order_id, status, notes, created_by, created_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderStatusHistory
	for rows.Next() {
		var entry models.OrderStatusHistory

		err := rows.Scan(
			&entry.ID,
			&entry.OrderID,
			&entry.Status,
			&entry.Notes,
			&entry.CreatedBy,
			&entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		history = append(history, entry)
	}

	return history, nil
}
