package common

import (
	"fmt"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert" // Using testify for easier assertions
)

func TestCartItemValidate(t *testing.T) {
	// Test case 1: Negative quantity
	cartItem := &CartItem{
		Quantity: -1,
	}
	err := cartItem.Validate()
	if err == nil || err.Error() != "quantity cannot be negative" {
		t.Errorf("Test case 1 failed: Expected error 'quantity cannot be negative', got %v", err)
	}

	// Test case 2: Negative unit price
	cartItem = &CartItem{
		Quantity:  1,
		UnitPrice: decimal.NewFromFloat(-10.0),
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "unit price cannot be negative" {
		t.Errorf("Test case 2 failed: Expected error 'unit price cannot be negative', got %v", err)
	}

	// Test case 3: Negative discount amount per unit
	cartItem = &CartItem{
		Quantity:                       1,
		DiscountValuePerDiscountedUnit: decimal.NewFromFloat(-5.0),
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "discount amount per unit cannot be negative" {
		t.Errorf("Test case 3 failed: Expected error 'discount amount per unit cannot be negative', got %v", err)
	}

	// Test case 4: Discounted units quantity exceeds total quantity
	cartItem = &CartItem{
		Quantity:                5,
		NumberOfUnitsDiscounted: 6,
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "discounted units quantity cannot exceed total quantity" {
		t.Errorf("Test case 4 failed: Expected error 'discounted units quantity cannot exceed total quantity', got %v", err)
	}

	// Test case 5: Invalid billing period unit
	cartItem = &CartItem{
		Quantity: 1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodUnknown,
		},
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "invalid billing period unit" {
		t.Errorf("Test case 5 failed: Expected error 'invalid billing period unit', got %v", err)
	}

	// Test case 6: Invalid duration unit for discount phase
	cartItem = &CartItem{
		Quantity: 1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					DurationUnit: TimePeriodUnknown,
				},
			},
		},
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "invalid duration unit for discount phase" {
		t.Errorf("Test case 6 failed: Expected error 'invalid duration unit for discount phase', got %v", err)
	}

	// Test case 7: Invalid duration for discount phase
	cartItem = &CartItem{
		Quantity: 1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					Duration: -1,
				},
			},
		},
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "invalid duration unit for discount phase" {
		t.Errorf("Test case 7 failed: Expected error 'invalid duration unit for discount phase', got %v", err)
	}

	// Test case 8: Discount phase duration unit exceeds billing period unit
	cartItem = &CartItem{
		Quantity: 1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					Duration:     1,
					DurationUnit: TimePeriodAnnual,
				},
			},
		},
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "discount phase duration unit cannot exceed billing period unit" {
		t.Errorf("Test case 8 failed: Expected error 'discount phase duration unit cannot exceed billing period unit', got %v", err)
	}

	// Test case 9: Invalid percentage discount phase
	cartItem = &CartItem{
		Quantity: 1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					DiscountType: Percentage,
					DurationUnit: TimePeriodMonthly,
					Duration:     1,
					// Missing required fields for percentage discount
				},
			},
		},
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "percentage discount rate must be positive" {
		t.Errorf("Test case 9 failed: Expected error 'percentage discount rate must be positive', got %v", err)
	}

	// Test case 10: Invalid fixed amount discount phase
	cartItem = &CartItem{
		Quantity: 1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					DiscountType: FixedAmount,
					// Missing required fields for fixed amount discount
				},
			},
		},
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "fixed amount discount must be positive" {
		t.Errorf("Test case 10 failed: Expected error 'fixed amount discount must be positive', got %v", err)
	}

	// Test case 11: Invalid duration for discount phase
	cartItem = &CartItem{
		Quantity: 1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases: []DiscountPhase{
				{
					DiscountType: Percentage,
					DurationUnit: TimePeriodMonthly,
					Duration:     0,
					// Missing required fields for fixed amount discount
				},
			},
		},
	}
	err = cartItem.Validate()
	if err == nil || err.Error() != "invalid duration for discount phase" {
		t.Errorf("Test case 10 failed: Expected error 'invalid duration for discount phase', got %v", err)
	}
}

