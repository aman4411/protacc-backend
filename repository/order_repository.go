package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/utils"
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
		Status:          models.OrderStatusPendingBookingPayment,
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
		SELECT h.id, h.order_id, h.status, h.notes, h.created_by, h.created_at,
		       u.first_name, u.last_name, u.email
		FROM order_status_history h
		LEFT JOIN users u ON h.created_by = u.id
		WHERE h.order_id = $1
		ORDER BY h.created_at DESC`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderStatusHistory
	for rows.Next() {
		var entry models.OrderStatusHistory
		entry.User = &models.User{} // Initialize user object

		var userFirstName, userLastName, userEmail *string

		err := rows.Scan(
			&entry.ID,
			&entry.OrderID,
			&entry.Status,
			&entry.Notes,
			&entry.CreatedBy,
			&entry.CreatedAt,
			&userFirstName,
			&userLastName,
			&userEmail,
		)
		if err != nil {
			return nil, err
		}

		// Populate user information if available
		if userFirstName != nil && userLastName != nil && userEmail != nil {
			entry.User.FirstName = *userFirstName
			entry.User.LastName = *userLastName
			entry.User.Email = *userEmail
		} else {
			entry.User = nil // No user information (system-generated entry)
		}

		history = append(history, entry)
	}

	return history, nil
}

// GetOrderByRazorpayOrderID retrieves an order by Razorpay order ID
func (r *OrderRepository) GetOrderByRazorpayOrderID(ctx context.Context, razorpayOrderID string) (*models.Order, error) {
	query := `SELECT id, user_id, order_number, total_amount, booking_amount, remaining_amount, 
	          status, payment_status, razorpay_order_id, razorpay_payment_id, payment_method, 
	          payment_gateway, notes, created_at, updated_at 
	          FROM orders WHERE razorpay_order_id = $1`

	var order models.Order
	err := r.db.QueryRow(ctx, query, razorpayOrderID).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.TotalAmount,
		&order.BookingAmount,
		&order.RemainingAmount,
		&order.Status,
		&order.PaymentStatus,
		&order.RazorpayOrderID,
		&order.RazorpayPaymentID,
		&order.PaymentMethod,
		&order.PaymentGateway,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &order, nil
}

// UpdateOrderPaymentStatus updates the payment status and related fields of an order
func (r *OrderRepository) UpdateOrderPaymentStatus(ctx context.Context, orderID int, paymentStatus bool, razorpayPaymentID string, newStatus models.OrderStatus) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update order payment status and remaining amount if full payment is completed
	var updateQuery string
	var args []interface{}

	if newStatus == models.OrderStatusFullPaymentReceived {
		// When full payment is received, set remaining_amount to 0
		utils.LogInfo("Setting remaining_amount to 0 for full payment", "orderID", orderID, "status", newStatus)
		updateQuery = `UPDATE orders 
		               SET payment_status = $1, razorpay_payment_id = $2, status = $3, remaining_amount = 0, updated_at = CURRENT_TIMESTAMP 
		               WHERE id = $4`
		args = []interface{}{paymentStatus, razorpayPaymentID, newStatus, orderID}
	} else {
		// For other status updates, don't change remaining_amount
		utils.LogInfo("Updating order status without changing remaining_amount", "orderID", orderID, "status", newStatus)
		updateQuery = `UPDATE orders 
		               SET payment_status = $1, razorpay_payment_id = $2, status = $3, updated_at = CURRENT_TIMESTAMP 
		               WHERE id = $4`
		args = []interface{}{paymentStatus, razorpayPaymentID, newStatus, orderID}
	}

	_, err = tx.Exec(ctx, updateQuery, args...)
	if err != nil {
		return err
	}

	// Insert order status history
	historyQuery := `INSERT INTO order_status_history (order_id, status, notes, created_by) 
	                 VALUES ($1, $2, $3, $4)`

	notes := "Payment status updated via Razorpay"
	if paymentStatus {
		notes = fmt.Sprintf("Payment successful. Payment ID: %s", razorpayPaymentID)
	} else {
		notes = "Payment failed or cancelled"
	}

	_, err = tx.Exec(ctx, historyQuery, orderID, newStatus, notes, nil)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// UpdateOrderRazorpayOrderID updates the Razorpay order ID for an order
func (r *OrderRepository) UpdateOrderRazorpayOrderID(ctx context.Context, orderID int, razorpayOrderID string) error {
	query := `UPDATE orders SET razorpay_order_id = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(ctx, query, razorpayOrderID, orderID)
	return err
}

// GetOrderByID retrieves an order by ID for a specific user
func (r *OrderRepository) GetOrderByID(ctx context.Context, orderID int, userID string) (*models.Order, error) {
	query := `SELECT id, user_id, order_number, total_amount, booking_amount, remaining_amount, 
	          status, payment_status, razorpay_order_id, razorpay_payment_id, payment_method, 
	          payment_gateway, notes, created_at, updated_at 
	          FROM orders WHERE id = $1 AND user_id = $2`

	var order models.Order
	err := r.db.QueryRow(ctx, query, orderID, userID).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.TotalAmount,
		&order.BookingAmount,
		&order.RemainingAmount,
		&order.Status,
		&order.PaymentStatus,
		&order.RazorpayOrderID,
		&order.RazorpayPaymentID,
		&order.PaymentMethod,
		&order.PaymentGateway,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `SELECT oi.id, oi.order_id, oi.service_id, oi.quantity, oi.price, oi.booking_amount, oi.created_at,
	               s.name, s.short_description
	               FROM order_items oi
	               JOIN services s ON oi.service_id = s.id
	               WHERE oi.order_id = $1`

	rows, err := r.db.Query(ctx, itemsQuery, orderID)
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

	order.Items = items
	return &order, nil
}

// GetOrderByNumber retrieves an order by its order number for a specific user
func (r *OrderRepository) GetOrderByNumber(ctx context.Context, orderNumber string, userID string) (*models.Order, error) {
	query := `SELECT id, user_id, order_number, total_amount, booking_amount, remaining_amount, 
	          status, payment_status, razorpay_order_id, razorpay_payment_id, payment_method, 
	          payment_gateway, notes, created_at, updated_at 
	          FROM orders WHERE order_number = $1 AND user_id = $2::uuid`

	var order models.Order
	err := r.db.QueryRow(ctx, query, orderNumber, userID).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.TotalAmount,
		&order.BookingAmount,
		&order.RemainingAmount,
		&order.Status,
		&order.PaymentStatus,
		&order.RazorpayOrderID,
		&order.RazorpayPaymentID,
		&order.PaymentMethod,
		&order.PaymentGateway,
		&order.Notes,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get order items
	itemsQuery := `SELECT oi.id, oi.order_id, oi.service_id, oi.quantity, oi.price, oi.booking_amount, oi.created_at,
	               s.name, s.short_description
	               FROM order_items oi
	               JOIN services s ON oi.service_id = s.id
	               WHERE oi.order_id = $1`

	rows, err := r.db.Query(ctx, itemsQuery, order.ID)
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

	order.Items = items
	return &order, nil
}
