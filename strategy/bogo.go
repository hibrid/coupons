package strategy

import (
	"fmt"

	"github.com/hibrid/coupons/common"
)

type BOGODiscountStrategy struct {
	DiscountedSku string
	FreeSku       string
	DiscountedQty int
	FreeQty       int

	// Additional fields for validation and result tracking
	isValid bool
	result  common.DiscountResult
}

func NewBogoDiscountStrategy(config *BOGODiscountStrategy) DiscountStrategy {
	return &BOGODiscountStrategy{
		DiscountedSku: config.DiscountedSku,
		FreeSku:       config.FreeSku,
		DiscountedQty: config.DiscountedQty,
		FreeQty:       config.FreeQty,
	}
}

func (b *BOGODiscountStrategy) GetName() string {
	return "BOGO"
}

/*
func (b *BOGODiscountStrategy) SetCampaignConfig(config campaign.CampaignConfig) {
	b.CampaignConfig = &config

}
*/
// Adjust BOGODiscountStrategy methods to use CampaignData interface
func (b *BOGODiscountStrategy) ValidateInputs(data CampaignData) bool {
	// Use data through the CampaignData interface
	_, _, isActive := data.GetCampaignDetails()
	b.isValid = isActive // Simplified example logic
	return b.isValid
}

func (b *BOGODiscountStrategy) ApplyDiscount(data CampaignData) {
	if !b.isValid {
		// If the inputs are not valid, do not proceed with applying the discount.
		return
	}

	// Example logic for applying the discount:
	// 1. Check if the discounted product is in the cart and has enough quantity.
	// 2. For every 'DiscountedQty', add 'FreeQty' of 'FreeSku' to the cart.
	// Note: This is simplified logic; real-world scenarios require handling inventory, pricing, etc.

	cart2 := data.GetCart()

	if discountedQty, exists := cart2.DoesItemExist(b.DiscountedSku); exists && discountedQty.GetQuantity() >= int64(b.DiscountedQty) {
		//numDiscounts := discountedQty.GetQuantity() / int64(b.DiscountedQty)
		//totalFreeItems := numDiscounts * int64(b.FreeQty)
		originalCartItem, err := cart2.GetCartItem(b.DiscountedSku)
		if err != nil {
			// Handle the error, e.g., by logging or setting the result accordingly.
			return
		}
		originalClone := originalCartItem.Clone()
		// Update the cart with the free items, this is pseudo-code and assumes an AddToCart method exists
		// TODO: Need to get the unit price from the cart or product catalog
		/*
			cartItem := cart.CartItem{
				SkuID:     b.FreeSku,
				Quantity:  totalFreeItems,
				UnitPrice: decimal.NewFromInt(0),
				//DiscountAmountPerDiscountedUnit: 0, // maybe we need to set this to the unit price of the free item
				//DiscountedUnitsQuantity:         0,
			}
			cartItem.UpdateTotals()
			cart2.AddItem(cartItem)
		*/
		// Populate the result with changes
		b.result = common.DiscountResult{
			OriginalValues: map[string]interface{}{"DiscountedQty": originalClone.GetDiscountedUnitsQuantity(), "FreeQty": 0},
			//ModifiedValues: map[string]interface{}{"DiscountedQty": cartItem.GetDiscountedUnitsQuantity(), "FreeQty": totalFreeItems},
			Description: fmt.Sprintf("BOGO Applied: Buy %d get %d free.", b.DiscountedQty, b.FreeQty),
		}
	} else {
		// If conditions are not met, you might want to handle this case, e.g., by logging or setting the result accordingly.
	}
}

func (b *BOGODiscountStrategy) GetDiscountResult() *common.DiscountResult {
	return &b.result
}
