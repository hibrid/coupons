package discount

import (
	"strconv"

	"github.com/hibrid/coupons/campaign"
	"github.com/hibrid/coupons/common"
	"github.com/hibrid/coupons/strategy"
)

type BOGODiscountStrategy struct {
	DiscountedSku string
	FreeSku       string
	DiscountedQty int
	FreeQty       int

	CampaignConfig *campaign.CampaignConfig
	// Additional fields for validation and result tracking
	isValid bool
	result  common.DiscountResult
}

func NewBogoDiscountStrategy(config *BOGODiscountStrategy) strategy.DiscountStrategy {
	return &BOGODiscountStrategy{
		DiscountedSku: config.DiscountedSku,
		FreeSku:       config.FreeSku,
		DiscountedQty: config.DiscountedQty,
		FreeQty:       config.FreeQty,
	}
}

func (b *BOGODiscountStrategy) SetCampaignConfig(config *campaign.CampaignConfig) {
	b.CampaignConfig = config

}

func (b *BOGODiscountStrategy) ValidateInputs() bool {
	// Placeholder validation logic
	b.isValid = true // Assuming validation passes for this example
	return b.isValid
}

func (b *BOGODiscountStrategy) ApplyDiscount() {
	if !b.isValid {
		return // Ensure discount logic is only executed if inputs are valid
	}
	// Apply the BOGO discount logic here
	// For example, adjusting quantities in a shopping cart

	// Populate the result with changes
	b.result = common.DiscountResult{
		OriginalValues: map[string]interface{}{"DiscountedQty": b.DiscountedQty, "FreeQty": 0},
		ModifiedValues: map[string]interface{}{"DiscountedQty": b.DiscountedQty, "FreeQty": b.FreeQty},
		Description:    "BOGO Applied: Buy " + strconv.Itoa(b.DiscountedQty) + " get " + strconv.Itoa(b.FreeQty) + " free.",
	}
}

func (b *BOGODiscountStrategy) GetDiscountResult() *common.DiscountResult {
	return &b.result
}
