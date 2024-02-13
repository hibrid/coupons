package recurring

import "github.com/hibrid/coupons/common"

type TrialPeriodDiscountConfig struct {
	TrialPeriodUnit   common.TimeUnit
	TrialPeriodLength int
	PostTrialPricing  float64
}

func (t *TrialPeriodDiscountConfig) GetDurationUnit() common.TimeUnit {
	return t.TrialPeriodUnit
}

func (t *TrialPeriodDiscountConfig) GetDurationLength() int {
	return t.TrialPeriodLength
}

func (t *TrialPeriodDiscountConfig) GetPricing() float64 {
	return t.PostTrialPricing
}
