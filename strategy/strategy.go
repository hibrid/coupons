package strategy

import (
	"github.com/hibrid/coupons/common"
)

type CampaignData interface {
	GetCampaignDetails() (startDate, endDate string, isActive bool)
	// Add other necessary methods to expose campaign data to strategies
}

type RecurringDiscountStrategy interface {
	GetDurationUnit() common.BillingPeriod
	GetDurationLength() int
	GetPricing() float64
}

// DiscountStrategy encapsulates the process of validating discount conditions,
// applying discounts, and retrieving the results of the discount application.
type DiscountStrategy interface {
	//SetCampaignConfig(*campaign.CampaignConfig) // Sets the campaign configuration
	ValidateInputs(data CampaignData) bool     // Validates the inputs needed for the discount strategy
	ApplyDiscount(data CampaignData)           // Executes the discount logic
	GetDiscountResult() *common.DiscountResult // Retrieves the result of the discount application
}
