package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aman4411/protacc-backend/models"
	"github.com/gosimple/slug"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ServiceRepository struct {
	db *pgxpool.Pool
}

func NewServiceRepository(db *pgxpool.Pool) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// Service Category Methods
func (r *ServiceRepository) CreateServiceCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {
	// Generate slug from name if not provided
	if category.Slug == "" {
		category.Slug = slug.Make(category.Name)
	}

	query := `
		INSERT INTO service_categories (name, slug, description, icon, status, priority)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
		category.Status,
		category.Priority,
	).Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *ServiceRepository) UpdateServiceCategory(ctx context.Context, category *models.ServiceCategory) (*models.ServiceCategory, error) {
	// Generate slug from name if not provided
	if category.Slug == "" {
		category.Slug = slug.Make(category.Name)
	}

	query := `
		UPDATE service_categories 
		SET name = $1, slug = $2, description = $3, icon = $4, status = $5, priority = $6, updated_at = NOW()
		WHERE id = $7
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
		category.Status,
		category.Priority,
		category.ID,
	).Scan(&category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (r *ServiceRepository) DeleteServiceCategory(ctx context.Context, categoryID int) error {
	query := `DELETE FROM service_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, categoryID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *ServiceRepository) GetServiceCategories(ctx context.Context) ([]models.ServiceCategory, error) {
	query := `
		SELECT id, name, slug, description, icon, status, priority, created_at, updated_at
		FROM service_categories
		WHERE status = 'active'
		ORDER BY priority ASC, name ASC`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := []models.ServiceCategory{}
	for rows.Next() {
		var category models.ServiceCategory
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Slug,
			&category.Description,
			&category.Icon,
			&category.Status,
			&category.Priority,
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

// SearchServices searches for services by name, description, or category
func (r *ServiceRepository) SearchServices(ctx context.Context, query string) ([]models.Service, error) {
	searchQuery := `
		SELECT 
			s.id, s.name, s.slug, s.description, s.price, s.min_delivery_days, s.max_delivery_days,
			s.features, s.requirements, s.suited_for, s.whats_included, s.category_id, s.status, s.priority, s.created_at, s.updated_at,
			sc.name as category_name, sc.slug as category_slug, sc.priority as category_priority
		FROM services s
		LEFT JOIN service_categories sc ON s.category_id = sc.id
		WHERE s.status = 'active' 
		AND (
			LOWER(s.name) LIKE LOWER($1) OR 
			LOWER(s.description) LIKE LOWER($1) OR
			LOWER(sc.name) LIKE LOWER($1)
		)
		ORDER BY 
			CASE 
				WHEN LOWER(s.name) LIKE LOWER($2) THEN 1
				WHEN LOWER(s.name) LIKE LOWER($1) THEN 2
				WHEN LOWER(s.description) LIKE LOWER($1) THEN 3
				ELSE 4
			END,
			s.priority ASC,
			sc.priority ASC,
			s.name`

	searchTerm := "%" + query + "%"
	exactTerm := query + "%"

	rows, err := r.db.Query(ctx, searchQuery, searchTerm, exactTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []models.Service{}
	for rows.Next() {
		var service models.Service
		var categoryName, categorySlug sql.NullString
		var categoryPriority sql.NullInt32

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Slug,
			&service.Description,
			&service.Price,
			&service.MinDeliveryDays,
			&service.MaxDeliveryDays,
			&service.Features,
			&service.Requirements,
			&service.SuitedFor,
			&service.WhatsIncluded,
			&service.CategoryID,
			&service.Status,
			&service.Priority,
			&service.CreatedAt,
			&service.UpdatedAt,
			&categoryName,
			&categorySlug,
			&categoryPriority,
		)
		if err != nil {
			return nil, err
		}

		// Populate category information if available
		if categoryName.Valid && categorySlug.Valid {
			service.Category = &models.ServiceCategory{
				ID:       service.CategoryID,
				Name:     categoryName.String,
				Slug:     categorySlug.String,
				Priority: int(categoryPriority.Int32),
			}
		}

		services = append(services, service)
	}

	return services, nil
}

// Service Methods
func (r *ServiceRepository) CreateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	// Generate slug from name if not provided
	if service.Slug == "" {
		service.GenerateSlug()
	}

	query := `
		INSERT INTO services (
			category_id, name, slug, description, short_description,
			features, requirements, price, booking_amount,
			min_delivery_days, max_delivery_days, icon, status, priority, suited_for, whats_included
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		service.CategoryID,
		service.Name,
		service.Slug,
		service.Description,
		service.ShortDescription,
		service.Features,
		service.Requirements,
		service.Price,
		service.BookingAmount,
		service.MinDeliveryDays,
		service.MaxDeliveryDays,
		service.Icon,
		service.Status,
		service.Priority,
		service.SuitedFor,
		service.WhatsIncluded,
	).Scan(&service.ID, &service.CreatedAt, &service.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *ServiceRepository) UpdateService(ctx context.Context, service *models.Service) (*models.Service, error) {
	// Generate slug from name if not provided
	if service.Slug == "" {
		service.GenerateSlug()
	}

	query := `
		UPDATE services 
		SET category_id = $1, name = $2, slug = $3::text, description = $4, short_description = $5,
			features = $6, requirements = $7, price = $8, booking_amount = $9,
			min_delivery_days = $10, max_delivery_days = $11, icon = $12, status = $13, priority = $14,
			suited_for = $15, whats_included = $16,
			previous_slugs = (
				SELECT ARRAY(
					SELECT DISTINCT s2
					FROM unnest(
						COALESCE(previous_slugs, ARRAY[]::text[])
						|| CASE WHEN slug::text <> $3::text THEN ARRAY[slug::text] ELSE ARRAY[]::text[] END
					) AS s2
					WHERE s2 <> $3::text AND s2 <> ''
				)
			),
			updated_at = NOW()
		WHERE id = $17
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		service.CategoryID,
		service.Name,
		service.Slug,
		service.Description,
		service.ShortDescription,
		service.Features,
		service.Requirements,
		service.Price,
		service.BookingAmount,
		service.MinDeliveryDays,
		service.MaxDeliveryDays,
		service.Icon,
		service.Status,
		service.Priority,
		service.SuitedFor,
		service.WhatsIncluded,
		service.ID,
	).Scan(&service.CreatedAt, &service.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return service, nil
}

func (r *ServiceRepository) DeleteService(ctx context.Context, serviceID int) error {
	query := `DELETE FROM services WHERE id = $1`

	result, err := r.db.Exec(ctx, query, serviceID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}

	return nil
}

func (r *ServiceRepository) GetServices(ctx context.Context, categoryID *int, categorySlug string) ([]models.Service, error) {
	query := `
		SELECT s.id, s.category_id, s.name, s.slug, s.description,
			s.short_description, s.features, s.requirements, s.suited_for, s.whats_included, s.price,
			s.booking_amount, s.min_delivery_days, s.max_delivery_days, s.icon, s.status, s.priority,
			s.created_at, s.updated_at,
			c.id, c.name, c.slug, c.description, c.icon, c.status, c.priority,
			COALESCE((SELECT AVG(rating) FROM reviews WHERE service_id = s.id AND status = 'published'), 0) AS avg_rating,
			COALESCE((SELECT COUNT(*) FROM reviews WHERE service_id = s.id AND status = 'published'), 0) AS review_count
		FROM services s
		JOIN service_categories c ON s.category_id = c.id
		WHERE s.status = 'active'
		AND ($1::int IS NULL OR s.category_id = $1)
		AND ($2::text = '' OR c.slug = $2)
		ORDER BY s.priority ASC, c.priority ASC, s.name ASC`

	rows, err := r.db.Query(ctx, query, categoryID, categorySlug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	services := []models.Service{}
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
			&service.SuitedFor,
			&service.WhatsIncluded,
			&service.Price,
			&service.BookingAmount,
			&service.MinDeliveryDays,
			&service.MaxDeliveryDays,
			&service.Icon,
			&service.Status,
			&service.Priority,
			&service.CreatedAt,
			&service.UpdatedAt,
			&service.Category.ID,
			&service.Category.Name,
			&service.Category.Slug,
			&service.Category.Description,
			&service.Category.Icon,
			&service.Category.Status,
			&service.Category.Priority,
			&service.AvgRating,
			&service.ReviewCount,
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
			s.short_description, s.features, s.requirements, s.suited_for, s.whats_included, s.price,
			s.booking_amount, s.min_delivery_days, s.max_delivery_days, s.icon, s.status,
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
		&service.SuitedFor,
		&service.WhatsIncluded,
		&service.Price,
		&service.BookingAmount,
		&service.MinDeliveryDays,
		&service.MaxDeliveryDays,
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
			s.short_description, s.features, s.requirements, s.suited_for, s.whats_included, s.price,
			s.booking_amount, s.min_delivery_days, s.max_delivery_days, s.icon, s.status,
			s.created_at, s.updated_at, s.previous_slugs,
			c.id, c.name, c.slug, c.description, c.icon, c.status
		FROM services s
		JOIN service_categories c ON s.category_id = c.id
		WHERE (s.slug = $1 OR $1 = ANY(s.previous_slugs)) AND s.status = 'active'`

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
		&service.SuitedFor,
		&service.WhatsIncluded,
		&service.Price,
		&service.BookingAmount,
		&service.MinDeliveryDays,
		&service.MaxDeliveryDays,
		&service.Icon,
		&service.Status,
		&service.CreatedAt,
		&service.UpdatedAt,
		&service.PreviousSlugs,
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

// Priority Management Methods
func (r *ServiceRepository) UpdateCategoryPriority(ctx context.Context, categoryID int, priority int) error {
	query := `UPDATE service_categories SET priority = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, priority, categoryID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *ServiceRepository) UpdateServicePriority(ctx context.Context, serviceID int, priority int) error {
	query := `UPDATE services SET priority = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.Exec(ctx, query, priority, serviceID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("service not found")
	}

	return nil
}
