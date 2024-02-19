package common

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert" // Using testify for easier assertions
)

func TestCalculateTrialPeriodDiscount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	cartItem := &CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       3,
			TrialPeriodUnit:   TimePeriodMonthly,
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	expectedDiscount := decimal.NewFromInt(3).Mul(unitPrice) // 3 months * $10
	discount := cartItem.calculateTotalTrialPeriodDiscount()
	assert.True(t, expectedDiscount.Equal(discount))
}

func TestNormalizeDuration(t *testing.T) {
	tests := []struct {
		name        string
		duration    int64
		fromUnit    TimeUnit
		toUnit      TimeUnit
		expected    float64
		expectError bool
	}{
		{
			name:        "Hours to Days",
			duration:    36, // 36 hours
			fromUnit:    TimePeriodHourly,
			toUnit:      TimePeriodDaily,
			expected:    1.5, // Ratio is 1.5 days
			expectError: false,
		},
		{
			name:        "Weeks to Days",
			duration:    2, // 2 weeks
			fromUnit:    TimePeriodWeekly,
			toUnit:      TimePeriodDaily,
			expected:    14, // 2 weeks * 7 days/week
			expectError: false,
		},
		{
			name:        "Days to Weeks",
			duration:    10, // 10 days
			fromUnit:    TimePeriodDaily,
			toUnit:      TimePeriodWeekly,
			expected:    1.428571, // Rounded up from 1.42857 weeks
			expectError: false,
		},
		{
			name:        "Monthly to BiAnnual",
			duration:    3, // 3 months
			fromUnit:    TimePeriodMonthly,
			toUnit:      TimePeriodBiAnnual,
			expected:    0.493151, // Rounded up from 0.5 of a half-year
			expectError: false,
		},
		{
			name:        "Invalid FromUnit",
			duration:    1,
			fromUnit:    TimePeriodUnknown,
			toUnit:      TimePeriodWeekly,
			expected:    0,
			expectError: true,
		},
		{
			name:        "Invalid ToUnit",
			duration:    1,
			fromUnit:    TimePeriodDaily,
			toUnit:      TimePeriodUnknown,
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NormalizeDuration(tt.duration, tt.fromUnit, tt.toUnit)
			if tt.expectError {
				if err == nil {
					t.Errorf("%s expected an error but got none", tt.name)
				}
			} else {
				if err != nil {
					t.Errorf("%s unexpected error: %v", tt.name, err)
				}
				if !decimal.NewFromFloat(result).Round(2).Equal(decimal.NewFromFloat(tt.expected).Round(2)) {
					t.Errorf("%s expected %f, got %f", tt.name, tt.expected, result)
				}
			}
		})
	}
}

func TestCalculatePhaseDiscountPercentInvalidConfig(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	cartItem := CartItem{
		UnitPrice: unitPrice,
	}
	phase := DiscountPhase{
		Duration:      2,          // 2 billing cycles
		DiscountValue: 50,         // 50% discount
		DiscountType:  Percentage, // Indicating the DiscountRate is a percentage
	}

	_, err := cartItem.calculateDiscountForPhase(phase)
	assert.Error(t, err, "Expected error when calculating a percentage discount")

}