func TestDiscountedUnitsQuantity(t *testing.T) {
	// Test case 1: GetDiscountedUnitsQuantity() returns the initial value (0)
	cartItem := &CartItem{}
	if qty := cartItem.GetDiscountedUnitsQuantity(); qty != 0 {
		t.Errorf("Test case 1 failed: Expected 0, got %d", qty)
	}

	// Test case 2: SetDiscountedUnitQuantity() sets the quantity correctly
	newQty := int64(5)
	err := cartItem.SetDiscountedUnitQuantity(newQty)
	if err != nil {
		t.Errorf("Test case 2 failed: Unexpected error: %v", err)
	}
	if qty := cartItem.GetDiscountedUnitsQuantity(); qty != newQty {
		t.Errorf("Test case 2 failed: Expected %d, got %d", newQty, qty)
	}

	// Test case 3: SetDiscountedUnitQuantity() doesn't allow negative quantity
	err = cartItem.SetDiscountedUnitQuantity(-2)
	expectedErr := "discounted units quantity cannot be negative"
	if err == nil {
		t.Errorf("Test case 3 failed: Expected error '%s', got nil", expectedErr)
	} else if err.Error() != expectedErr {
		t.Errorf("Test case 3 failed: Expected error '%s', got '%s'", expectedErr, err.Error())
	}
	// Ensure the quantity wasn't changed
	if qty := cartItem.GetDiscountedUnitsQuantity(); qty != newQty {
		t.Errorf("Test case 3 failed: Quantity should remain %d, got %d", newQty, qty)
	}
}

func TestSetDiscountAmountPerUnitFromString(t *testing.T) {
	// Test case 1: Valid discount amount string
	cartItem := &CartItem{}
	discountAmountStr := "10.50"
	err := cartItem.SetDiscountAmountPerUnitFromString(discountAmountStr)
	if err != nil {
		t.Errorf("Test case 1 failed: Unexpected error: %v", err)
	}
	expectedDiscountAmount, _ := decimal.NewFromString(discountAmountStr)
	if !cartItem.DiscountValuePerDiscountedUnit.Equal(expectedDiscountAmount) {
		t.Errorf("Test case 1 failed: Expected discount amount %s, got %s", expectedDiscountAmount.String(), cartItem.DiscountValuePerDiscountedUnit.String())
	}

	// Test case 2: Invalid discount amount string
	cartItem = &CartItem{}
	invalidDiscountAmountStr := "invalid"
	err = cartItem.SetDiscountAmountPerUnitFromString(invalidDiscountAmountStr)
	if err == nil {
		t.Error("Test case 2 failed: Expected an error for invalid discount amount string, but got nil")
	} else if err.Error() != "invalid discount amount format" {
		t.Errorf("Test case 2 failed: Expected error 'invalid discount amount format', got '%s'", err.Error())
	}

	// Test case 3: Negative discount amount string
	cartItem = &CartItem{}
	negativeDiscountAmountStr := "-10.50"
	err = cartItem.SetDiscountAmountPerUnitFromString(negativeDiscountAmountStr)
	if err == nil {
		t.Error("Test case 3 failed: Expected an error for negative discount amount, but got nil")
	} else if err.Error() != "discount amount cannot be negative" {
		t.Errorf("Test case 3 failed: Expected error 'discount amount cannot be negative', got '%s'", err.Error())
	}
}

