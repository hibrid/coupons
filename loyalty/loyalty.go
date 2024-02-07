package loyalty

import "github.com/hibrid/coupons/common"

type LoyaltyConfig struct {
	RequirePreviousPurchase                    bool                 `json:"requirePreviousPurchase"`
	RequiredPreviousPurchaseAmount             float64              `json:"requiredPreviousPurchaseAmount"`
	RequiredPreviousPurchaseAmountWithinUnit   common.BillingPeriod `json:"requiredPreviousPurchaseAmountWithinUnit"`
	RequiredPreviousPurchaseAmountWithinLength int                  `json:"requiredPreviousPurchaseAmountWithinLength"`
	RequirePreviousSubscription                bool                 `json:"requirePreviousSubscription"`
	RequirePreviousSubscriptionWithinUnit      common.BillingPeriod `json:"requirePreviousSubscriptionWithinUnit"`
	RequirePreviousSubscriptionWithinLength    int                  `json:"requirePreviousSubscriptionWithinLength"`
	RequireCurrentSubscription                 bool                 `json:"requireCurrentSubscription"`
	RequireCurrentSubscriptionLengthUnit       common.BillingPeriod `json:"requireCurrentSubscriptionLengthUnit"`
	RequireCurrentSubscriptionLength           int                  `json:"requireCurrentSubscriptionLength"`
}
