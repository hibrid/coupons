package recurring

import "github.com/hibrid/coupons/common"

type FixedPriceSubscriptionConfig struct {
	FixedPriceDurationUnit   common.BillingPeriod
	FixedPriceDurationLength int
	FixedPrice               float64
}

func (f *FixedPriceSubscriptionConfig) GetDurationUnit() common.BillingPeriod {
	return f.FixedPriceDurationUnit
}

func (f *FixedPriceSubscriptionConfig) GetDurationLength() int {
	return f.FixedPriceDurationLength
}

func (f *FixedPriceSubscriptionConfig) GetPricing() float64 {
	return f.FixedPrice
}
