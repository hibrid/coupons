package common

import "errors"

type Cart struct {
	CartItems []CartItem `json:"cartItems"`
}

// create the functions to work with Cart
func (c *Cart) GetCartItems() []CartItem {
	return c.CartItems
}

func (c *Cart) SetCartItems(cartItems []CartItem) {
	c.CartItems = cartItems
}

func (c *Cart) GetItemTotal(item CartItem) float64 {

	return item.GetNetTotalAmount()
}

func (c *Cart) AddCartItem(cartItem CartItem) *CartItem {
	// check if the cart item already exists
	if item, exists := c.DoesItemExist(cartItem.SkuID); exists {
		// if it exists, update the quantity and total amount
		item.Quantity += cartItem.Quantity
		item.UpdateTotals()
		return item
	}
	c.CartItems = append(c.CartItems, cartItem)
	return &cartItem
}

func (c *Cart) RemoveCartItem(skuID string) {
	for i, item := range c.CartItems {
		if item.SkuID == skuID {
			c.CartItems = append(c.CartItems[:i], c.CartItems[i+1:]...)
			break
		}
	}
}

func (c *Cart) GetTotalAmount() float64 {
	var total float64
	for _, item := range c.CartItems {
		total += item.TotalGrossAmount
	}
	return total
}

func (c *Cart) GetTotalDiscountAmount() float64 {
	var total float64
	for _, item := range c.CartItems {
		total += item.DiscountAmountPerDiscountedUnit
	}
	return total
}

func (c *Cart) GetTotalDiscountedUnits() int {
	var total int
	for _, item := range c.CartItems {
		total += item.DiscountedUnitsQuantity
	}
	return total
}

func (c *Cart) UpdateCartItem(skuID string, quantity int, unitPrice, discountAmount float64, discountedUnits int) {
	for i, item := range c.CartItems {
		if item.SkuID == skuID {
			c.CartItems[i].SetQuantity(quantity)
			c.CartItems[i].SetUnitPrice(unitPrice)
			c.CartItems[i].SetDiscountAmountPerUnit(discountAmount)
			c.CartItems[i].SetDiscountedUnitQuantity(discountedUnits)
			c.CartItems[i].UpdateTotals()
			break
		}
	}
}

func (c *Cart) DecrementCartItem(skuID string, quantity int) (int, error) {
	for i, item := range c.CartItems {
		if item.SkuID == skuID {
			return c.CartItems[i].DecrementCartItemQuantity(quantity)
		}
	}
	return 0, errors.New("item not found")
}

func (c *Cart) IncrementCartItem(skuID string, quantity int) (int, error) {
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
