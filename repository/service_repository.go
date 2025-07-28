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
		INSERT INTO service_categories (name, slug, description, icon, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
		category.Status,
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
		SET name = $1, slug = $2, description = $3, icon = $4, status = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Slug,
		category.Description,
		category.Icon,
		category.Status,
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

// SearchServices searches for services by name, description, or category
func (r *ServiceRepository) SearchServices(ctx context.Context, query string) ([]models.Service, error) {
	searchQuery := `
		SELECT 
			s.id, s.name, s.slug, s.description, s.price, s.estimated_delivery_days,
			s.features, s.requirements, s.category_id, s.status, s.created_at, s.updated_at,
			sc.name as category_name, sc.slug as category_slug
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
			s.name`

	searchTerm := "%" + query + "%"
	exactTerm := query + "%"

	rows, err := r.db.Query(ctx, searchQuery, searchTerm, exactTerm)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var services []models.Service
	for rows.Next() {
		var service models.Service
		var categoryName, categorySlug sql.NullString

		err := rows.Scan(
			&service.ID,
			&service.Name,
			&service.Slug,
			&service.Description,
			&service.Price,
			&service.EstimatedDeliveryDays,
			&service.Features,
			&service.Requirements,
			&service.CategoryID,
			&service.Status,
			&service.CreatedAt,
			&service.UpdatedAt,
			&categoryName,
			&categorySlug,
		)
		if err != nil {
			return nil, err
		}

		// Populate category information if available
		if categoryName.Valid && categorySlug.Valid {
			service.Category = &models.ServiceCategory{
				ID:   service.CategoryID,
				Name: categoryName.String,
				Slug: categorySlug.String,
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
			estimated_delivery_days, icon, status
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
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
		service.EstimatedDeliveryDays,
		service.Icon,
		service.Status,
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
		SET category_id = $1, name = $2, slug = $3, description = $4, short_description = $5,
			features = $6, requirements = $7, price = $8, booking_amount = $9,
			estimated_delivery_days = $10, icon = $11, status = $12, updated_at = NOW()
		WHERE id = $13
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
		service.EstimatedDeliveryDays,
		service.Icon,
		service.Status,
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
			s.short_description, s.features, s.requirements, s.price,
			s.booking_amount, s.estimated_delivery_days, s.icon, s.status,
			s.created_at, s.updated_at,
			c.id, c.name, c.slug, c.description, c.icon, c.status
		FROM services s
		JOIN service_categories c ON s.category_id = c.id
		WHERE s.status = 'active'
		AND ($1::int IS NULL OR s.category_id = $1)
		AND ($2::text = '' OR c.slug = $2)
		ORDER BY s.name`

	rows, err := r.db.Query(ctx, query, categoryID, categorySlug)
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