func TestCalculateDiscountForPhase(t *testing.T) {
	// Test case 1: Percentage discount, recurring application
	cartItem := &CartItem{
		UnitPrice: decimal.NewFromFloat(100),
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	phase := DiscountPhase{
		DiscountValue: 20,
		Duration:      1,
		DurationUnit:  TimePeriodMonthly,
		Application:   Recurring,
	}
	phase.SetDiscountType(Percentage)
	if phase.GetDiscountType() != Percentage {
		t.Errorf("Test case 1 failed: Expected DiscountType to be Percentage, got %v", phase.GetDiscountType())
	}
	phase.SetApplicableNumberOfBillingCycles(12)
	if phase.GetApplicableNumberOfBillingCycles() != 12 {
		t.Errorf("Test case 1 failed: Expected ApplicableNumberOfBillingCycles to be 12, got %d", phase.GetApplicableNumberOfBillingCycles())
	}
	discount, err := cartItem.calculateDiscountForPhase(&phase)
	if err != nil {
		t.Errorf("Error in test case 1: %v", err)
	}
	expectedDiscount := decimal.NewFromFloat(20)
	if !discount.Equal(expectedDiscount) {
		t.Errorf("Test case 1 failed: Expected discount %s, got %s", expectedDiscount.String(), discount.String())
	}
	/*
		if len(cartItem.Logs) == 0 {
			t.Errorf("Test case 1 failed: No logs recorded")
		}
	*/
	// Test case 2: Fixed amount discount, spread application
	cartItem = &CartItem{
		UnitPrice: decimal.NewFromFloat(100),
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	phase = DiscountPhase{
		DiscountType:                    FixedAmount,
		DiscountValue:                   30,
		Duration:                        2,                 // Not relevant for fixed value
		DurationUnit:                    TimePeriodMonthly, // Not relevant for fixed value
		Application:                     OneTime,
		ApplicableNumberOfBillingCycles: 6,
	}
	phase.SetDescription("Fixed amount discount of $30 for 6 months")
	if phase.GetDescription() != "Fixed amount discount of $30 for 6 months" {
		t.Errorf("Test case 2 failed: Expected description to be 'Fixed amount discount of $30 for 6 months', got %s", phase.GetDescription())
	}

	discount, err = cartItem.calculateDiscountForPhase(&phase)
	if err != nil {
		t.Errorf("Error in test case 2: %v", err)
	}
	expectedDiscount = decimal.NewFromFloat(180)
	if !discount.Equal(expectedDiscount) {
		t.Errorf("Test case 2 failed: Expected discount %s, got %s", expectedDiscount.String(), discount.String())
	}
	/*
		if len(cartItem.Logs) == 0 {
			t.Errorf("Test case 2 failed: No logs recorded")
		}*/

	// Test case 3: Percentage discount, invalid duration unit
	cartItem = &CartItem{
		UnitPrice: decimal.NewFromFloat(100),
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}
	phase = DiscountPhase{
		DiscountType:                    Percentage,
		DiscountValue:                   20,
		Duration:                        1,
		DurationUnit:                    TimePeriodUnknown,
		Application:                     Recurring,
		ApplicableNumberOfBillingCycles: 12,
	}
	_, err = cartItem.calculateDiscountForPhase(&phase)
	if err == nil {
		t.Errorf("Test case 3 failed: Expected error for invalid duration unit")
	}
	//fmt.Println("Error in test case 3:", err)
	phase.SetDurationUnit(TimePeriodMonthly)
	phase.SetDiscountType(Unknown)
	_, err = cartItem.calculateDiscountForPhase(&phase)
	assert.Error(t, err, "Expected error when calculating a discount for an unknown discount type")
}

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
	discount, err := cartItem.calculateTotalTrialPeriodDiscount()
	assert.NoError(t, err, "Expected no error when calculating a trial period discount")
	assert.True(t, expectedDiscount.Equal(discount))
	cartItem.Subscription.TrialPeriodUnit = TimePeriodUnknown
	_, err = cartItem.calculateTotalTrialPeriodDiscount()
	assert.Error(t, err, "Expected error when calculating a trial period discount")
	cartItem.Subscription.TrialPeriodUnit = TimePeriodMonthly
	cartItem.Subscription.TrialPeriod = 0
	_, err = cartItem.calculateTotalTrialPeriodDiscount()
	assert.NoError(t, err, "Expected no error when calculating a trial period discount")
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
	discount, err := cartItem.calculateTotalSubscriptionDiscounts()
	assert.NoError(t, err, "Expected no error when calculating a subscription discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f got: %f", expectedDiscountAsFloat, discountAsFloat))
	cartItem.Subscription.DiscountPhases[0].DiscountValue = -1
	_, err = cartItem.calculateTotalSubscriptionDiscounts()
	assert.Error(t, err, "Expected error when calculating a subscription discount")
}

func TestSetUnitPriceFromString(t *testing.T) {
	// Create a CartItem instance
	cartItem := &CartItem{
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodMonthly,
		},
	}

	// Test case 1: Valid unit price string
	unitPriceStr := "10.99"
	err := cartItem.SetUnitPriceFromString(unitPriceStr)
	if err != nil {
		t.Errorf("Test case 1 failed: Expected no error, got %v", err)
	}
	expectedUnitPrice, _ := decimal.NewFromString(unitPriceStr)
	if !cartItem.UnitPrice.Equal(expectedUnitPrice) {
		t.Errorf("Test case 1 failed: Expected unit price %s, got %s", expectedUnitPrice.String(), cartItem.UnitPrice.String())
	}

	// Test case 2: Invalid unit price string
	unitPriceStr = "invalid"
	err = cartItem.SetUnitPriceFromString(unitPriceStr)
	if err == nil {
		t.Errorf("Test case 2 failed: Expected an error, got nil")
	}
	expectedErrorMessage := "can't convert invalid to decimal"
	if err.Error() != expectedErrorMessage {
		t.Errorf("Test case 2 failed: Expected error message \"%s\", got \"%s\"", expectedErrorMessage, err.Error())
	}
}

