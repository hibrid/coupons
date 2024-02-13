package loyalty

import "github.com/hibrid/coupons/common"

type LoyaltyConfig struct {
	RequirePreviousPurchase                    bool            `json:"requirePreviousPurchase"`
	RequiredPreviousPurchaseAmount             float64         `json:"requiredPreviousPurchaseAmount"`
	RequiredPreviousPurchaseAmountWithinUnit   common.TimeUnit `json:"requiredPreviousPurchaseAmountWithinUnit"`
	RequiredPreviousPurchaseAmountWithinLength int             `json:"requiredPreviousPurchaseAmountWithinLength"`
	RequirePreviousSubscription                bool            `json:"requirePreviousSubscription"`
	RequirePreviousSubscriptionWithinUnit      common.TimeUnit `json:"requirePreviousSubscriptionWithinUnit"`
	RequirePreviousSubscriptionWithinLength    int             `json:"requirePreviousSubscriptionWithinLength"`
	RequireCurrentSubscription                 bool            `json:"requireCurrentSubscription"`
	RequireCurrentSubscriptionLengthUnit       common.TimeUnit `json:"requireCurrentSubscriptionLengthUnit"`
	RequireCurrentSubscriptionLength           int             `json:"requireCurrentSubscriptionLength"`
}