func TestCalculatePhaseDiscountPercent(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:      2, // 2 * DurationUnit
		DurationUnit:  TimePeriodDaily,
		DiscountValue: 50,         // 50% discount of the cost of two days
		DiscountType:  Percentage, // Indicating the DiscountRate is a percentage
		Application:   Recurring,
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	expectedDiscount := decimal.NewFromFloat(.33) // 2 cycles * $10 * 50%
	discount, err := cartItem.calculateDiscountForPhase(phase)
	printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
	expectedDiscount = decimal.NewFromFloat(.33) // 2 cycles * $10 * 50%
	discount, err = cartItem.calculateTotalDiscountValue(phase)
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	expectedDiscountAsFloat, _ = expectedDiscount.Round(2).Float64()
	discountAsFloat, _ = discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculatePhaseDiscountPercentRecurringValid(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        2, // 2 * DurationUnit
		DurationUnit:                    TimePeriodMonthly,
		DiscountValue:                   50,         // 50% discount of the cost of two months
		DiscountType:                    Percentage, // Indicating the DiscountRate is a percentage
		Application:                     Recurring,
		ApplicableNumberOfBillingCycles: 3,
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	expectedDiscount := decimal.NewFromFloat(10) // 2 cycles * $10 * 50%
	discount, err := cartItem.calculateDiscountForPhase(phase)
	printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
	expectedDiscount = decimal.NewFromFloat(5) // 2 cycles * $10 * 50%
	phase.ApplicableNumberOfBillingCycles = 1
	discount, err = cartItem.calculateDiscountForPhase(phase)
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	expectedDiscountAsFloat, _ = expectedDiscount.Round(2).Float64()
	discountAsFloat, _ = discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculatePhaseDiscountPercentRecurringInvalid(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        32, // 2 * DurationUnit
		DurationUnit:                    TimePeriodDaily,
		DiscountValue:                   50,         // 50% discount of the cost of two months
		DiscountType:                    Percentage, // Indicating the DiscountRate is a percentage
		Application:                     Recurring,
		ApplicableNumberOfBillingCycles: 3,
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	discount, err := cartItem.calculateDiscountForPhase(phase)
	printLnDecimalToString(discount, "discount")
	assert.Error(t, err, "Expected error when calculating a percentage discount")

}

func TestCalculatePhaseDiscountPercentSpread(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        32, // 2 * DurationUnit
		DurationUnit:                    TimePeriodDaily,
		DiscountValue:                   50,         // 50% discount of the cost of two months
		DiscountType:                    Percentage, // Indicating the DiscountRate is a percentage
		Application:                     Spread,
		ApplicableNumberOfBillingCycles: 2,
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	discount, err := cartItem.calculateDiscountForPhase(phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected error when calculating a percentage discount")
	expectedDiscount := decimal.NewFromFloat(5.33).Round(2) // 1 cycles * $10 * 50% and 2 days from the second billing cycle
	assert.True(t, expectedDiscount.Equal(discount.Round(2)))

	phase.ApplicableNumberOfBillingCycles = 1
	discount, err = cartItem.calculateDiscountForPhase(phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected error when calculating a percentage discount")
	expectedDiscount = decimal.NewFromFloat(5).Round(2) // 1 cycles * $10 * 50% and 2 days from the second billing cycle
	assert.True(t, expectedDiscount.Equal(discount.Round(2)))
}

func TestCalculatePhaseDiscountPercentTimeInvalidConfig(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	phase := DiscountPhase{
		Duration:      2,          // 2 billing cycles
		DiscountValue: 50,         // 50% discount
		DiscountType:  Percentage, // Indicating the DiscountRate is a percentage
	}
	_, err := cartItem.calculateDiscountForPhase(phase)
	assert.Error(t, err, "Expected an error when calculating a percentage discount based on time")
	//
}

func TestCalculatePhaseDiscountPercentTime(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	phase := DiscountPhase{
		Duration:      2,               // 2 Duration Units
		DiscountValue: 100,             // 100% discount
		DiscountType:  Percentage,      // Indicating the DiscountValue is a percentage
		DurationUnit:  TimePeriodDaily, // The total discount period is DurationUnit * Duration
		Application:   Recurring,       // Apply to Recurring - TODO: We are missing the number of billing cycles
	}
	expectedDiscount := decimal.NewFromFloat(.67) // 1 cycles * $10 / 30 days * 100%
	discount, err := cartItem.calculateDiscountForPhase(phase)
	assert.NoError(t, err, "Expected no error when calculating a percentage discount based on time")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculatePhaseDiscount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	//discountedUnitsQuantity, _ := decimal.NewFromString("2")
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	phase := DiscountPhase{
		Duration:      2,               // 2 billing cycles (billing period unit)
		DurationUnit:  TimePeriodDaily, //period of time to discount
		DiscountValue: 10,              // this will be converted to a 100% percentage discount for the duration since timebased implies 100% discount for the time period
		DiscountType:  TimeBased,       // Indicating the Discountvalue is timebased ( (DiscountValue * DurationUnit)*Duration )
	}
	expectedDiscount := decimal.NewFromFloat(.67) // 1 cycles * $10 * 50%
	discount, err := cartItem.calculateDiscountForPhase(phase)
	assert.NoError(t, err, "Expected no error when calculating a time-based discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculateSubscriptionDiscount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Quantity:  1,
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			BillingPeriodUnit: TimePeriodMonthly,
			TrialPeriod:       1, // 1 month trial
			TrialPeriodUnit:   TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					Duration:      2,
					DiscountValue: 100, // 50% for 2 days
					DurationUnit:  TimePeriodDaily,
					DiscountType:  Percentage,
				},
			},
		},
	}
	//trialDiscount := unitPrice                 // $10 for trial
	phaseDiscount := decimal.NewFromFloat(.67) // $10.67 for 1 month free trial and 2 days at 50%
	expectedDiscount := phaseDiscount          //trialDiscount.Add(phaseDiscount)
	discount, err := cartItem.calculateTotalSubscriptionDiscount()
	assert.NoError(t, err, "Expected no error when calculating a subscription discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculateNonSubscriptionDiscount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	discountAmountPerUnit, _ := decimal.NewFromString("2")
	cartItem := CartItem{
		UnitPrice:                       unitPrice,
		Quantity:                        5,
		DiscountAmountPerDiscountedUnit: discountAmountPerUnit,
		DiscountedUnitsQuantity:         3,
	}
	expectedDiscount := discountAmountPerUnit.Mul(decimal.NewFromInt(3)) // 3 units * $2 discount
	discount, err := cartItem.calculateNonSubscriptionDiscount()
	assert.NoError(t, err, "Expected no error when calculating a non-subscription discount")
	assert.Equal(t, expectedDiscount, discount)
}

func TestCartItem_Validate(t *testing.T) {
	tests := []struct {
		name    string
		item    *CartItem
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid CartItem",
			item: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(10.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(1.0),
				DiscountedUnitsQuantity:         1,
				Subscription: SubscriptionInfo{
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			wantErr: false,
		},
		{
			name: "Negative Quantity",
			item: &CartItem{
				Quantity:                        -1,
				UnitPrice:                       decimal.NewFromFloat(10.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(1.0),
				DiscountedUnitsQuantity:         0,
			},
			wantErr: true,
			errMsg:  "quantity cannot be negative",
		},
		{
			name: "Negative Unit Price",
			item: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(-10.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(1.0),
				DiscountedUnitsQuantity:         1,
			},
			wantErr: true,
			errMsg:  "unit price cannot be negative",
		},
		{
			name: "Negative Discount Amount Per Unit",
			item: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(10.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(-1.0),
				DiscountedUnitsQuantity:         1,
			},
			wantErr: true,
			errMsg:  "discount amount per unit cannot be negative",
		},
		{
			name: "Negative Discounted Units Quantity",
			item: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(10.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(1.0),
				DiscountedUnitsQuantity:         -1,
			},
			wantErr: true,
			errMsg:  "discounted units quantity cannot be negative",
		},
		{
			name: "Discounted Units Quantity Exceeds Total Quantity",
			item: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(10.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(1.0),
				DiscountedUnitsQuantity:         2,
			},
			wantErr: true,
			errMsg:  "discounted units quantity cannot exceed total quantity",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCartItem_Clone(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10.0")
	discountAmountPerUnit, _ := decimal.NewFromString("1.0")
	totalDiscountAmount, _ := decimal.NewFromString("1.0")
	totalGrossAmount, _ := decimal.NewFromString("20.0")
	original := &CartItem{
		SkuID:                           "test-sku",
		Quantity:                        2,
		UnitPrice:                       unitPrice,
		DiscountAmountPerDiscountedUnit: discountAmountPerUnit,
		DiscountedUnitsQuantity:         1,
		TotalDiscountAmount:             totalDiscountAmount,
		TotalGrossAmount:                totalGrossAmount,
	}
	cloned := original.Clone()

	assert.Equal(t, original.SkuID, cloned.SkuID)
	assert.Equal(t, original.Quantity, cloned.Quantity)
	assert.True(t, original.UnitPrice.Equal(cloned.UnitPrice))
	assert.True(t, original.DiscountAmountPerDiscountedUnit.Equal(cloned.DiscountAmountPerDiscountedUnit))
	assert.Equal(t, original.DiscountedUnitsQuantity, cloned.DiscountedUnitsQuantity)
	assert.True(t, original.TotalDiscountAmount.Equal(cloned.TotalDiscountAmount))
	assert.True(t, original.TotalGrossAmount.Equal(cloned.TotalGrossAmount))
	assert.False(t, original == cloned, "Cloned cart item should not be the same instance as the original")
}

func TestCartItem_SetGetSkuID(t *testing.T) {
	cartItem := &CartItem{}
	skuID := "test-sku"
	cartItem.SetSkuID(skuID)
	assert.Equal(t, skuID, cartItem.GetSkuID(), "GetSkuID should return the SKU ID set by SetSkuID")
	cartItem.SetSkuID("") // Test error case
	assert.NotEqual(t, "", cartItem.GetSkuID(), "GetSkuID should not return an empty string after SetSkuID with an empty string")
}

func TestCartItem_SetGetQuantity(t *testing.T) {
	cartItem := &CartItem{}
	var quantity int64 = 5
	cartItem.SetQuantity(quantity)
	assert.Equal(t, quantity, cartItem.GetQuantity(), "GetQuantity should return the quantity set by SetQuantity")
	cartItem.SetQuantity(-1) // Test error case
	assert.NotEqual(t, -1, cartItem.GetQuantity(), "GetQuantity should not return -1 after SetQuantity with -1")
}

func TestCartItem_SetGetUnitPrice(t *testing.T) {
	cartItem := &CartItem{}
	unitPriceStr := "10.00" // Use string representation for precise decimal initialization
	unitPrice, _ := decimal.NewFromString(unitPriceStr)
	cartItem.SetUnitPrice(unitPrice)
	assert.True(t, unitPrice.Equal(cartItem.GetUnitPrice()), "GetUnitPrice should return the unit price set by SetUnitPrice")

	negativeUnitPrice, _ := decimal.NewFromString("-1.00")
	cartItem.SetUnitPrice(negativeUnitPrice) // Attempt to set a negative unit price
	assert.False(t, negativeUnitPrice.Equal(cartItem.GetUnitPrice()), "GetUnitPrice should not return -1 after SetUnitPrice with -1")
}

func TestCartItem_GetGrossTotalAmount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10.00")
	cartItem := &CartItem{
		UnitPrice: unitPrice,
		Quantity:  2,
	}
	expectedTotal := unitPrice.Mul(decimal.NewFromInt(2)) // 2 * unit price
	assert.True(t, expectedTotal.Equal(cartItem.GetGrossTotalAmount()), "GetGrossTotalAmount should return the correct total amount")
}

func TestCartItem_GetNetTotalAmount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10.00")
	discountAmountPerUnit, _ := decimal.NewFromString("1.00")
	cartItem := &CartItem{
		UnitPrice:                       unitPrice,
		Quantity:                        2,
		DiscountAmountPerDiscountedUnit: discountAmountPerUnit,
		DiscountedUnitsQuantity:         1,
	}
	grossTotal := unitPrice.Mul(decimal.NewFromInt(2))                // 2 * unit price
	discountTotal := discountAmountPerUnit.Mul(decimal.NewFromInt(1)) // 1 * discount per unit
	expectedNetTotal := grossTotal.Sub(discountTotal)                 // gross total - discount
	nextTotal, err := cartItem.GetNetTotalAmount()
	assert.NoError(t, err, "GetNetTotalAmount should not return an error")
	assert.True(t, expectedNetTotal.Equal(nextTotal), "GetNetTotalAmount should return the correct net total amount")
}

func TestCalculateTotalDiscountAmount(t *testing.T) {
	tests := []struct {
		name             string
		cartItem         *CartItem
		expectedDiscount decimal.Decimal // Expected discount as string to initialize decimal.Decimal
		expectError      bool
	}{
		{
			name: "No discount when discount per unit is zero",
			cartItem: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(100.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(0),
				DiscountedUnitsQuantity:         1,
				Subscription: SubscriptionInfo{
					IsRecurring:       false,
					TrialPeriod:       0,
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			expectedDiscount: decimal.NewFromFloat(0).Round(2),
			expectError:      false,
		},
		{
			name: "No discount when quantity is zero",
			cartItem: &CartItem{
				Quantity:                        0,
				UnitPrice:                       decimal.NewFromFloat(100.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(5),
				DiscountedUnitsQuantity:         1,
				Subscription: SubscriptionInfo{
					IsRecurring:       false,
					TrialPeriod:       0,
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			expectedDiscount: decimal.NewFromFloat(0).Round(2),
			expectError:      true,
		},
		{
			name: "No discount when unit price is zero",
			cartItem: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.Zero,
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(5),
				DiscountedUnitsQuantity:         1,
				Subscription: SubscriptionInfo{
					IsRecurring:       false,
					TrialPeriod:       0,
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			expectedDiscount: decimal.NewFromFloat(0).Round(2),
			expectError:      false,
		},
		{
			name: "No discount when discounted units quantity is zero",
			cartItem: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(100.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(5),
				DiscountedUnitsQuantity:         0,
				Subscription: SubscriptionInfo{
					IsRecurring:       false,
					TrialPeriod:       0,
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			expectedDiscount: decimal.NewFromFloat(0).Round(2),
			expectError:      false,
		},
		{
			name: "Correct discount calculation",
			cartItem: &CartItem{
				Quantity:                        2,
				UnitPrice:                       decimal.NewFromFloat(100.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(5),
				DiscountedUnitsQuantity:         1,
				Subscription: SubscriptionInfo{
					IsRecurring:       false,
					TrialPeriod:       0,
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			expectedDiscount: decimal.NewFromFloat(5).Round(2),
			expectError:      false,
		},
		{
			name: "Discount does not exceed total gross amount",
			cartItem: &CartItem{
				Quantity:                        1,
				UnitPrice:                       decimal.NewFromFloat(50.0),
				DiscountAmountPerDiscountedUnit: decimal.NewFromFloat(100), // Intentionally high to trigger the cap
				DiscountedUnitsQuantity:         1,
				Subscription: SubscriptionInfo{
					IsRecurring:       false,
					TrialPeriod:       0,
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			expectedDiscount: decimal.NewFromFloat(50).Round(2), // Capped at the total gross amount
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			discount, err := tt.cartItem.GetTotalDiscountAmount()
			if tt.expectError {
				assert.Error(t, err, "GetTotalDiscountAmount should return an error")
			} else {
				assert.NoError(t, err, "GetTotalDiscountAmount should not return an error")
			}
			assert.True(t, tt.expectedDiscount.Equal(discount.Round(2)))
		})
	}
}

func TestCartItem_GetSetTotalDiscountAmountNoPhases(t *testing.T) {
	cartItem := &CartItem{
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	discountAmountStr := "2.00"
	discountAmount, _ := decimal.NewFromString(discountAmountStr)

	// Setting up the cart item
	err := cartItem.SetQuantity(1)
	assert.NoError(t, err, "Setting a valid quantity should not result in an error")
	unitPrice, _ := decimal.NewFromString("10.00")
	cartItem.SetUnitPrice(unitPrice)

	// Setting and getting discount amount per unit
	err = cartItem.SetDiscountAmountPerUnit(decimal.NewFromInt(-1)) // Attempt to set a negative discount amount
	assert.Error(t, err, "Attempting to set a negative discount amount should result in an error")

	err = cartItem.SetDiscountAmountPerUnit(discountAmount) // Valid discount amount
	assert.NoError(t, err, "Setting a valid discount amount should not result in an error")

	_, err = cartItem.GetTotalDiscountAmount()
	assert.Error(t, err, "Getting the discount amount should result in an error")
}

func TestCartItem_GetSetTotalDiscountAmountValid(t *testing.T) {
	cartItem := &CartItem{
		Subscription: SubscriptionInfo{
			IsRecurring: true,
			TrialPeriod: 0,

			BillingPeriodUnit: TimePeriodDaily,
			DiscountPhases: []DiscountPhase{
				{
					Duration:      12,  // 12 hour
					DiscountValue: 100, // 50% for 12 hours
					DurationUnit:  TimePeriodHourly,
					DiscountType:  Percentage,
				},
			},
		},
	}
	discountAmountStr := "5.0"
	discountAmount, _ := decimal.NewFromString(discountAmountStr)

	// Setting up the cart item
	err := cartItem.SetQuantity(1)
	assert.NoError(t, err, "Setting a valid quantity should not result in an error")
	unitPrice, _ := decimal.NewFromString("10.00")
	cartItem.SetUnitPrice(unitPrice)

	// Setting and getting discount amount per unit
	//err = cartItem.SetDiscountAmountPerUnit(decimal.NewFromInt(-1)) // Attempt to set a negative discount amount
	//assert.Error(t, err, "Attempting to set a negative discount amount should result in an error")

	err = cartItem.SetDiscountAmountPerUnit(discountAmount) // Valid discount amount
	assert.NoError(t, err, "Setting a valid discount amount should not result in an error")

	receivedDiscountAmount, err := cartItem.GetTotalDiscountAmount()
	assert.NoError(t, err, "Getting the discount amount should not result in an error")

	printLnDecimalToString(receivedDiscountAmount, "receivedDiscountAmount")
	discountAmountfloat, _ := discountAmount.Round(2).Float64()
	receivedDiscountAmountfloat, _ := receivedDiscountAmount.Round(2).Float64()
	assert.True(t, discountAmount.Round(2).Equal(receivedDiscountAmount.Round(2)), "The received discount amount should match the set value", discountAmountfloat, receivedDiscountAmountfloat)

}

func TestCartItem_GetSetTotalDiscountAmountValidPhaseDurationMatchesBilling(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10.00")
	cartItem := &CartItem{
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			BillingPeriodUnit: TimePeriodDaily,
			DiscountPhases: []DiscountPhase{
				{
					Duration:      12, // 12 hour
					DiscountValue: 50, // 50% for 12 hours
					DurationUnit:  TimePeriodHourly,
					DiscountType:  Percentage,
				},
			},
		},
		UnitPrice: unitPrice,
	}
	discountAmountStr := "2.5"
	discountAmount, _ := decimal.NewFromString(discountAmountStr)

	// Setting up the cart item
	err := cartItem.SetQuantity(1)
	assert.NoError(t, err, "Setting a valid quantity should not result in an error")

	cartItem.SetUnitPrice(unitPrice)

	// Setting and getting discount amount per unit
	err = cartItem.SetDiscountAmountPerUnit(decimal.NewFromInt(-1)) // Attempt to set a negative discount amount
	assert.Error(t, err, "Attempting to set a negative discount amount should result in an error")

	err = cartItem.SetDiscountAmountPerUnit(discountAmount) // Valid discount amount
	assert.NoError(t, err, "Setting a valid discount amount should not result in an error")

	receivedDiscountAmount, err := cartItem.GetTotalDiscountAmount()
	assert.NoError(t, err, "Getting the discount amount should not result in an error")

	printLnDecimalToString(receivedDiscountAmount, "receivedDiscountAmount")
	discountAmountfloat, _ := discountAmount.Round(2).Float64()
	receivedDiscountAmountfloat, _ := receivedDiscountAmount.Round(2).Float64()
	assert.True(t, discountAmount.Round(2).Equal(receivedDiscountAmount.Round(2)), "The received discount amount should match the set value", discountAmountfloat, receivedDiscountAmountfloat)
}

func TestCartItem_GetSetTotalDiscountAmountValid2(t *testing.T) {
	cartItem := &CartItem{
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					Duration:      15,
					DiscountValue: 100, // 50% for 12 hours
					DurationUnit:  TimePeriodDaily,
					DiscountType:  Percentage,
				},
			},
		},
	}
	discountAmountStr := "5"
	discountAmount, _ := decimal.NewFromString(discountAmountStr)

	// Setting up the cart item
	err := cartItem.SetQuantity(1)
	assert.NoError(t, err, "Setting a valid quantity should not result in an error")
	unitPrice, _ := decimal.NewFromString("10.00")
	cartItem.SetUnitPrice(unitPrice)

	// Setting and getting discount amount per unit
	//err = cartItem.SetDiscountAmountPerUnit(decimal.NewFromInt(-1)) // Attempt to set a negative discount amount
	//assert.Error(t, err, "Attempting to set a negative discount amount should result in an error")

	err = cartItem.SetDiscountAmountPerUnit(discountAmount) // Valid discount amount
	assert.NoError(t, err, "Setting a valid discount amount should not result in an error")

	receivedDiscountAmount, err := cartItem.GetTotalDiscountAmount()
	assert.NoError(t, err, "Getting the discount amount should not result in an error")

	printLnDecimalToString(receivedDiscountAmount, "receivedDiscountAmount")
	discountAmountfloat, _ := discountAmount.Round(2).Float64()
	receivedDiscountAmountfloat, _ := receivedDiscountAmount.Round(2).Float64()
	assert.True(t, discountAmount.Round(2).Equal(receivedDiscountAmount.Round(2)), "The received discount amount should match the set value", discountAmountfloat, receivedDiscountAmountfloat)
}

func TestCartItem_GetSetTotalDiscountAmountValid3(t *testing.T) {
	cartItem := &CartItem{
		Subscription: SubscriptionInfo{
			IsRecurring:       true,
			TrialPeriod:       0,
			BillingPeriodUnit: TimePeriodHourly,
			DiscountPhases: []DiscountPhase{
				{
					Duration:      1,
					DiscountValue: 100, // 50% for 12 hours
					DurationUnit:  TimePeriodDaily,
					DiscountType:  Percentage,
				},
			},
		},
	}
	discountAmountStr := "5.00"
	discountAmount, _ := decimal.NewFromString(discountAmountStr)

	// Setting up the cart item
	err := cartItem.SetQuantity(1)
	assert.NoError(t, err, "Setting a valid quantity should not result in an error")
	unitPrice, _ := decimal.NewFromString("10.00")
	cartItem.SetUnitPrice(unitPrice)

	// Setting and getting discount amount per unit
	//err = cartItem.SetDiscountAmountPerUnit(decimal.NewFromInt(-1)) // Attempt to set a negative discount amount
	//assert.Error(t, err, "Attempting to set a negative discount amount should result in an error")

	err = cartItem.SetDiscountAmountPerUnit(discountAmount) // Valid discount amount
	assert.NoError(t, err, "Setting a valid discount amount should not result in an error")

	_, err = cartItem.GetTotalDiscountAmount()
	assert.Error(t, err, "Getting the discount amount should result in an error")

}

func TestCartItem_IncrementDecrementCartItemQuantity(t *testing.T) {
	cartItem := &CartItem{Quantity: 5}

	// Increment
	incrementedQuantity, err := cartItem.IncrementCartItemQuantity(2)
	expectedQuantityAfterIncrement := int64(7)
	assert.Nil(t, err, "IncrementCartItemQuantity should not return an error for valid increment")
	assert.Equal(t, expectedQuantityAfterIncrement, incrementedQuantity, "IncrementCartItemQuantity should correctly increment the quantity")

	// Error case for increment
	qty, err := cartItem.IncrementCartItemQuantity(-1)
	assert.Equal(t, expectedQuantityAfterIncrement, qty, "IncrementCartItemQuantity should remain the same")
	assert.NotNil(t, err, "IncrementCartItemQuantity should return an error when incrementing with a negative quantity")

	// Decrement
	decrementedQuantity, err := cartItem.DecrementCartItemQuantity(4)
	expectedQuantityAfterDecrement := int64(3)
	assert.Nil(t, err, "DecrementCartItemQuantity should not return an error for valid decrement")
	assert.Equal(t, expectedQuantityAfterDecrement, decrementedQuantity, "DecrementCartItemQuantity should correctly decrement the quantity")

	// Error case for decrement
	_, err = cartItem.DecrementCartItemQuantity(4)
	assert.NotNil(t, err, "DecrementCartItemQuantity should return an error when decrementing below 0")

	originalQty := cartItem.GetQuantity()
	// Error case for negative increment
	qty, err = cartItem.DecrementCartItemQuantity(-1)
	assert.Equal(t, originalQty, qty, "DecrementCartItemQuantity should remain the same")
	assert.NotNil(t, err, "DecrementCartItemQuantity should return an error when decrementing with a negative quantity")
}

// Additional tests can include edge cases such as setting negative quantities, unit prices, or discount amounts,
// and verifying the behavior of `UpdateTotals` to ensure it correctly updates the cart item's totals based on its properties.
