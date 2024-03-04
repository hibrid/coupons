package common

import (
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a new cart with predefined items
func newTestCart() *Cart {
	return &Cart{
		CartItems: []CartItem{
			{
				SkuID:                          "test1",
				Quantity:                       2,
				UnitPrice:                      decimal.NewFromFloat(10.0),
				DiscountDescription:            "Test Discount",
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(1.0),
				NumberOfUnitsDiscounted:        1,
				Subscription:                   SubscriptionInfo{BillingPeriodUnit: TimePeriodNoBilling},
			},
		},
	}
}

func TestCart_GetCartItems(t *testing.T) {
	cart := newTestCart()
	items := cart.GetCartItems()
	assert.Equal(t, 1, len(items), "Expected 1 item in the cart")
}

func TestCart_SetCartItems(t *testing.T) {
	cart := &Cart{}
	err := cart.SetCartItems([]CartItem{})
	assert.NotNil(t, err, "Expected an error when setting empty cart items")

	newItems := []CartItem{
		{SkuID: "test2", Quantity: 1, UnitPrice: decimal.NewFromFloat(20.0)},
	}
	err = cart.SetCartItems(newItems)
	assert.Nil(t, err, "Expected no error when setting non-empty cart items")
	assert.Equal(t, 1, len(cart.CartItems), "Expected 1 item in the cart after setting items")
}

func TestCart_AddItem(t *testing.T) {
	cart := newTestCart()
	newItem := CartItem{
		SkuID: "test2", Quantity: 1, UnitPrice: decimal.NewFromFloat(20.0),
		Subscription: SubscriptionInfo{BillingPeriodUnit: TimePeriodNoBilling},
	}
	err := cart.AddItem(newItem)
	assert.Nil(t, err, "Expected no error when adding a new item")
	assert.Equal(t, 2, len(cart.CartItems), "Expected 2 items in the cart after adding an item")

	// Test adding an item with an existing SKU
	existingItem := CartItem{SkuID: "test1", Quantity: 1, UnitPrice: decimal.NewFromFloat(10.0),
		Subscription: SubscriptionInfo{BillingPeriodUnit: TimePeriodNoBilling},
	}
	err = cart.AddItem(existingItem)
	assert.Nil(t, err, "Expected no error when adding an existing item")
	assert.Equal(t, int64(3), cart.CartItems[0].Quantity, "Expected quantity of the first item to be 3 after adding existing item")
}

func TestCart_RemoveItem(t *testing.T) {
	cart := newTestCart()
	err := cart.RemoveItem("test1")
	assert.Nil(t, err, "Expected no error when removing an existing item")
	assert.Equal(t, 0, len(cart.CartItems), "Expected 0 items in the cart after removing an item")

	err = cart.RemoveItem("nonexistent")
	assert.NotNil(t, err, "Expected an error when removing a nonexistent item")
}

// Continue with similar structure for other Cart methods:
// - TestCart_GetItemNetTotal
// - TestCart_GetItemGrossTotal
// - TestCart_GetTotalAmount
// - TestCart_GetTotalDiscountAmount
// - TestCart_GetTotalDiscountedUnits
// - TestCart_UpdateCartItem
// - TestCart_DecrementCartItem
// - TestCart_IncrementCartItem
// - TestCart_GetCartItem
// - TestCart_DoesItemExist
// - TestCart_ClearCart
// - TestCart_GetCartSize
// - TestCart_IsEmpty

// Each of these tests should verify the correct behavior of the method under test,
// checking for expected outcomes and handling of edge cases.
