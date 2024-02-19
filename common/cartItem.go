package common

import (
	"errors"
	"fmt"
	"math"

	"github.com/shopspring/decimal"
)

type DiscountType int

const (
	Unknown DiscountType = iota
	Percentage
	TimeBased
	FixedAmount
)

type DiscountApplication string

const (
	UnknownApplication DiscountApplication = "unknown"
	Recurring          DiscountApplication = "recurring"
	Spread             DiscountApplication = "spread"
)

type DiscountPhase struct {
	Duration                        int64               `json:"duration"` // Duration in units of TimeUnit
	DurationUnit                    TimeUnit            `json:"durationUnit"`
	DiscountValue                   float64             `json:"discountValue"` // Could be a percentage or a fixed amount
	DiscountType                    DiscountType        `json:"discountType"`  // Indicates Percentage, TimeBased, or FixedAmount
	Description                     string              `json:"description"`
	ApplicableNumberOfBillingCycles int64               `json:"applicateNumberOfBillingCycles"` // numsber of billing cycles this is applicable to
	Application                     DiscountApplication `json:"application"`                    // Recurring or Spread
}

func (dp *DiscountPhase) GetApplication() DiscountApplication {
	return dp.Application
}

func (dp *DiscountPhase) SetApplication(application DiscountApplication) {
	dp.Application = application
}

func (dp *DiscountPhase) GetDiscountValue() float64 {
	return dp.DiscountValue
}

func (dp *DiscountPhase) SetDiscountValue(value float64) {
	dp.DiscountValue = value
}

func (dp *DiscountPhase) GetDurationUnit() TimeUnit {
	return dp.DurationUnit
}

func (dp *DiscountPhase) SetDurationUnit(unit TimeUnit) {
	dp.DurationUnit = unit
}

func (dp *DiscountPhase) GetApplicableNumberOfBillingCycles() int64 {
	return dp.ApplicableNumberOfBillingCycles
}

func (dp *DiscountPhase) SetApplicableNumberOfBillingCycles(cycles int64) {
	dp.ApplicableNumberOfBillingCycles = cycles
}

func (dp *DiscountPhase) GetDescription() string {
	return dp.Description
}

func (dp *DiscountPhase) SetDescription(description string) {
	dp.Description = description
}

func (dp *DiscountPhase) GetDuration() int64 {
	return dp.Duration
}

func (dp *DiscountPhase) SetDuration(duration int64, durationUnit TimeUnit) {
	dp.Duration = duration
	dp.DurationUnit = durationUnit
}

func (dp *DiscountPhase) GetDiscountType() DiscountType {
	return dp.DiscountType
}

func (dp *DiscountPhase) ValidatePercentageDiscount() error {
	if dp.DiscountValue <= 0 {
		return errors.New("percentage discount rate must be positive")
	}
	if dp.DiscountValue > 100 {
		return errors.New("percentage discount rate cannot exceed 100")
	}
	if dp.DurationUnit == TimePeriodUnknown {
		return errors.New("invalid duration unit for percentage discount")
	}
	if dp.Duration <= 0 {
		return errors.New("invalid duration for percentage discount")
	}
	if dp.Application == UnknownApplication {
		return errors.New("invalid application for percentage discount")
	}
	//if dp.ApplicableNumberOfBillingCycles <= 0 {
	//	return errors.New("invalid number of applicable billing cycles for percentage discount")
	//}
	return nil
}

func (dp *DiscountPhase) ValidateTimeBasedDiscount() error {
	if dp.Duration <= 0 {
		return errors.New("time-based discount must have valid duration")
	}
	if dp.DurationUnit == TimePeriodUnknown {
		return errors.New("invalid duration unit for time-based discount")
	}
	if dp.Application == UnknownApplication {
		return errors.New("invalid application for time-based discount")
	}
	//if dp.ApplicableNumberOfBillingCycles <= 0 {
	//	return errors.New("invalid number of applicable billing cycles for time-based discount")
	//}
	return nil
}

