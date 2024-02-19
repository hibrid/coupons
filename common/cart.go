package common

import (
	"errors"

	"github.com/shopspring/decimal"
)

type Cart struct {
	CartItems []CartItem `json:"cartItems"`
}

// create the functions to work with Cart
func (c *Cart) GetCartItems() []CartItem {
	return c.CartItems
}

func (c *Cart) SetCartItems(cartItems []CartItem) error {
	if len(cartItems) == 0 {
		return errors.New("cart items cannot be empty")
	}
	c.CartItems = cartItems
	return nil
}

func (c *Cart) GetItemTotal(item CartItem) (float64, error) {
	netTotalAmountDecimal, err := item.GetNetTotalAmount() // Assuming GetNetTotalAmount returns decimal.Decimal
	if err != nil {
		// Handle the error, e.g., by logging or returning a default value
		return 0, err
	}
	// Round to two decimal places for cents precision
	netTotalAmountDecimal = netTotalAmountDecimal.Round(2)
	// Convert decimal.Decimal to float64
	netTotalAmountFloat64, _ := netTotalAmountDecimal.Float64()
	return netTotalAmountFloat64, nil
}

func (c *Cart) AddCartItem(cartItem CartItem) (*CartItem, error) {
	if err := cartItem.Validate(); err != nil {
		return nil, err
	}
	// check if the cart item already exists
	if item, exists := c.DoesItemExist(cartItem.SkuID); exists {
		// if it exists, update the quantity and total amount
		item.Quantity += cartItem.Quantity
		item.UpdateTotals()
		return item, nil
	}
	c.CartItems = append(c.CartItems, cartItem)
	return &cartItem, nil
}

func (c *Cart) RemoveCartItem(skuID string) error {
	for i, item := range c.CartItems {
		if item.SkuID == skuID {
			c.CartItems = append(c.CartItems[:i], c.CartItems[i+1:]...)
			return nil
		}
	}
	return errors.New("item not found")
}

func (c *Cart) GetTotalAmount() float64 {
	totalDecimal := decimal.Zero // Use decimal for precise summing
	for _, item := range c.CartItems {
		// Ensure TotalGrossAmount is a decimal.Decimal
		totalDecimal = totalDecimal.Add(item.TotalGrossAmount)
	}
	// Round to 2 decimal places for currency precision and convert to float64
	total, _ := totalDecimal.Round(2).Float64()
	return total
}

func (c *Cart) GetTotalDiscountAmount() float64 {
	totalDecimal := decimal.Zero // Use decimal for precise summing
	for _, item := range c.CartItems {
		// Ensure DiscountAmountPerDiscountedUnit is a decimal.Decimal
		totalDecimal = totalDecimal.Add(item.DiscountAmountPerDiscountedUnit)
	}
	// Round to 2 decimal places for currency precision and convert to float64
	total, _ := totalDecimal.Round(2).Float64()
	return total
}

func (c *Cart) GetTotalDiscountedUnits() int64 {
	var total int64
	for _, item := range c.CartItems {
		total += item.DiscountedUnitsQuantity
	}
	return total
}

func (c *Cart) UpdateCartItem(skuID string, quantity int64, unitPriceStr, discountAmountStr string, discountedUnits int64) {
	unitPrice, _ := decimal.NewFromString(unitPriceStr)
	discountAmount, _ := decimal.NewFromString(discountAmountStr)

	for i, item := range c.CartItems {
		if item.SkuID == skuID {
			c.CartItems[i].SetQuantity(quantity)
			c.CartItems[i].SetUnitPrice(unitPrice)                  // Adjusted to accept decimal.Decimal
			c.CartItems[i].SetDiscountAmountPerUnit(discountAmount) // Adjusted to accept decimal.Decimal
			c.CartItems[i].SetDiscountedUnitQuantity(discountedUnits)
			c.CartItems[i].UpdateTotals()
			break
		}
	}
}

func (c *Cart) DecrementCartItem(skuID string, quantity int64) (int64, error) {
	for i, item := range c.CartItems {
		if item.SkuID == skuID {
			return c.CartItems[i].DecrementCartItemQuantity(quantity)
		}
	}
	return 0, errors.New("item not found")
}

func (c *Cart) IncrementCartItem(skuID string, quantity int64) (int64, error) {
	for i, item := range c.CartItems {
		if item.SkuID == skuID {
			c.CartItems[i].IncrementCartItemQuantity(quantity)
			return c.CartItems[i].GetQuantity(), nil
		}
	}
	return 0, errors.New("item not found")
}

func (c *Cart) GetCartItem(skuID string) (*CartItem, error) {
	for _, item := range c.CartItems {
		if item.SkuID == skuID {
			return &item, nil
		}
	}
	return nil, errors.New("item not found")
}

func (c *Cart) DoesItemExist(skuID string) (*CartItem, bool) {
	for _, item := range c.CartItems {
		if item.SkuID == skuID {
			return &item, true
		}
	}
	return nil, false
}

func (c *Cart) ClearCart() {
	c.CartItems = []CartItem{}
}

func (c *Cart) GetCartSize() int {
	return len(c.CartItems)
}

func (c *Cart) IsEmpty() bool {
	return len(c.CartItems) == 0
}

/*
func (c *Cart) ValidateCartItemToAdd(cartItem CartItem) error {
	// First, validate the cart item itself
	if err := cartItem.Validate(); err != nil {
		return err
	}

	// Check for duplicate SKUs and validate against cart-specific rules
	for _, item := range c.CartItems {
		if item.SkuID == cartItem.SkuID {
			// Assuming you want to prevent adding duplicate SKUs with different details
			// Alternatively, you could check for and merge items with the same SKU here
			return fmt.Errorf("duplicate SKU ID '%s' found in cart", cartItem.SkuID)
		}
	}

	// Add any additional cart-specific validations here
	// For example, enforcing a maximum number of items in the cart
	if len(c.CartItems) >= 10 { // Example limit
		return errors.New("cannot add more than 10 items to the cart")
	}

	// Example: enforcing a maximum total discount amount or quantity
	totalDiscountedUnits := cartItem.DiscountedUnitsQuantity
	for _, item := range c.CartItems {
		totalDiscountedUnits += item.DiscountedUnitsQuantity
	}
	if totalDiscountedUnits > 50 { // Example limit
		return errors.New("adding this item exceeds the total allowed discounted units in the cart")
	}

	return nil
}
*/