func TestCartItem_SetGetQuantity(t *testing.T) {
	cartItem := &CartItem{}
	var quantity int64 = 5
	cartItem.SetQuantity(quantity)
	assert.Equal(t, quantity, cartItem.GetQuantity(), "GetQuantity should return the quantity set by SetQuantity")
	cartItem.SetQuantity(-1) // Test error case
	assert.NotEqual(t, -1, cartItem.GetQuantity(), "GetQuantity should not return -1 after SetQuantity with -1")
	cartItem.IsSubscription = true
	cartItem.SetQuantity(2) // Test error case
}

func TestDiscountDescription(t *testing.T) {
	// Create a CartItem instance
	cartItem := &CartItem{}

	// Test case 1: GetDiscountDescription() returns an empty string initially
	if desc := cartItem.GetDiscountDescription(); desc != "" {
		t.Errorf("Test case 1 failed: Expected empty string, got %s", desc)
	}

	// Test case 2: SetDiscountDescription() sets the description correctly
	description := "10% off for the first month"
	cartItem.SetDiscountDescription(description)
	if desc := cartItem.GetDiscountDescription(); desc != description {
		t.Errorf("Test case 2 failed: Expected %s, got %s", description, desc)
	}

	// Test case 3: SetDiscountDescription() updates the description correctly
	newDescription := "15% off for the first month"
	cartItem.SetDiscountDescription(newDescription)
	if desc := cartItem.GetDiscountDescription(); desc != newDescription {
		t.Errorf("Test case 3 failed: Expected %s, got %s", newDescription, desc)
	}
}

func TestDiscountTypeString(t *testing.T) {
	// Test case 1: Percentage discount
	discountType := Percentage
	expectedString := "Percentage"
	if result := discountType.String(); result != expectedString {
		t.Errorf("Test case 1 failed: Expected %s, got %s", expectedString, result)
	}

	// Test case 2: TimeBased discount
	discountType = TimeBased
	expectedString = "TimeBased"
	if result := discountType.String(); result != expectedString {
		t.Errorf("Test case 2 failed: Expected %s, got %s", expectedString, result)
	}

	// Test case 3: FixedAmount discount
	discountType = FixedAmount
	expectedString = "FixedAmount"
	if result := discountType.String(); result != expectedString {
		t.Errorf("Test case 3 failed: Expected %s, got %s", expectedString, result)
	}

	// Test case 4: Unknown discount
	discountType = DiscountType(-1) // Invalid discount type
	expectedString = "Unknown"
	if result := discountType.String(); result != expectedString {
		t.Errorf("Test case 4 failed: Expected %s, got %s", expectedString, result)
	}
}

func TestSetDiscountType(t *testing.T) {
	dp := &DiscountPhase{
		DiscountType:  FixedAmount,
		DiscountValue: 10,
		DurationUnit:  TimePeriodDaily,
		Duration:      30,
	}
	err := dp.SetDiscountType(Percentage)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if dp.DiscountType != Percentage {
		t.Errorf("Expected DiscountType to be Percentage, got %v", dp.DiscountType)
	}

	err = dp.SetDiscountType(TimeBased)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if dp.DiscountType != TimeBased {
		t.Errorf("Expected DiscountType to be TimeBased, got %v", dp.DiscountType)
	}
	if dp.DiscountValue != 100 { // all timebased discounts are 100%
		t.Errorf("Expected DiscountValue to be 100, got %v", dp.DiscountValue)
	}

	dp.SetDuration(0, TimePeriodUnknown)
	dp.SetApplication(OneTime)
	dp.SetApplicableNumberOfBillingCycles(1)
	err = dp.SetDiscountType(FixedAmount)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if dp.DiscountType != FixedAmount {
		t.Errorf("Expected DiscountType to be FixedAmount, got %v", dp.DiscountType)
	}

}

func TestValidateFixedAmountDiscount(t *testing.T) {
	dp := &DiscountPhase{
		DiscountValue:                   50,
		DurationUnit:                    TimePeriodUnknown,
		Duration:                        0,
		Application:                     Recurring,
		ApplicableNumberOfBillingCycles: 3,
	}
	err := dp.ValidateFixedAmountDiscount()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	dp.SetDiscountValue(-1)
	err = dp.ValidateFixedAmountDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	dp.SetDiscountValue(1)
	dp.SetDurationUnit(TimePeriodDaily)
	err = dp.ValidateFixedAmountDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	dp.SetDiscountValue(1)
	dp.SetDuration(1, TimePeriodUnknown)
	err = dp.ValidateFixedAmountDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetDiscountValue(50)
	dp.SetDuration(0, TimePeriodUnknown)
	dp.SetApplication(UnknownApplication)
	err = dp.ValidateFixedAmountDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetApplication(OneTime)
	dp.SetApplicableNumberOfBillingCycles(-1)
	err = dp.ValidateFixedAmountDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetApplicableNumberOfBillingCycles(2)
	err = dp.ValidateFixedAmountDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetDuration(30, TimePeriodDaily)

	// Add more test cases covering different scenarios and error cases...
}