// SetDiscountType sets the discount type and resets the fields based on the discount type.
func (dp *DiscountPhase) SetDiscountType(discountType DiscountType) error {
	// validate the fields based on the discount type
	// Percentage, TimeBased, or FixedAmount
	switch discountType {
	case Percentage:
		return dp.ValidatePercentageDiscount()
	case TimeBased:
		dp.DiscountValue = 0
		dp.DurationUnit = TimePeriodUnknown
		dp.Application = UnknownApplication
		dp.ApplicableNumberOfBillingCycles = 0
	case FixedAmount:
		dp.DiscountValue = 0
		dp.DurationUnit = TimePeriodUnknown
		dp.Duration = 0
		dp.Application = UnknownApplication
		dp.ApplicableNumberOfBillingCycles = 0
	}
	dp.DiscountType = discountType
	return nil
}

func (dp *DiscountPhase) ValidateAndReset() error {
	switch dp.DiscountType {
	case Percentage:
		if dp.DiscountValue <= 0 {
			return errors.New("percentage discount rate must be positive")
		}
		// Reset fields not related to percentage-based discounts
		dp.DiscountValue = 0
	case TimeBased:
		if dp.Duration <= 0 || dp.DurationUnit == TimePeriodUnknown {
			return errors.New("time-based discount must have valid duration and unit")
		}
		// Reset fields not related to time-based discounts
		dp.DiscountValue = 0
	case FixedAmount:
		if dp.DiscountValue <= 0 {
			return errors.New("fixed amount discount must be positive")
		}
		// Reset fields not related to fixed amount discounts
		dp.DiscountValue = 0
	default:
		return errors.New("invalid discount type")
	}
	return nil
}

type SubscriptionInfo struct {
	IsRecurring       bool            `json:"isRecurring"`
	BillingPeriodUnit TimeUnit        `json:"billingCyclePeriod"`
	TrialPeriod       int64           `json:"trialPeriod"` // Represented in units of TimeUnit, 0 if no trial
	TrialPeriodUnit   TimeUnit        `json:"trialPeriodUnit"`
	DiscountPhases    []DiscountPhase `json:"discountPhases"`
}

// CartItem represents an item in a shopping cart.
// It includes the SKU ID, quantity, unit price, and discount information.
// The discount information includes the discount amount per discounted unit, the number of discounted units,
// the total discount amount, and the total gross amount.

type CartItem struct {
	SkuID                           string           `json:"skuId"`
	Quantity                        int64            `json:"quantity"`
	UnitPrice                       decimal.Decimal  `json:"unitPrice"`                       // Updated to use decimal.Decimal
	DiscountAmountPerDiscountedUnit decimal.Decimal  `json:"discountAmountPerDiscountedUnit"` // Updated
	DiscountedUnitsQuantity         int64            `json:"discountedUnitsQuantity"`
	DiscountDescription             string           `json:"discountDescription"`
	TotalDiscountAmount             decimal.Decimal  `json:"totalDiscountAmount"` // Updated
	TotalGrossAmount                decimal.Decimal  `json:"totalGrossAmount"`    // Updated
	Subscription                    SubscriptionInfo `json:"subscription"`
}

/*
func (c *CartItem) CalculateTotalDiscountAmount() float64 {
	// This method needs to be updated to calculate the discount based on the DiscountPhases
	// For simplicity, this example will not implement the full logic, but you'll want to:
	// 1. Calculate any trial period discounts
	// 2. Apply any DiscountPhases
	// 3. Adjust the total discount amount accordingly
	return 0 // Placeholder return
}

func (c *CartItem) UpdateTotals() {
	// Update this method to consider the subscription phases in the total calculations
	c.TotalDiscountAmount = c.CalculateTotalDiscountAmount()
	// Calculate the gross amount, potentially adjusting for subscription phases
	c.TotalGrossAmount = float64(c.Quantity) * c.UnitPrice // Simplified; adjust as needed
}
*/

// add methods for adding/removing DiscountPhases, handling trials, etc.

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

func (c *CartItem) GetDiscountDescription() string {
	return c.DiscountDescription
}

func (c *CartItem) SetDiscountDescription(description string) {
	c.DiscountDescription = description
}

// create the functions to work with CartItem
func (c *CartItem) GetSkuID() string {
	return c.SkuID
}

func (c *CartItem) SetSkuID(skuID string) error {
	if skuID == "" {
		return errors.New("skuID cannot be empty")
	}
	c.SkuID = skuID
	return nil
}

func (c *CartItem) GetQuantity() int64 {
	return c.Quantity
}

