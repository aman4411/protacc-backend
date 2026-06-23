package models

import "time"

type CouponDiscountType string
type CouponApplicationMode string

const (
	CouponTypePercentage CouponDiscountType = "percentage"
	CouponTypeFlat       CouponDiscountType = "flat"

	// CouponModeFinal applies the whole discount to the final/remaining payment.
	CouponModeFinal CouponApplicationMode = "final"
	// CouponModeProportional splits the discount across booking and final payments.
	CouponModeProportional CouponApplicationMode = "proportional"
	// CouponModeBooking applies the discount to the booking payment first.
	CouponModeBooking CouponApplicationMode = "booking"
)

type Coupon struct {
	ID                int                   `json:"id"`
	Code              string                `json:"code"`
	Description       string                `json:"description"`
	DiscountType      CouponDiscountType    `json:"discount_type"`
	DiscountValue     float64               `json:"discount_value"`
	MaxDiscountAmount *float64              `json:"max_discount_amount,omitempty"`
	MinOrderAmount    float64               `json:"min_order_amount"`
	ApplicationMode   CouponApplicationMode `json:"application_mode"`
	UsageLimit        *int                  `json:"usage_limit,omitempty"`
	PerUserLimit      *int                  `json:"per_user_limit,omitempty"`
	ValidFrom         *time.Time            `json:"valid_from,omitempty"`
	ValidUntil        *time.Time            `json:"valid_until,omitempty"`
	IsActive          bool                  `json:"is_active"`
	IsVisible         bool                  `json:"is_visible"`
	ShowOnHomepage    bool                  `json:"show_on_homepage"`

	// Targeting: empty means the coupon applies to everything. Otherwise it
	// applies only to items in these categories and/or services.
	ApplicableCategoryIDs []int32   `json:"applicable_category_ids"`
	ApplicableServiceIDs  []int32   `json:"applicable_service_ids"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`

	// Derived (admin display): how many times the coupon has been redeemed on paid orders.
	UsedCount int `json:"used_count"`
}

// CartLine is one item in a user's cart, used for coupon eligibility checks.
type CartLine struct {
	ServiceID     int
	CategoryID    int
	Price         float64
	BookingAmount float64
}

// Restricted reports whether the coupon targets specific categories/services.
func (c *Coupon) Restricted() bool {
	return len(c.ApplicableCategoryIDs) > 0 || len(c.ApplicableServiceIDs) > 0
}

// AppliesTo reports whether a cart line is eligible for this coupon.
func (c *Coupon) AppliesTo(line CartLine) bool {
	if !c.Restricted() {
		return true
	}
	for _, id := range c.ApplicableServiceIDs {
		if int(id) == line.ServiceID {
			return true
		}
	}
	for _, id := range c.ApplicableCategoryIDs {
		if int(id) == line.CategoryID {
			return true
		}
	}
	return false
}

// EligibleSubtotal sums the prices of cart lines this coupon applies to.
func (c *Coupon) EligibleSubtotal(lines []CartLine) float64 {
	total := 0.0
	for _, l := range lines {
		if c.AppliesTo(l) {
			total += l.Price
		}
	}
	return total
}

// DiscountedAmounts is the result of applying a coupon to a cart, split across
// the two-stage (booking + final) payment.
type DiscountedAmounts struct {
	Subtotal        float64 `json:"subtotal"`
	DiscountAmount  float64 `json:"discount_amount"`
	TotalAmount     float64 `json:"total_amount"`     // subtotal - discount
	BookingAmount   float64 `json:"booking_amount"`   // due now
	RemainingAmount float64 `json:"remaining_amount"` // due at final stage
}

// round2 rounds to 2 decimal places (paise).
func round2(v float64) float64 {
	return float64(int64(v*100+0.5)) / 100
}

// ComputeDiscountedAmounts applies a discount of `discount` to a cart with the
// given subtotal and base booking amount, splitting it across the booking and
// final payments according to the application mode. It guarantees non-negative
// booking/remaining and that total == subtotal - discount.
func ComputeDiscountedAmounts(subtotal, bookingBase, discount float64, mode CouponApplicationMode) DiscountedAmounts {
	if discount < 0 {
		discount = 0
	}
	if discount > subtotal {
		discount = subtotal
	}

	total := round2(subtotal - discount)
	var booking, remaining float64

	switch mode {
	case CouponModeBooking:
		// Discount eats the booking first; any overflow reduces the final payment.
		booking = bookingBase - discount
		if booking < 0 {
			booking = 0
		}
		booking = round2(booking)
		remaining = round2(total - booking)
	case CouponModeProportional:
		share := 0.0
		if subtotal > 0 {
			share = bookingBase / subtotal
		}
		booking = round2(bookingBase - discount*share)
		if booking < 0 {
			booking = 0
		}
		remaining = round2(total - booking)
	default: // CouponModeFinal
		// Booking unchanged; discount comes off the final payment. Cap so the
		// remaining never goes negative.
		if discount > subtotal-bookingBase {
			discount = round2(subtotal - bookingBase)
			if discount < 0 {
				discount = 0
			}
			total = round2(subtotal - discount)
		}
		booking = round2(bookingBase)
		remaining = round2(total - booking)
	}

	if remaining < 0 {
		remaining = 0
	}
	return DiscountedAmounts{
		Subtotal:        round2(subtotal),
		DiscountAmount:  round2(subtotal - total),
		TotalAmount:     total,
		BookingAmount:   booking,
		RemainingAmount: remaining,
	}
}
