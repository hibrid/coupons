package recurring

import "github.com/hibrid/coupons/common"

type RecurringDiscountConfig struct {
	DiscountDurationUnit   common.BillingPeriod
	DiscountDurationLength int
	DiscountValue          float64
}

func (r *RecurringDiscountConfig) GetDurationUnit() common.BillingPeriod {
	return r.DiscountDurationUnit
}

func (r *RecurringDiscountConfig) GetDurationLength() int {
	return r.DiscountDurationLength
}

func (r *RecurringDiscountConfig) GetPricing() float64 {
	return r.DiscountValue
}