func (c *CartItem) UpdateTotals() error {
	_, err := c.GetNetTotalAmount()
	return err
}

func (c *CartItem) SetQuantity(quantity int64) error {
	if quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	if c.Subscription.IsRecurring && quantity > 1 {
		return errors.New("quantity cannot exceed 1 for recurring subscriptions")
	}
	c.Quantity = quantity
	c.UpdateTotals()
	return nil
}

// GetUnitPrice returns the unit price as a decimal.Decimal.
func (c *CartItem) GetUnitPrice() decimal.Decimal {
	return c.UnitPrice
}

// SetUnitPrice sets the unit price from a decimal.Decimal value.
func (c *CartItem) SetUnitPrice(unitPrice decimal.Decimal) error {
	if unitPrice.IsNegative() {
		return errors.New("unit price cannot be negative")
	}
	c.UnitPrice = unitPrice
	return c.UpdateTotals()
}

// Optionally, if you want to allow setting the unit price from a string (e.g., when parsing user input):
func (c *CartItem) SetUnitPriceFromString(unitPriceStr string) error {
	unitPrice, err := decimal.NewFromString(unitPriceStr)
	if err != nil {
		return err // The error will indicate the invalid decimal format
	}
	return c.SetUnitPrice(unitPrice)
}

func (c *CartItem) GetGrossTotalAmount() decimal.Decimal {
	return c.UnitPrice.Mul(decimal.NewFromInt(int64(c.Quantity)))
}

func (c *CartItem) GetNetTotalAmount() (decimal.Decimal, error) {
	totalDiscountAmount, err := c.calculateTotalDiscountAmount()
	if err != nil {
		return decimal.Zero, err
	}
	return c.GetGrossTotalAmount().Sub(totalDiscountAmount), nil
}

func (c *CartItem) GetTotalDiscountAmount() (decimal.Decimal, error) {
	err := c.Validate()
	if err != nil {
		return decimal.Zero, err
	}
	return c.calculateTotalDiscountAmount()
}

func (c *CartItem) calculateTotalDiscountAmount() (decimal.Decimal, error) {
	if c.Quantity == 0 || c.UnitPrice.IsZero() {
		return decimal.Zero, nil
	}

	if c.Subscription.IsRecurring {
		return c.calculateTotalSubscriptionDiscount()
	} else {
		return c.calculateNonSubscriptionDiscount()
	}
}

func (c *CartItem) calculateTotalDiscountValue(phase DiscountPhase) (decimal.Decimal, error) {
	unitPrice := c.UnitPrice
	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	var discountPerCycle decimal.Decimal
	switch phase.DiscountType {
	case Percentage:
		if phase.DiscountValue > 100 {
			return decimal.Zero, errors.New("percentage discount value cannot exceed 100")
		}
		percentageRate := decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))
		discountPerCycle = unitPrice.Mul(percentageRate)
	default:
		return decimal.Zero, errors.New("unsupported discount type")
	}

	// Calculate the total discount value
	totalDiscount := discountPerCycle
	if normalizedDurationRatio > 1 {
		// For discounts that span multiple billing cycles
		cycleCount := math.Ceil(normalizedDurationRatio) // Ensure we cover partial cycles as full ones
		totalDiscount = discountPerCycle.Mul(decimal.NewFromFloat(cycleCount))
	} else {
		// For discounts within a single billing cycle, apply the discount ratio directly
		totalDiscount = discountPerCycle.Mul(decimal.NewFromFloat(normalizedDurationRatio))
	}

	return totalDiscount, nil
}

func (c *CartItem) calculateTotalSubscriptionDiscount() (decimal.Decimal, error) {
	totalDiscount := decimal.Zero // Assuming this starts from a base value, adjust as necessary

	if c.Subscription.DiscountPhases == nil || len(c.Subscription.DiscountPhases) == 0 {
		return totalDiscount, errors.New("no discount phases found; discount phases are required for recurring subscriptions")
	}

	for _, phase := range c.Subscription.DiscountPhases {
		phaseDiscount, err := c.calculateDiscountForPhase(phase)
		if err != nil {
			// Instead of continuing, we return on the first error encountered according to your revised approach
			return decimal.Zero, fmt.Errorf("error calculating discount for phase: %v", err)
		}

		// Accumulate the phase discount into the total discount
		totalDiscount = totalDiscount.Add(phaseDiscount)
	}

	// Cap the total discount to not exceed the total unit price, if applicable
	// Assuming there's a mechanism to calculate or retrieve the total unit price or gross amount
	// totalUnitPrice := ... // Define how to obtain this
	// if totalDiscount.GreaterThan(totalUnitPrice) {
	//     totalDiscount = totalUnitPrice
	// }

	return totalDiscount, nil
}

