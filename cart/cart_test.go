package cart

import (
	"testing"

	"github.com/hibrid/coupons/common"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

// TestCart_GetItemNetTotal_ValidationError tests the scenario where GetNetTotalAmount results in an error.
func TestCart_GetItemNetTotal_ValidationError(t *testing.T) {
	// Assuming a mock or faulty CartItem that would fail during GetNetTotalAmount calculation.
	faultyItem := createFaultyCartItemForValidation()
	netTotal, err := new(Cart).GetItemNetTotal(faultyItem)
	assert.NotEqual(t, nil, err, "Expected an error when calculating net total of an invalid cart item")
	assert.Equal(t, float64(0), netTotal, "Expected net total to be 0 when there's a validation error")
}

// Since GetItemGrossTotal does not directly deal with errors, its test for validation error is not applicable as per the given function signature.

func TestCart_AddItem_ValidationError(t *testing.T) {
	cart := new(Cart)
	faultyItem := createFaultyCartItemForValidation()
	err := cart.AddItem(faultyItem)
	assert.NotEqual(t, nil, err, "Expected an error when adding an invalid cart item")
}

// Utility function to create a faulty CartItem for validation failure
func createFaultyCartItemForValidation() CartItem {
	// Create a CartItem with conditions that would fail validation, e.g., negative quantity.
	// Adjust according to your CartItem.Validate() implementation
	return CartItem{
		Quantity:  -1, // Assuming negative quantity would fail validation
		UnitPrice: decimal.NewFromFloat(10.0),
	}
}

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
				Subscription:                   SubscriptionInfo{BillingPeriodUnit: common.TimePeriodNoBilling},
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
		Subscription: SubscriptionInfo{BillingPeriodUnit: common.TimePeriodNoBilling},
	}
	err := cart.AddItem(newItem)
	assert.Nil(t, err, "Expected no error when adding a new item")
	assert.Equal(t, 2, len(cart.CartItems), "Expected 2 items in the cart after adding an item")

	// Test adding an item with an existing SKU
	existingItem := CartItem{SkuID: "test1", Quantity: 1, UnitPrice: decimal.NewFromFloat(10.0),
		Subscription: SubscriptionInfo{BillingPeriodUnit: common.TimePeriodNoBilling},
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

func TestCart_GetItemNetTotal(t *testing.T) {
	cart := newTestCart()
	netTotal, err := cart.GetItemNetTotal(cart.CartItems[0])
	assert.Nil(t, err, "Should not encounter an error calculating net total")
	expectedNetTotal := 19.0 // Assuming a single unit discount of 1.0 from a 10.0 unit price, times 2 units
	assert.Equal(t, expectedNetTotal, netTotal, "Net total calculation error")
}

func TestCart_GetItemGrossTotal(t *testing.T) {
	cart := newTestCart()
	grossTotal := cart.GetItemGrossTotal(cart.CartItems[0])

	expectedGrossTotal := 20.0 // 10.0 unit price times 2 units
	assert.Equal(t, expectedGrossTotal, grossTotal, "Gross total calculation error")
}

func TestCart_GetTotalGrossAmount(t *testing.T) {
	cart := newTestCart()
	totalAmount := cart.GetTotalGrossAmount()

	assert.Equal(t, float64(20), totalAmount, "Total amount calculation error")
}

func TestCart_GetTotalAmount(t *testing.T) {
	cart := newTestCart()
	totalAmount, err := cart.GetTotalNetAmount()
	assert.Nil(t, err, "Should not encounter an error calculating net total")
	assert.Equal(t, float64(19), totalAmount, "Total amount calculation error")
	cart.CartItems[0].UnitPrice = decimal.NewFromFloat(-10.0)
	_, err = cart.GetTotalNetAmount()
	assert.NotNil(t, err, "Should encounter an error calculating net total")
}

func TestCart_GetTotalDiscountAmount(t *testing.T) {
	cart := newTestCart()
	totalDiscountAmount := cart.GetTotalDiscountAmount()
	assert.Equal(t, 1.0, totalDiscountAmount, "Total discount amount calculation error")
}

func TestCart_GetTotalDiscountedUnits(t *testing.T) {
	cart := newTestCart()
	totalDiscountedUnits := cart.GetTotalDiscountedUnits()
	assert.Equal(t, int64(1), totalDiscountedUnits, "Total discounted units calculation error")
}

func TestCart_UpdateCartItem(t *testing.T) {
	cart := newTestCart()
	cart.UpdateCartItem("test1", 3, "15.0", "2.0", 2)
	//assert.Nil(t, err, "Updating cart item should not produce an error")
	updatedItem, _ := cart.GetCartItem("test1")
	assert.Equal(t, int64(3), updatedItem.Quantity, "Item quantity was not updated correctly")
	updatedPrice, _ := updatedItem.UnitPrice.Round(2).Float64()
	assert.Equal(t, float64(15.0), updatedPrice, "Item unit price was not updated correctly")
}

func TestCart_DecrementCartItem(t *testing.T) {
	cart := newTestCart()
	newQuantity, err := cart.DecrementCartItem("test1", 1)
	assert.Nil(t, err, "Decrementing cart item should not produce an error")
	assert.Equal(t, int64(1), newQuantity, "Item quantity should be decremented by 1")

	_, err = cart.DecrementCartItem("nonexistent", 1)
	assert.NotNil(t, err, "Decrementing a nonexistent item should produce an error")
}

func TestCart_IncrementCartItem(t *testing.T) {
	cart := newTestCart()
	newQuantity, err := cart.IncrementCartItem("test1", 1)
	assert.Nil(t, err, "Incrementing cart item should not produce an error")
	assert.Equal(t, int64(3), newQuantity, "Item quantity should be incremented by 1")

	_, err = cart.IncrementCartItem("nonexistent", 1)
	assert.NotNil(t, err, "Incrementing a nonexistent item should produce an error")
}

func TestCart_GetCartItem(t *testing.T) {
	cart := newTestCart()
	item, err := cart.GetCartItem("test1")
	assert.Nil(t, err, "Getting an existing cart item should not produce an error")
	assert.Equal(t, "test1", item.SkuID, "Retrieved item SKU ID mismatch")

	_, err = cart.GetCartItem("nonexistent")
	assert.NotNil(t, err, "Getting a nonexistent item should produce an error")
}

func TestCart_DoesItemExist(t *testing.T) {
	cart := newTestCart()
	_, exists := cart.DoesItemExist("test1")
	assert.True(t, exists, "Item should exist in the cart")

	_, exists = cart.DoesItemExist("nonexistent")
	assert.False(t, exists, "Item should not exist in the cart")
}

func TestCart_ClearCart(t *testing.T) {
	cart := newTestCart()
	cart.ClearCart()
	assert.Equal(t, 0, len(cart.CartItems), "Cart should be empty after clearing")
}

func TestCart_GetCartSize(t *testing.T) {
	cart := newTestCart()
	size := cart.GetCartSize()
	assert.Equal(t, 1, size, "Cart size should reflect the number of items in the cart")
}

func TestCart_IsEmpty(t *testing.T) {
	cart := newTestCart()
	isEmpty := cart.IsEmpty()
	assert.False(t, isEmpty, "Cart should not be empty with items in it")

	cart.ClearCart()
	isEmpty = cart.IsEmpty()
	assert.True(t, isEmpty, "Cart should be empty after items are cleared")
}