func TestValidatePercentageDiscount(t *testing.T) {
	dp := &DiscountPhase{
		DiscountValue:                   50,
		DurationUnit:                    TimePeriodDaily,
		Duration:                        30,
		Application:                     Recurring,
		ApplicableNumberOfBillingCycles: 3,
	}
	err := dp.ValidatePercentageDiscount()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	dp.SetDiscountValue(-1)
	err = dp.ValidatePercentageDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetDiscountValue(101)
	err = dp.ValidatePercentageDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetDiscountValue(50)
	dp.SetDurationUnit(TimePeriodUnknown)
	err = dp.ValidatePercentageDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetDurationUnit(TimePeriodDaily)
	dp.SetDuration(0, TimePeriodDaily)
	err = dp.ValidatePercentageDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	dp.SetDuration(30, TimePeriodDaily)
	dp.SetApplication(UnknownApplication)
	err = dp.ValidatePercentageDiscount()
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Add more test cases covering different scenarios and error cases...
}

func TestGetSetApplication(t *testing.T) {
	dp := &DiscountPhase{Application: Recurring}
	application := dp.GetApplication()
	if application != Recurring {
		t.Errorf("Expected Application to be Recurring, got %v", application)
	}

	dp.SetApplication(Spread)
	application = dp.GetApplication()
	if application != Spread {
		t.Errorf("Expected Application to be Spread, got %v", application)
	}
}

func TestGetSetDuration(t *testing.T) {
	dp := &DiscountPhase{Duration: 2}
	application := dp.GetDuration()
	if application != 2 {
		t.Errorf("Expected Application to be Recurring, got %v", application)
	}

	dp.SetDuration(3, TimePeriodMonthly)
	application = dp.GetDuration()
	if application != 3 {
		t.Errorf("Expected Application to be Spread, got %v", application)
	}
	durationUnit := dp.GetDurationUnit()
	if durationUnit != TimePeriodMonthly {
		t.Errorf("Expected DurationUnit to be Monthly, got %v", durationUnit)
	}

	dp.SetDurationUnit(TimePeriodWeekly)
	durationUnit = dp.GetDurationUnit()
	if durationUnit != TimePeriodWeekly {
		t.Errorf("Expected DurationUnit to be Weekly, got %v", durationUnit)
	}

}

func TestGetSetDiscountValue(t *testing.T) {
	dp := &DiscountPhase{DiscountValue: 10.5}
	value := dp.GetDiscountValue()
	if value != 10.5 {
		t.Errorf("Expected DiscountValue to be 10.5, got %v", value)
	}

	dp.SetDiscountValue(15.75)
	value = dp.GetDiscountValue()
	if value != 15.75 {
		t.Errorf("Expected DiscountValue to be 15.75, got %v", value)
	}
}

func TestValidateTimeBasedDiscount(t *testing.T) {
	dp := &DiscountPhase{
		DurationUnit:                    TimePeriodMonthly,
		Duration:                        10,
		Application:                     Spread,
		ApplicableNumberOfBillingCycles: 5,
		DiscountValue:                   100,
	}
	err := dp.ValidateTimeBasedDiscount()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
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

	_, err := cartItem.calculateDiscountForPhase(&phase)
	assert.Error(t, err, "Expected error when calculating a percentage discount")

}

func TestCalculatePhaseDiscountNegativeValue(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	cartItem := CartItem{
		UnitPrice: unitPrice,
	}
	phase := DiscountPhase{
		Duration:      2,          // 2 billing cycles
		DiscountValue: -50,        // 50% discount
		DiscountType:  Percentage, // Indicating the DiscountRate is a percentage
	}

	_, err := cartItem.calculateDiscountForPhase(&phase)
	assert.Error(t, err, "Expected error when calculating a percentage discount")

}