/*
func (c *CartItem) calculateTotalSubscriptionDiscount() (decimal.Decimal, error) {
	totalTrialDiscount := c.calculateTotalTrialPeriodDiscount()
	var totalSubscriptionDiscount decimal.Decimal
	printLnDecimalToString(totalTrialDiscount, "totalTrialDiscount")
	if c.Subscription.DiscountPhases == nil || len(c.Subscription.DiscountPhases) == 0 {
		return totalTrialDiscount, errors.New("no discount phases found and its required for recurring subscriptions")
	}
	for _, phase := range c.Subscription.DiscountPhases {
		phaseDiscount, err := c.calculatePerPhaseDiscount(phase)
		printLnDecimalToString(phaseDiscount, "phaseDiscount")
		// Multiply the phase discount by the number of billing cycles (duration of the phase in billing periods)
		if phase.DurationUnit < c.Subscription.BillingPeriodUnit {
			// does the phase duration unit multiplied by the phase duration result in a time period greater than the billing period unit?
			normalizedRation, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
			if err != nil {
				return decimal.Zero, err
			}
			// if the ratio is one then we
			// we need to make sure that the phase duration unit is the same as the billing period unit
			//
			if normalizedRation < 1 {
				totalSubscriptionDiscount = phaseDiscount
			} else if normalizedRation == 1 {
				totalSubscriptionDiscount = phaseDiscount.Mul(decimal.NewFromInt(int64(phase.Duration)))
			} else if normalizedRation > 1 {
				totalSubscriptionDiscount = phaseDiscount.Mul(decimal.NewFromFloat(normalizedRation))
			}

		}
		//phaseDiscount = phaseDiscount.Mul(decimal.NewFromInt(int64(phase.Duration)))
		printLnDecimalToString(phaseDiscount, "phaseDiscount") //phaseDiscount is 75
		if err != nil {
			return decimal.Zero, fmt.Errorf("error calculating phase %s discount: %v", phase.Description, err)
		}
		printLnDecimalToString(totalTrialDiscount, "totalTrialDiscount") //totalTrialDiscount is 0
		totalTrialDiscount = totalTrialDiscount.Round(2)
		printLnDecimalToString(totalTrialDiscount, "totalTrialDiscount") //totalTrialDiscount is 0
		totalSubscriptionDiscount = totalSubscriptionDiscount.Round(2).Add(totalTrialDiscount.Add(phaseDiscount.Round(2)))
		printLnDecimalToString(totalTrialDiscount, "totalTrialDiscount") //WHY IS totalTrialDiscount 75?
	}

	// Cap the total discount to not exceed the total gross amount

	return totalSubscriptionDiscount, nil
}
*/

func (c *CartItem) calculateTotalTrialPeriodDiscount() decimal.Decimal {
	if c.Subscription.TrialPeriod > 0 {
		normalizedTrialDuration, err := NormalizeDuration(int64(c.Subscription.TrialPeriod), c.Subscription.TrialPeriodUnit, c.Subscription.BillingPeriodUnit)
		if err != nil {
			return decimal.Zero
		}
		return c.UnitPrice.Mul(decimal.NewFromInt(int64(normalizedTrialDuration)))
	}
	return decimal.Zero
}

