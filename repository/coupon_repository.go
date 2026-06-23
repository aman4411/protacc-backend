package repository

import (
	"context"
	"errors"

	"github.com/aman4411/protacc-backend/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CouponRepository struct {
	db *pgxpool.Pool
}

func NewCouponRepository(db *pgxpool.Pool) *CouponRepository {
	return &CouponRepository{db: db}
}

// Orders that count as a redemption: booking paid or beyond (not unpaid/cancelled).
const redeemedOrderFilter = `status NOT IN ('pending_payment', 'pending_booking_payment', 'cancelled')`

func scanCoupon(row pgx.Row, c *models.Coupon) error {
	return row.Scan(
		&c.ID, &c.Code, &c.Description, &c.DiscountType, &c.DiscountValue,
		&c.MaxDiscountAmount, &c.MinOrderAmount, &c.ApplicationMode,
		&c.UsageLimit, &c.PerUserLimit, &c.ValidFrom, &c.ValidUntil,
		&c.IsActive, &c.IsVisible, &c.ShowOnHomepage,
		&c.ApplicableCategoryIDs, &c.ApplicableServiceIDs, &c.CreatedAt, &c.UpdatedAt,
	)
}

const couponCols = `id, code, description, discount_type, discount_value, max_discount_amount,
	min_order_amount, application_mode, usage_limit, per_user_limit, valid_from, valid_until,
	is_active, is_visible, show_on_homepage, applicable_category_ids, applicable_service_ids, created_at, updated_at`

func (r *CouponRepository) Create(ctx context.Context, c *models.Coupon) (*models.Coupon, error) {
	query := `
		INSERT INTO coupons (code, description, discount_type, discount_value, max_discount_amount,
			min_order_amount, application_mode, usage_limit, per_user_limit, valid_from, valid_until, is_active, is_visible, show_on_homepage, applicable_category_ids, applicable_service_ids)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRow(ctx, query,
		c.Code, c.Description, c.DiscountType, c.DiscountValue, c.MaxDiscountAmount,
		c.MinOrderAmount, c.ApplicationMode, c.UsageLimit, c.PerUserLimit, c.ValidFrom, c.ValidUntil, c.IsActive, c.IsVisible, c.ShowOnHomepage,
		c.ApplicableCategoryIDs, c.ApplicableServiceIDs,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CouponRepository) Update(ctx context.Context, c *models.Coupon) (*models.Coupon, error) {
	query := `
		UPDATE coupons SET code=$1, description=$2, discount_type=$3, discount_value=$4,
			max_discount_amount=$5, min_order_amount=$6, application_mode=$7, usage_limit=$8,
			per_user_limit=$9, valid_from=$10, valid_until=$11, is_active=$12, is_visible=$13, show_on_homepage=$14,
			applicable_category_ids=$15, applicable_service_ids=$16, updated_at=NOW()
		WHERE id=$17
		RETURNING created_at, updated_at`
	err := r.db.QueryRow(ctx, query,
		c.Code, c.Description, c.DiscountType, c.DiscountValue, c.MaxDiscountAmount,
		c.MinOrderAmount, c.ApplicationMode, c.UsageLimit, c.PerUserLimit, c.ValidFrom, c.ValidUntil, c.IsActive, c.IsVisible, c.ShowOnHomepage,
		c.ApplicableCategoryIDs, c.ApplicableServiceIDs, c.ID,
	).Scan(&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (r *CouponRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx, `DELETE FROM coupons WHERE id = $1`, id)
	return err
}

func (r *CouponRepository) GetByCode(ctx context.Context, code string) (*models.Coupon, error) {
	var c models.Coupon
	err := scanCoupon(r.db.QueryRow(ctx, `SELECT `+couponCols+` FROM coupons WHERE code = $1`, code), &c)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

// List returns all coupons with a derived used_count (paid redemptions).
func (r *CouponRepository) List(ctx context.Context) ([]models.Coupon, error) {
	query := `
		SELECT ` + couponCols + `,
			COALESCE((SELECT COUNT(*) FROM orders o WHERE o.coupon_code = c.code AND o.` + redeemedOrderFilter + `), 0) AS used_count
		FROM coupons c
		ORDER BY c.created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	coupons := []models.Coupon{}
	for rows.Next() {
		var c models.Coupon
		if err := rows.Scan(
			&c.ID, &c.Code, &c.Description, &c.DiscountType, &c.DiscountValue,
			&c.MaxDiscountAmount, &c.MinOrderAmount, &c.ApplicationMode,
			&c.UsageLimit, &c.PerUserLimit, &c.ValidFrom, &c.ValidUntil,
			&c.IsActive, &c.IsVisible, &c.ShowOnHomepage,
			&c.ApplicableCategoryIDs, &c.ApplicableServiceIDs, &c.CreatedAt, &c.UpdatedAt, &c.UsedCount,
		); err != nil {
			return nil, err
		}
		coupons = append(coupons, c)
	}
	return coupons, nil
}

// ListVisible returns active, in-window coupons flagged to show on the cart.
func (r *CouponRepository) ListVisible(ctx context.Context) ([]models.Coupon, error) {
	return r.listPublic(ctx, "is_visible")
}

// ListForHomepage returns active, in-window coupons flagged for the homepage banner.
func (r *CouponRepository) ListForHomepage(ctx context.Context) ([]models.Coupon, error) {
	return r.listPublic(ctx, "show_on_homepage")
}

// listPublic returns public-safe coupon fields for an active, in-window coupon
// matching the given boolean flag column.
func (r *CouponRepository) listPublic(ctx context.Context, flagColumn string) ([]models.Coupon, error) {
	query := `
		SELECT code, description, discount_type, discount_value, max_discount_amount, min_order_amount
		FROM coupons
		WHERE is_active = true AND ` + flagColumn + ` = true
		  AND (valid_from IS NULL OR valid_from <= NOW())
		  AND (valid_until IS NULL OR valid_until >= NOW())
		ORDER BY created_at DESC`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	coupons := []models.Coupon{}
	for rows.Next() {
		var c models.Coupon
		if err := rows.Scan(&c.Code, &c.Description, &c.DiscountType, &c.DiscountValue, &c.MaxDiscountAmount, &c.MinOrderAmount); err != nil {
			return nil, err
		}
		coupons = append(coupons, c)
	}
	return coupons, nil
}

// CountRedemptions returns total paid redemptions of a coupon.
func (r *CouponRepository) CountRedemptions(ctx context.Context, code string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE coupon_code = $1 AND `+redeemedOrderFilter, code).Scan(&n)
	return n, err
}

// CountUserRedemptions returns paid redemptions of a coupon by one user.
func (r *CouponRepository) CountUserRedemptions(ctx context.Context, code, userID string) (int, error) {
	var n int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM orders WHERE coupon_code = $1 AND user_id = $2 AND `+redeemedOrderFilter, code, userID).Scan(&n)
	return n, err
}
