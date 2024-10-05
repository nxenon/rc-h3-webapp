package models

type CouponCode struct {
	CouponId        int
	CouponValue     string
	IsValid         int // 0 or 1
	DiscountPercent int
}

type ApplyCouponRequest struct {
	CouponValue string `json:"couponValue"`
}