func (c *CartItem) calculateDiscountForPhase(phase DiscountPhase) (decimal.Decimal, error) {
	unitPrice := c.UnitPrice // Assuming this is the price for a single billing cycle
	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	var discount, percentageRate decimal.Decimal

	switch phase.DiscountType {
	case TimeBased:
		// set the discountValue to 100% if the discount type is time based and fallthrough to the percentage case
		phase.DiscountValue = 100
		fallthrough
	case Percentage:
		// Calculate the discount as a percentage of the unit price.
		percentageRate = decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))
		discount = unitPrice.Mul(percentageRate)

		if phase.Application == Recurring {
			// Prorate the discount based on the normalized duration ratio for the first billing cycle.
			discount = discount.Mul(decimal.NewFromFloat(normalizedDurationRatio))
			printLnDecimalToString(discount, "discount")
			if normalizedDurationRatio > 1 {
				// validate that the phase duration unit is less than or equal to the billing period unit
				// for recurring, we want to apply a discount for each billing cycle so if normalizedDurationRatio is greater than 1
				// then the period unit must be less than or equal to the billing period unit so the discount is applied to each billing cycle
				if phase.DurationUnit != c.Subscription.BillingPeriodUnit && phase.DurationUnit.HourValue()*phase.Duration > c.Subscription.BillingPeriodUnit.HourValue() {
					return decimal.Zero, errors.New("phase duration unit period cannot be greater than the billing period unit for recurring discounts")
				}
				ValidNumberOfBilingCycles, _ := decimal.NewFromInt(phase.ApplicableNumberOfBillingCycles).Round(2).Float64()
				if ValidNumberOfBilingCycles < normalizedDurationRatio {
					// cap the normalizedDurationRatio to the number of applicable billing cycles
					normalizedDurationRatio = ValidNumberOfBilingCycles
				}
				//additionalDiscount := unitPrice.Mul(percentageRate).Mul(decimal.NewFromFloat(additionalCycles))
				discount = discount.Mul(decimal.NewFromFloat(normalizedDurationRatio)).Mul(percentageRate)
				printLnDecimalToString(discount, "discount")
			}
		} else if phase.Application == Spread || phase.Application == UnknownApplication || phase.Application == "" {
			// Spread the discount evenly across the duration of the phase.
			discount = discount.Mul(decimal.NewFromFloat(normalizedDurationRatio))
			if normalizedDurationRatio > 1 {

				ValidNumberOfBilingCycles, _ := decimal.NewFromInt(phase.ApplicableNumberOfBillingCycles).Round(2).Float64()
				var remainingCycles float64
				if ValidNumberOfBilingCycles < normalizedDurationRatio {
					// cap the normalizedDurationRatio to the number of applicable billing cycles
					remainingCycles = normalizedDurationRatio - ValidNumberOfBilingCycles
					normalizedDurationRatio = ValidNumberOfBilingCycles
				}
				discountAdjustment := unitPrice.Mul(percentageRate).Mul(decimal.NewFromFloat(remainingCycles))
				discount = discount.Sub(discountAdjustment)
			}
		}

	case FixedAmount:
		// Apply the fixed amount directly, prorated by the duration ratio if it's less than a full cycle.
		fixedAmountDiscount := decimal.NewFromFloat(phase.DiscountValue)
		if phase.Application == Recurring {
			if normalizedDurationRatio > 1 {
				// validate that the phase duration unit is less than or equal to the billing period unit
				// for recurring, we want to apply a discount for each billing cycle so if normalizedDurationRatio is greater than 1
				// then the period unit must be less than or equal to the billing period unit so the discount is applied to each billing cycle
				if phase.DurationUnit != c.Subscription.BillingPeriodUnit && phase.DurationUnit.HourValue()*phase.Duration > c.Subscription.BillingPeriodUnit.HourValue() {
					return decimal.Zero, errors.New("phase duration unit cannot be greater than the billing period unit for recurring discounts")
				}
				additionalCycles := math.Floor(normalizedDurationRatio) - 1
				if phase.ApplicableNumberOfBillingCycles > 0 {
					// Cap the additional cycles to the specified number of billing cycles
					additionalCycles = math.Min(additionalCycles, float64(phase.ApplicableNumberOfBillingCycles))
				}
				additionalDiscount := fixedAmountDiscount.Mul(decimal.NewFromFloat(additionalCycles))
				discount = discount.Add(additionalDiscount)
			}
		} else if phase.Application == Spread {
			if normalizedDurationRatio > 1 {
				// If the phase spans multiple cycles and is spread, apply the full discount to each billing cycle and prorate the rest.
				additionalCycles := math.Floor(normalizedDurationRatio) - 1
				if phase.ApplicableNumberOfBillingCycles > 0 {
					// Cap the additional cycles to the specified number of billing cycles
					additionalCycles = math.Min(additionalCycles, float64(phase.ApplicableNumberOfBillingCycles))
				}
				additionalDiscount := fixedAmountDiscount.Mul(decimal.NewFromFloat(additionalCycles))
				discount = discount.Add(additionalDiscount)
			}
		}
	default:
		return decimal.Zero, errors.New("unsupported discount type")
	}

	return discount, nil
}