func TestCalculatePhaseFixedTooBigValue(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        0,                 // Not relevant for fixed value
		DurationUnit:                    TimePeriodUnknown, // Not relevant for fixed value
		DiscountValue:                   15,                // 50% discount of the cost of two days
		DiscountType:                    FixedAmount,       // Indicating the DiscountRate is a percentage
		Application:                     OneTime,           // Only recurring and Onetime are valid for fixed value. Default is onetime
		ApplicableNumberOfBillingCycles: 1,                 // 1 billing cycle
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       false,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	expectedDiscount := decimal.NewFromFloat(10) // 2 cycles * $10 * 50%
	discount, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	assert.True(t, len(phase.DiscountsPerBillingCycle) == 1, "The discount per billing cycle should be 1")
	for _, amount := range phase.DiscountsPerBillingCycle {
		assert.Equal(t, float64(10), amount, "The discount per billing cycle should be 5")
	}

	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculatePhaseFixedValue(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        2,               // Not relevant for fixed value
		DurationUnit:                    TimePeriodDaily, // Not relevant for fixed value
		DiscountValue:                   5,               // 50% discount of the cost of two days
		DiscountType:                    FixedAmount,     // Indicating the DiscountRate is a percentage
		Application:                     OneTime,         // Only recurring and Onetime are valid for fixed value. Default is onetime
		ApplicableNumberOfBillingCycles: 1,               // 1 billing cycle
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       false,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	expectedDiscount := decimal.NewFromFloat(5) // 2 cycles * $10 * 50%
	discount, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	assert.True(t, len(phase.DiscountsPerBillingCycle) == 1, "The discount per billing cycle should be 1")
	for _, amount := range phase.DiscountsPerBillingCycle {
		assert.Equal(t, float64(5), amount, "The discount per billing cycle should be 5")
	}

	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculatePhaseFixedInvalidApplicationValue(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        2,                  // Not relevant for fixed value
		DurationUnit:                    TimePeriodDaily,    // Not relevant for fixed value
		DiscountValue:                   5,                  // 50% discount of the cost of two days
		DiscountType:                    FixedAmount,        // Indicating the DiscountRate is a percentage
		Application:                     UnknownApplication, // Only recurring and Onetime are valid for fixed value. Default is onetime
		ApplicableNumberOfBillingCycles: 1,                  // 1 billing cycle
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       false,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	_, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.Error(t, err, "Expected error when calculating a fixed discount")

}

func TestCalculatePhaseFixedValueRecurring(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        2,               // Not relevant for fixed value
		DurationUnit:                    TimePeriodDaily, // Not relevant for fixed value
		DiscountValue:                   5,               // $5 discount
		DiscountType:                    FixedAmount,     // Indicating the DiscountRate is a fixedamount
		Application:                     Recurring,       // Only recurring and Onetime are valid for fixed value
		ApplicableNumberOfBillingCycles: 2,               // 2 billing cycles
	}
	cartItem := CartItem{
		UnitPrice: unitPrice,
		Subscription: SubscriptionInfo{
			IsRecurring:       false,
			TrialPeriod:       0,
			TrialPeriodUnit:   TimePeriodNoBilling,
			BillingPeriodUnit: TimePeriodMonthly,
			DiscountPhases:    []DiscountPhase{phase},
		},
	}

	expectedDiscount := decimal.NewFromFloat(10) // 2 cycles * $10 * 50%
	discount, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")

	assert.True(t, len(phase.DiscountsPerBillingCycle) == 2, "The discount per billing cycle should be 2")
	for _, amount := range phase.DiscountsPerBillingCycle {
		assert.Equal(t, float64(5), amount, "The discount per billing cycle should be 5")
	}

	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
	/*
		expectedDiscount = decimal.NewFromFloat(.33) // 2 cycles * $10 * 50%
		discount, err = cartItem.calculateTotalDiscountValue(phase)
		assert.NoError(t, err, "Expected no error when calculating a percentage discount")
		expectedDiscountAsFloat, _ = expectedDiscount.Round(2).Float64()
		discountAsFloat, _ = discount.Round(2).Float64()
		assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
	*/
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
	discount, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
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
	discount, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected no error when calculating a percentage discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f and got: %f", expectedDiscountAsFloat, discountAsFloat))
	expectedDiscount = decimal.NewFromFloat(5) // 2 cycles * $10 * 50%
	phase.ApplicableNumberOfBillingCycles = 1
	discount, err = cartItem.calculateDiscountForPhase(&phase)
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

	_, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
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

	discount, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected error when calculating a percentage discount")
	expectedDiscount := decimal.NewFromFloat(5.33).Round(2) // 1 cycles * $10 * 50% and 2 days from the second billing cycle
	assert.True(t, expectedDiscount.Equal(discount.Round(2)))

	phase.ApplicableNumberOfBillingCycles = 1
	discount, err = cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected error when calculating a percentage discount")
	expectedDiscount = decimal.NewFromFloat(5).Round(2) // 1 cycles * $10 * 50% and 2 days from the second billing cycle
	assert.True(t, expectedDiscount.Equal(discount.Round(2)))
}

func TestCalculatePhaseDiscountPercentBigDiscount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	phase := DiscountPhase{
		Duration:                        32, // 2 * DurationUnit
		DurationUnit:                    TimePeriodDaily,
		DiscountValue:                   100,        // 50% discount of the cost of two months
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

	discount, err := cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected error when calculating a percentage discount")
	expectedDiscount := decimal.NewFromFloat(10.67).Round(2) // discount capped at the unit price
	assert.True(t, expectedDiscount.Equal(discount.Round(2)))

	phase.ApplicableNumberOfBillingCycles = 1
	discount, err = cartItem.calculateDiscountForPhase(&phase)
	//printLnDecimalToString(discount, "discount")
	assert.NoError(t, err, "Expected error when calculating a percentage discount")
	expectedDiscount = decimal.NewFromFloat(10).Round(2) // 1 cycles * $10 * 50% and 2 days from the second billing cycle
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
	_, err := cartItem.calculateDiscountForPhase(&phase)
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
	discount, err := cartItem.calculateDiscountForPhase(&phase)
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
	discount, err := cartItem.calculateDiscountForPhase(&phase)
	assert.NoError(t, err, "Expected no error when calculating a time-based discount")
	expectedDiscountAsFloat, _ := expectedDiscount.Round(2).Float64()
	discountAsFloat, _ := discount.Round(2).Float64()
	assert.True(t, expectedDiscount.Round(2).Equal(discount.Round(2)), fmt.Sprintf("The calculated discount should match the expected discount wanted: %f got: %f", expectedDiscountAsFloat, discountAsFloat))
}

func TestCalculateNonSubscriptionDiscount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10")
	discountAmountPerUnit, _ := decimal.NewFromString("2")
	cartItem := CartItem{
		UnitPrice:                      unitPrice,
		Quantity:                       5,
		DiscountValuePerDiscountedUnit: discountAmountPerUnit,
		NumberOfUnitsDiscounted:        3,
	}
	expectedDiscount := discountAmountPerUnit.Mul(decimal.NewFromInt(3)) // 3 units * $2 discount
	discount, err := cartItem.calculateNonSubscriptionDiscount()
	assert.NoError(t, err, "Expected no error when calculating a non-subscription discount")
	assert.Equal(t, expectedDiscount, discount)

	// Additional test case when DiscountType is Percentage
	cartItem.DiscountType = Percentage
	expectedPercentageDiscount := unitPrice.Mul(decimal.NewFromFloat(0.02)).Mul(decimal.NewFromInt(3)) // 2% discount * 3 units
	discount, err = cartItem.calculateNonSubscriptionDiscount()
	assert.NoError(t, err, "Expected no error when calculating a non-subscription discount")
	assert.Equal(t, expectedPercentageDiscount.Round(2), discount.Round(2))
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
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(10.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(1.0),
				NumberOfUnitsDiscounted:        1,
				Subscription: SubscriptionInfo{
					BillingPeriodUnit: TimePeriodMonthly,
				},
			},
			wantErr: false,
		},
		{
			name: "Negative Quantity",
			item: &CartItem{
				Quantity:                       -1,
				UnitPrice:                      decimal.NewFromFloat(10.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(1.0),
				NumberOfUnitsDiscounted:        0,
			},
			wantErr: true,
			errMsg:  "quantity cannot be negative",
		},
		{
			name: "Negative Unit Price",
			item: &CartItem{
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(-10.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(1.0),
				NumberOfUnitsDiscounted:        1,
			},
			wantErr: true,
			errMsg:  "unit price cannot be negative",
		},
		{
			name: "Negative Discount Amount Per Unit",
			item: &CartItem{
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(10.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(-1.0),
				NumberOfUnitsDiscounted:        1,
			},
			wantErr: true,
			errMsg:  "discount amount per unit cannot be negative",
		},
		{
			name: "Negative Discounted Units Quantity",
			item: &CartItem{
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(10.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(1.0),
				NumberOfUnitsDiscounted:        -1,
			},
			wantErr: true,
			errMsg:  "discounted units quantity cannot be negative",
		},
		{
			name: "Discounted Units Quantity Exceeds Total Quantity",
			item: &CartItem{
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(10.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(1.0),
				NumberOfUnitsDiscounted:        2,
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
	//totalDiscountAmount, _ := decimal.NewFromString("1.0")
	//totalGrossAmount, _ := decimal.NewFromString("20.0")
	original := &CartItem{
		SkuID:                          "test-sku",
		Quantity:                       2,
		UnitPrice:                      unitPrice,
		DiscountValuePerDiscountedUnit: discountAmountPerUnit,
		NumberOfUnitsDiscounted:        1,
		//TotalDiscountAmount:            totalDiscountAmount,
		//TotalGrossAmount:               totalGrossAmount,
	}
	cloned := original.Clone()

	assert.Equal(t, original.SkuID, cloned.SkuID)
	assert.Equal(t, original.Quantity, cloned.Quantity)
	assert.True(t, original.UnitPrice.Equal(cloned.UnitPrice))
	assert.True(t, original.DiscountValuePerDiscountedUnit.Equal(cloned.DiscountValuePerDiscountedUnit))
	assert.Equal(t, original.NumberOfUnitsDiscounted, cloned.NumberOfUnitsDiscounted)
	//assert.True(t, original.TotalDiscountAmount.Equal(cloned.TotalDiscountAmount))
	//assert.True(t, original.TotalGrossAmount.Equal(cloned.TotalGrossAmount))
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
	cartItem.UnitPrice = decimal.NewFromInt(-1)
	assert.True(t, decimal.Zero.Equal(cartItem.GetGrossTotalAmount()), "GetGrossTotalAmount should return 0 when the unit price is negative")
}

func TestCartItem_GetNetTotalAmount(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10.00")
	discountAmountPerUnit, _ := decimal.NewFromString("1.00")
	cartItem := &CartItem{
		UnitPrice:                      unitPrice,
		Quantity:                       2,
		DiscountValuePerDiscountedUnit: discountAmountPerUnit,
		NumberOfUnitsDiscounted:        1,
		Subscription: SubscriptionInfo{
			BillingPeriodUnit: TimePeriodNoBilling,
		},
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
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(100.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(0),
				NumberOfUnitsDiscounted:        1,
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
				Quantity:                       0,
				UnitPrice:                      decimal.NewFromFloat(100.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(5),
				NumberOfUnitsDiscounted:        1,
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
				Quantity:                       1,
				UnitPrice:                      decimal.Zero,
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(5),
				NumberOfUnitsDiscounted:        1,
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
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(100.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(5),
				NumberOfUnitsDiscounted:        0,
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
				Quantity:                       2,
				UnitPrice:                      decimal.NewFromFloat(100.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(5),
				NumberOfUnitsDiscounted:        1,
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
				Quantity:                       1,
				UnitPrice:                      decimal.NewFromFloat(50.0),
				DiscountValuePerDiscountedUnit: decimal.NewFromFloat(100), // Intentionally high to trigger the cap
				NumberOfUnitsDiscounted:        1,
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
		IsSubscription: true,
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
		IsSubscription: true,
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

	//printLnDecimalToString(receivedDiscountAmount, "receivedDiscountAmount")
	discountAmountfloat, _ := discountAmount.Round(2).Float64()
	receivedDiscountAmountfloat, _ := receivedDiscountAmount.Round(2).Float64()
	assert.True(t, discountAmount.Round(2).Equal(receivedDiscountAmount.Round(2)), "The received discount amount should match the set value", discountAmountfloat, receivedDiscountAmountfloat)

}

func TestCartItem_GetSetTotalDiscountAmountValidPhaseDurationMatchesBilling(t *testing.T) {
	unitPrice, _ := decimal.NewFromString("10.00")
	cartItem := &CartItem{
		IsSubscription: true,
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

	//printLnDecimalToString(receivedDiscountAmount, "receivedDiscountAmount")
	discountAmountfloat, _ := discountAmount.Round(2).Float64()
	receivedDiscountAmountfloat, _ := receivedDiscountAmount.Round(2).Float64()
	assert.True(t, discountAmount.Round(2).Equal(receivedDiscountAmount.Round(2)), "The received discount amount should match the set value", discountAmountfloat, receivedDiscountAmountfloat)
}

func TestCartItem_GetSetTotalDiscountAmountValid2(t *testing.T) {
	cartItem := &CartItem{
		IsSubscription: true,
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

	//printLnDecimalToString(receivedDiscountAmount, "receivedDiscountAmount")
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
