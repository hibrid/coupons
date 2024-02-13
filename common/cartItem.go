package common

import "errors"

type CartItem struct {
	SkuID                           string  `json:"skuId"`
	Quantity                        int     `json:"quantity"`
	UnitPrice                       float64 `json:"unitPrice"`
	DiscountAmountPerDiscountedUnit float64 `json:"discountAmountPerDiscountedUnit"`
	DiscountedUnitsQuantity         int     `json:"discountedUnitsQuantity"`
	TotalDiscountAmount             float64 `json:"getTotalDiscountAmount"`
	TotalGrossAmount                float64 `json:"totalGrossAmount"`
}

// Clone creates a copy of the CartItem instance, ensuring that modifications to the copy
// do not affect the original instance.
func (c *CartItem) Clone() *CartItem {
	return &CartItem{
		SkuID:                           c.SkuID,
		Quantity:                        c.Quantity,
		UnitPrice:                       c.UnitPrice,
		DiscountAmountPerDiscountedUnit: c.DiscountAmountPerDiscountedUnit,
		DiscountedUnitsQuantity:         c.DiscountedUnitsQuantity,
		TotalDiscountAmount:             c.TotalDiscountAmount,
		TotalGrossAmount:                c.TotalGrossAmount,
	}
}

// create the functions to work with CartItem
func (c *CartItem) GetSkuID() string {
	return c.SkuID
}

func (c *CartItem) SetSkuID(skuID string) {
	c.SkuID = skuID
}

func (c *CartItem) GetQuantity() int {
	return c.Quantity
}

func (c *CartItem) UpdateTotals() {
	c.GetNetTotalAmount()
}

func (c *CartItem) SetQuantity(quantity int) {
	c.Quantity = quantity
	c.UpdateTotals()
}

func (c *CartItem) GetUnitPrice() float64 {
	return c.UnitPrice
}

func (c *CartItem) SetUnitPrice(unitPrice float64) {
	c.UnitPrice = unitPrice
	c.UpdateTotals()
}

func (c *CartItem) GetGrossTotalAmount() float64 {
	if c.UnitPrice == 0 {
		return 0
	}
	if c.TotalGrossAmount != c.UnitPrice*float64(c.Quantity) {
		c.TotalGrossAmount = c.UnitPrice * float64(c.Quantity)
	}
	return c.TotalGrossAmount
}

func (c *CartItem) GetNetTotalAmount() float64 {
	return c.GetGrossTotalAmount() - c.GetTotalDiscountAmount()
}

func (c *CartItem) GetTotalDiscountAmount() float64 {
	if c.DiscountAmountPerDiscountedUnit > 0 {
		discountedAmount := c.DiscountAmountPerDiscountedUnit * float64(c.DiscountedUnitsQuantity)
		if discountedAmount > c.TotalGrossAmount {
			c.DiscountAmountPerDiscountedUnit = 0
		}
		c.DiscountAmountPerDiscountedUnit = discountedAmount
	}
	return c.DiscountAmountPerDiscountedUnit
}

func (c *CartItem) SetDiscountAmountPerUnit(discountAmount float64) {
	c.DiscountAmountPerDiscountedUnit = discountAmount
	c.UpdateTotals()
}

func (c *CartItem) GetDiscountedUnitsQuantity() int {
	return c.DiscountedUnitsQuantity
}

func (c *CartItem) SetDiscountedUnitQuantity(discountedUnitsQty int) {
	c.DiscountedUnitsQuantity = discountedUnitsQty
	c.UpdateTotals()
}

func (c *CartItem) IncrementCartItemQuantity(quantity int) int {
	c.Quantity += quantity
	c.UpdateTotals()
	return c.Quantity
}

func (c *CartItem) DecrementCartItemQuantity(quantity int) (int, error) {
	if c.Quantity-quantity < 0 {
		return 0, errors.New("quantity cannot be less than 0")
	}
	c.Quantity -= quantity
	c.UpdateTotals()
	return c.Quantity, nil
}