/*
func (c *CartItem) calculateDiscountForPhase(phase DiscountPhase) (decimal.Decimal, error) {
	unitPrice := c.UnitPrice // Assuming this is the price for a single billing cycle
	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	cycleCount := math.Ceil(normalizedDurationRatio) // Number of billing cycles covered

	var totalDiscount decimal.Decimal
	switch phase.DiscountType {
	case Percentage:
		percentageRate := decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))
		discountPerCycle := unitPrice.Mul(percentageRate)
		printLnDecimalToString(discountPerCycle, "discountPerCycle")
		if phase.Application == Recurring {
			totalDiscount = discountPerCycle.Mul(decimal.NewFromFloat(cycleCount))
		} else { // Spread
			totalDiscount = discountPerCycle // Apply once and spread the effect if needed
		}

	case TimeBased:
		// For simplicity, assuming time-based discounts apply a fixed rate for the duration
		// This case might need special handling based on how time-based discounts are defined in your system
		return decimal.Zero, errors.New("time-based discounts require specific handling")

	case FixedAmount:
		fixedAmountDiscount := decimal.NewFromFloat(phase.DiscountValue)
		if phase.Application == Recurring {
			totalDiscount = fixedAmountDiscount.Mul(decimal.NewFromFloat(cycleCount))
		} else { // Spread
			// Spread the total fixed amount evenly across the cycleCount
			totalDiscount = fixedAmountDiscount.Div(decimal.NewFromFloat(cycleCount)).Mul(decimal.NewFromFloat(cycleCount))
		}

	default:
		return decimal.Zero, errors.New("unsupported discount type")
	}

	// Adjust for partial cycles if needed, especially for spread discounts
	if phase.Application == Spread && normalizedDurationRatio < 1 {
		totalDiscount = totalDiscount.Mul(decimal.NewFromFloat(normalizedDurationRatio))
	}

	return totalDiscount, nil
}
*/
/*
func (c *CartItem) calculateDiscountForPhase(phase DiscountPhase) (decimal.Decimal, error) {
	unitPrice := c.UnitPrice
	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	cycleCount := math.Ceil(normalizedDurationRatio) // Number of billing cycles covered

	var totalDiscount decimal.Decimal
	switch phase.Application {
	case Recurring:
		// For recurring discounts, apply the discount rate for each cycle up to n
		percentageRate := decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))
		discountPerCycle := unitPrice.Mul(percentageRate)
		totalDiscount = discountPerCycle.Mul(decimal.NewFromFloat(cycleCount))

	case Spread:
		// For spread discounts, divide the total discount amount evenly across n cycles
		totalDiscountAmount := decimal.NewFromFloat(phase.DiscountValue) // Assuming this is the total amount to spread
		discountPerCycle := totalDiscountAmount.Div(decimal.NewFromFloat(cycleCount))
		totalDiscount = discountPerCycle.Mul(decimal.NewFromFloat(cycleCount))

	default:
		return decimal.Zero, errors.New("unsupported discount application type")
	}

	return totalDiscount, nil
}
*/
/*
func (c *CartItem) calculateDiscountForPhase(phase DiscountPhase) (decimal.Decimal, error) {
	unitPrice := c.UnitPrice // Assuming UnitPrice is already a decimal.Decimal
	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	var discount decimal.Decimal
	switch phase.DiscountType {
	case Percentage:
		if phase.DiscountValue > 100 {
			return decimal.Zero, errors.New("percentage discount value cannot exceed 100")
		}
		percentageRate := decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))

		if normalizedDurationRatio <= 1 {
			// For discounts covering less than or exactly one billing cycle
			discount = unitPrice.Mul(percentageRate).Mul(decimal.NewFromFloat(normalizedDurationRatio))
		} else {
			// For discounts spanning multiple billing cycles, distribute the discount evenly across the cycles
			// Here, you may decide to apply the full discount rate to each cycle or adjust the discount rate
			// based on the number of cycles covered. This example applies the full rate to each cycle.
			discount = unitPrice.Mul(percentageRate) // Apply the discount rate to the unit price for each cycle
			// Optionally, adjust the discount based on the total number of cycles if needed
		}
	default:
		// Handle other types as needed
		return decimal.Zero, errors.New("unsupported discount type")
	}

	// Ensure the discount does not exceed the unit price per billing cycle
	if discount.GreaterThan(unitPrice) {
		discount = unitPrice
	}

	return discount, nil
}
*/
/*
func (c *CartItem) calculatePerPhaseDiscount(phase DiscountPhase) (decimal.Decimal, error) {
	unitPrice := c.UnitPrice
	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	// Normalize the phase duration to the billing period unit
	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	var discount decimal.Decimal
	switch phase.DiscountType {
	case Percentage:
		if phase.DiscountValue > 100 {
			return decimal.Zero, errors.New("percentage discount value cannot exceed 100")
		}
		percentageRate := decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))

		// Apply the discount rate to the portion of the billing period covered by the phase
		discount = unitPrice.Mul(percentageRate).Mul(decimal.NewFromFloat(normalizedDurationRatio))

	default:
		// For simplicity, other types are not detailed. Implement according to your requirements.
		return decimal.Zero, errors.New("unsupported discount type")
	}

	// Ensure the discount does not exceed the unit price for partial billing periods
	if discount.GreaterThan(unitPrice) {
		discount = unitPrice
	}

	return discount, nil
}
*/

