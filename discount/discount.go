package discount

import "github.com/hibrid/coupons/common"

type DiscountConfig struct {
	Type            common.CouponDiscountType `json:"type"`
	Value           float64                   `json:"value"`
	MinimumPurchase float64                   `json:"minimumPurchase,omitempty"`
	DiscountedSku   string                    `json:"discountedSku,omitempty"`
	FreeSku         string                    `json:"freeSku,omitempty"`
	DiscountedQty   int                       `json:"discountedQty,omitempty"`
	FreeQty         int                       `json:"freeQty,omitempty"`
}
