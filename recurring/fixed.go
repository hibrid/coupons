package recurring

import "github.com/hibrid/coupons/common"

type FixedPriceSubscriptionConfig struct {
	FixedPriceDurationUnit   common.TimeUnit
	FixedPriceDurationLength int
	FixedPrice               float64
}

func (f *FixedPriceSubscriptionConfig) GetDurationUnit() common.TimeUnit {
	return f.FixedPriceDurationUnit
}

func (f *FixedPriceSubscriptionConfig) GetDurationLength() int {
	return f.FixedPriceDurationLength
}

func (f *FixedPriceSubscriptionConfig) GetPricing() float64 {
	return f.FixedPrice
}