func (c *CartItem) calculatePerPhaseDiscountOld(phase DiscountPhase) (decimal.Decimal, error) {

	unitPrice := c.UnitPrice // Assuming UnitPrice is already a decimal.Decimal

	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	// Default to the billing period unit if the phase duration unit is not set
	if phase.DurationUnit == TimePeriodUnknown {
		if phase.Duration > 0 {
			return decimal.Zero, errors.New("invalid phase DurationUnit. Please provide a valid duration unit")
		}
	}

	// Normalize the phase duration to the billing period unit
	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	// Return an error if the phase period unit is greater than the billing period unit
	// This is inferred from a normalizedDurationRatio greater than 1, meaning the phase covers more than one billing cycle
	if normalizedDurationRatio > 1 {
		return decimal.Zero, errors.New("phase duration unit cannot be greater than the billing period unit")
	}

	var discount decimal.Decimal
	switch phase.DiscountType {
	case Percentage:
		// Validate percentage discount value
		if phase.DiscountValue > 100 {
			return decimal.Zero, errors.New("percentage discount value cannot exceed 100")
		}
		// Calculate discount for the billing cycle as a percentage of the unit price
		percentageRate := decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))
		printLnDecimalToString(percentageRate, "percentageRate")
		discount = unitPrice.Mul(percentageRate)
		printLnDecimalToString(discount, "discount")
		discount = discount.Mul(decimal.NewFromFloat(normalizedDurationRatio))
		printLnDecimalToString(discount, "discount")
		/*
			// Calculate the discount based on the normalized duration as a percentage of the billing period.
			printLnDecimalToString(unitPrice, "unitPrice")

			discountPerDurationUnit := unitPrice.Mul(percentageRate).Mul(decimal.NewFromFloat(normalizedPhaseDuration))
			printLnDecimalToString(discountPerDurationUnit, "discountPerDurationUnit")
			// if the phase duration time unit is less that the billing period unit, we need to adjust the discount
			if phase.DurationUnit < c.Subscription.BillingPeriodUnit {
				discount = discountPerDurationUnit
			} else {
				// this returns the whole value of the discount because the phase duration unit is greater than the billing period unit
				// if the phase duration time unit is greater than the billing period unit, we need to adjust the discount
				//should this return an error because why would the phase duration unit be greater than the billing period unit
				discount = discountPerDurationUnit.Div(decimal.NewFromFloat(float64(normalizedPhaseDuration)))
			}
			//discount = discountPerDurationUnit.Mul(decimal.NewFromInt(phase.Duration))
			printLnDecimalToString(discount, "discount")
		*/
	case TimeBased:
		// Handling for TimeBased discounts should be defined based on specific rules
		//return decimal.Zero, errors.New("time-based discounts not implemented in this context")
		discount = unitPrice.Mul(decimal.NewFromFloat(normalizedDurationRatio))
	default:
		// For direct discounts or other types
		discount = decimal.NewFromFloat(phase.DiscountValue)
		if discount.GreaterThan(unitPrice) {
			return decimal.Zero, errors.New("discount value cannot exceed unit price")
		}
	}

	// Ensure the discount does not exceed the unit price
	if discount.GreaterThan(unitPrice) {
		discount = unitPrice
	}

	return discount, nil
}

