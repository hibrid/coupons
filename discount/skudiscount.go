package discount

import "github.com/hibrid/coupons/common"

type SkuDiscountConfig struct {
	DiscountedSku string  // SKU of the item to which the discount applies
	DiscountValue float64 // The discount value, which could be a percentage or a fixed amount
	DiscountType  common.CouponDiscountType
}
