package referral

import "github.com/hibrid/coupons/common"

type ReferralConfig struct {
	GiveReferrerDiscount bool                      `json:"giveReferrerDiscount"`
	ReferrerDiscountType common.CouponDiscountType `json:"referrerDiscountType"`
	ReferrerDiscount     float64                   `json:"referrerDiscount"`
	ReferrerSku          string                    `json:"referrerSku"`
	GiveRefereeDiscount  bool                      `json:"giveRefereeDiscount"`
	RefereeDiscountType  common.CouponDiscountType `json:"refereeDiscountType"`
	RefereeDiscount      float64                   `json:"refereeDiscount"`
	RefereeSku           string                    `json:"refereeSku"`
}