func (c *CartItem) calculateNonSubscriptionDiscount() (decimal.Decimal, error) {
	if c.DiscountedUnitsQuantity == 0 || c.DiscountAmountPerDiscountedUnit.IsZero() {
		return decimal.Zero, nil
	}

	discountedAmount := c.DiscountAmountPerDiscountedUnit.Mul(decimal.NewFromInt(int64(c.DiscountedUnitsQuantity)))
	// Cap the discount to not exceed the total gross amount
	grossTotal := c.GetGrossTotalAmount()
	if discountedAmount.GreaterThan(grossTotal) {
		return grossTotal, nil
	}
	return discountedAmount, nil
}

func (c *CartItem) SetDiscountAmountPerUnit(discountAmount decimal.Decimal) error {
	if discountAmount.IsNegative() {
		return errors.New("discount amount cannot be negative")
	}
	c.DiscountAmountPerDiscountedUnit = discountAmount
	c.UpdateTotals()
	return nil
}

func (c *CartItem) SetDiscountAmountPerUnitFromString(discountAmountStr string) error {
	discountAmount, err := decimal.NewFromString(discountAmountStr)
	if err != nil {
		return errors.New("invalid discount amount format")
	}
	if discountAmount.IsNegative() {
		return errors.New("discount amount cannot be negative")
	}
	c.DiscountAmountPerDiscountedUnit = discountAmount
	c.UpdateTotals()
	return nil
}

func (c *CartItem) GetDiscountedUnitsQuantity() int64 {
	return c.DiscountedUnitsQuantity
}

func (c *CartItem) SetDiscountedUnitQuantity(discountedUnitsQty int64) error {
	if discountedUnitsQty < 0 {
		return errors.New("discounted units quantity cannot be negative")
	}
	c.DiscountedUnitsQuantity = discountedUnitsQty
	c.UpdateTotals()
	return nil
}

func (c *CartItem) IncrementCartItemQuantity(quantity int64) (int64, error) {
	if quantity < 0 {
		return c.Quantity, errors.New("quantity cannot be negative")
	}
	c.Quantity += quantity
	c.UpdateTotals()
	return c.Quantity, nil
}

func (c *CartItem) DecrementCartItemQuantity(quantity int64) (int64, error) {
	if c.Quantity-quantity < 0 || quantity < 0 {
		return c.Quantity, errors.New("quantity cannot be less than 0")
	}
	c.Quantity -= quantity
	c.UpdateTotals()
	return c.Quantity, nil
}

func (c *CartItem) Validate() error {
	if c.Quantity < 0 {
		return errors.New("quantity cannot be negative")
	}
	if c.UnitPrice.IsNegative() {
		return errors.New("unit price cannot be negative")
	}
	if c.DiscountAmountPerDiscountedUnit.IsNegative() {
		return errors.New("discount amount per unit cannot be negative")
	}
	if c.DiscountedUnitsQuantity < 0 {
		return errors.New("discounted units quantity cannot be negative")
	}
	if c.DiscountedUnitsQuantity > c.Quantity {
		return errors.New("discounted units quantity cannot exceed total quantity")
	}
	if c.Subscription.BillingPeriodUnit == TimePeriodUnknown {
		return errors.New("invalid billing period unit")
	}
	// range through the discount phases and validate them against the billing period unit
	for _, phase := range c.Subscription.DiscountPhases {
		if phase.DurationUnit == TimePeriodUnknown {
			return errors.New("invalid duration unit for discount phase")
		}
		if phase.Duration <= 0 {
			return errors.New("invalid duration for discount phase")
		}
		if phase.DurationUnit > c.Subscription.BillingPeriodUnit {
			return errors.New("discount phase duration unit cannot exceed billing period unit")
		}
	}
	return nil
}
