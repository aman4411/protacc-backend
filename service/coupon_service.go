package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aman4411/protacc-backend/cache"
	"github.com/aman4411/protacc-backend/models"
	"github.com/aman4411/protacc-backend/repository"
)

// Cache prefix + TTL for the public coupon banners. Kept short because coupons
// are time-sensitive. Validate() and admin List() are never cached.
const (
	couponCachePrefix = "coupons:"
	couponCacheTTL    = 5 * time.Minute
)

type CouponService struct {
	repo  *repository.CouponRepository
	cache *cache.Cache
}

func NewCouponService(repo *repository.CouponRepository, c *cache.Cache) *CouponService {
	return &CouponService{repo: repo, cache: c}
}

func (s *CouponService) invalidate() {
	if s.cache != nil {
		s.cache.InvalidatePrefix(couponCachePrefix)
	}
}

func normalizeCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

// Validate checks a coupon against a user + cart and returns the coupon and the
// raw discount amount (before the booking/final split). The discount is computed
// on the subtotal of cart items the coupon applies to.
func (s *CouponService) Validate(ctx context.Context, code, userID string, lines []models.CartLine) (*models.Coupon, float64, error) {
	c, err := s.repo.GetByCode(ctx, normalizeCode(code))
	if err != nil {
		return nil, 0, err
	}
	if c == nil {
		return nil, 0, fmt.Errorf("invalid coupon code")
	}
	if !c.IsActive {
		return nil, 0, fmt.Errorf("this coupon is no longer active")
	}

	now := time.Now()
	if c.ValidFrom != nil && now.Before(*c.ValidFrom) {
		return nil, 0, fmt.Errorf("this coupon is not active yet")
	}
	if c.ValidUntil != nil && now.After(*c.ValidUntil) {
		return nil, 0, fmt.Errorf("this coupon has expired")
	}

	// Subtotal of the items this coupon applies to (whole cart if unrestricted).
	subtotal := c.EligibleSubtotal(lines)
	if c.Restricted() && subtotal <= 0 {
		return nil, 0, fmt.Errorf("this coupon does not apply to the items in your cart")
	}
	if subtotal < c.MinOrderAmount {
		return nil, 0, fmt.Errorf("this coupon requires a minimum eligible order of ₹%.0f", c.MinOrderAmount)
	}

	if c.UsageLimit != nil {
		used, err := s.repo.CountRedemptions(ctx, c.Code)
		if err != nil {
			return nil, 0, err
		}
		if used >= *c.UsageLimit {
			return nil, 0, fmt.Errorf("this coupon has reached its usage limit")
		}
	}
	if userID != "" && c.PerUserLimit != nil {
		used, err := s.repo.CountUserRedemptions(ctx, c.Code, userID)
		if err != nil {
			return nil, 0, err
		}
		if used >= *c.PerUserLimit {
			return nil, 0, fmt.Errorf("you have already used this coupon")
		}
	}

	discount := 0.0
	switch c.DiscountType {
	case models.CouponTypePercentage:
		discount = subtotal * c.DiscountValue / 100
		if c.MaxDiscountAmount != nil && discount > *c.MaxDiscountAmount {
			discount = *c.MaxDiscountAmount
		}
	case models.CouponTypeFlat:
		discount = c.DiscountValue
	}
	if discount > subtotal {
		discount = subtotal
	}
	return c, discount, nil
}

// CRUD passthroughs for admin.
func (s *CouponService) List(ctx context.Context) ([]models.Coupon, error) {
	return s.repo.List(ctx)
}

// ListVisible returns coupons flagged to show publicly on the cart.
func (s *CouponService) ListVisible(ctx context.Context) ([]models.Coupon, error) {
	return cache.Load(s.cache, couponCachePrefix+"visible", couponCacheTTL, func() ([]models.Coupon, error) {
		return s.repo.ListVisible(ctx)
	})
}

// ListForHomepage returns coupons flagged for the homepage campaign banner.
func (s *CouponService) ListForHomepage(ctx context.Context) ([]models.Coupon, error) {
	return cache.Load(s.cache, couponCachePrefix+"homepage", couponCacheTTL, func() ([]models.Coupon, error) {
		return s.repo.ListForHomepage(ctx)
	})
}

func (s *CouponService) Create(ctx context.Context, c *models.Coupon) (*models.Coupon, error) {
	if err := s.normalizeAndValidate(c); err != nil {
		return nil, err
	}
	res, err := s.repo.Create(ctx, c)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *CouponService) Update(ctx context.Context, c *models.Coupon) (*models.Coupon, error) {
	if err := s.normalizeAndValidate(c); err != nil {
		return nil, err
	}
	res, err := s.repo.Update(ctx, c)
	if err == nil {
		s.invalidate()
	}
	return res, err
}

func (s *CouponService) Delete(ctx context.Context, id int) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.invalidate()
	}
	return err
}

func (s *CouponService) normalizeAndValidate(c *models.Coupon) error {
	c.Code = normalizeCode(c.Code)
	if c.Code == "" {
		return fmt.Errorf("coupon code is required")
	}
	if c.DiscountType != models.CouponTypePercentage && c.DiscountType != models.CouponTypeFlat {
		return fmt.Errorf("invalid discount type")
	}
	if c.DiscountValue <= 0 {
		return fmt.Errorf("discount value must be greater than 0")
	}
	if c.DiscountType == models.CouponTypePercentage && c.DiscountValue > 100 {
		return fmt.Errorf("percentage discount cannot exceed 100")
	}
	if c.ApplicationMode == "" {
		c.ApplicationMode = models.CouponModeFinal
	}
	if c.ApplicationMode != models.CouponModeFinal &&
		c.ApplicationMode != models.CouponModeProportional &&
		c.ApplicationMode != models.CouponModeBooking {
		return fmt.Errorf("invalid application mode")
	}
	return nil
}
