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

func (DiscountType DiscountType) String() string {
	switch DiscountType {
	case Percentage:
		return "Percentage"
	case TimeBased:
		return "TimeBased"
	case FixedAmount:
		return "FixedAmount"
	default:
		return "Unknown"
	}
}

type DiscountApplication string

const (
	UnknownApplication DiscountApplication = "unknown"
	Recurring          DiscountApplication = "recurring"
	Spread             DiscountApplication = "spread"
	OneTime            DiscountApplication = "one-time"
)

type DiscountPhase struct {
	Duration                        int64               `json:"duration"` // Duration in units of TimeUnit
	DurationUnit                    TimeUnit            `json:"durationUnit"`
	DiscountValue                   float64             `json:"discountValue"` // Could be a percentage or a fixed amount
	DiscountType                    DiscountType        `json:"discountType"`  // Indicates Percentage, TimeBased, or FixedAmount
	Description                     string              `json:"description"`
	ApplicableNumberOfBillingCycles int64               `json:"applicateNumberOfBillingCycles"` // numsber of billing cycles this is applicable to
	Application                     DiscountApplication `json:"application"`                    // Recurring or Spread
	Logs                            []string            `json:"logs"`
	DiscountsPerBillingCycle        map[int64]float64   `json:"discountPerBillingCycle"`
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
		phaseDiscount, err := c.calculateDiscountForPhase(&phase)
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

func (c *CartItem) calculateDiscountForPhase(phase *DiscountPhase) (decimal.Decimal, error) {
	unitPrice := c.UnitPrice

	// Validate discount value
	if phase.DiscountValue <= 0 {
		return decimal.Zero, errors.New("discount value must be positive")
	}

	// correct the duration unit to the billing period unit if fixed amount
	if phase.DiscountType == FixedAmount {
		phase.DurationUnit = c.Subscription.BillingPeriodUnit
		phase.Duration = phase.ApplicableNumberOfBillingCycles
	}

	// Normalize duration
	normalizedDurationRatio, err := NormalizeDuration(phase.Duration, phase.DurationUnit, c.Subscription.BillingPeriodUnit)
	if err != nil {
		return decimal.Zero, err
	}

	// Log start of calculation
	log := []string{}

	log = append(log, "Start calculation:")
	log = append(log, fmt.Sprintf("Unit Price: %s", unitPrice.String()))
	log = append(log, fmt.Sprintf("Discount Type: %s", phase.DiscountType))
	log = append(log, fmt.Sprintf("Normalized Duration Ratio: %f", normalizedDurationRatio))

	// discounts per billing cycles
	discounts := make(map[int64]decimal.Decimal)

	// Calculate discount based on discount type
	var discount, percentageRate decimal.Decimal
	ValidNumberOfBilingCycles, _ := decimal.NewFromInt(phase.ApplicableNumberOfBillingCycles).Round(2).Float64()
	switch phase.DiscountType {
	case Percentage, TimeBased:
		// If time based, set the discount value to 100 to represent a full discount for the time period
		if phase.DiscountType == TimeBased {
			phase.DiscountValue = 100
			log = append(log, "Set Discount Value 100 because this is a TimeBased discount")
		}

		// Calculate the discount as a percentage of the unit price.
		percentageRate = decimal.NewFromFloat(phase.DiscountValue).Div(decimal.NewFromFloat(100))
		discount = unitPrice.Mul(percentageRate)
		discounts[1] = discount

		log = append(log, fmt.Sprintf("Discount Type: %s", phase.DiscountType))
		log = append(log, fmt.Sprintf("Percentage Rate: %s", percentageRate.String()))
		log = append(log, fmt.Sprintf("Discount: %s", discount.String()))

		if phase.Application == Recurring {
			// Prorate the discount based on the normalized duration ratio for the first billing cycle.
			discount = discount.Mul(decimal.NewFromFloat(normalizedDurationRatio))
			discounts[1] = discount
			log = append(log, fmt.Sprintf("Discount after proration: %s", discount.String()))

			if normalizedDurationRatio > 1 {
				// validate that the phase duration unit is less than or equal to the billing period unit
				// for recurring, we want to apply a discount for each billing cycle so if normalizedDurationRatio is greater than 1
				// then the period unit must be less than or equal to the billing period unit so the discount is applied to each billing cycle
				if phase.DurationUnit != c.Subscription.BillingPeriodUnit && phase.DurationUnit.HourValue()*phase.Duration > c.Subscription.BillingPeriodUnit.HourValue() {
					return decimal.Zero, errors.New("phase duration unit period cannot be greater than the billing period unit for recurring discounts")
				}

				if ValidNumberOfBilingCycles < normalizedDurationRatio {
					// cap the normalizedDurationRatio to the number of applicable billing cycles
					normalizedDurationRatio = ValidNumberOfBilingCycles
				}
				//additionalDiscount := unitPrice.Mul(percentageRate).Mul(decimal.NewFromFloat(additionalCycles))
				discount = discount.Mul(decimal.NewFromFloat(normalizedDurationRatio)).Mul(percentageRate)
				for i := 2; i <= int(normalizedDurationRatio); i++ {
					discounts[int64(i)] = discount
				}
				log = append(log, fmt.Sprintf("Discount after recurring discount: %s", discount.String()))
			}
		} else if phase.Application == Spread || phase.Application == UnknownApplication || phase.Application == "" {
			// Spread the discount evenly across the duration of the phase.
			discount = discount.Mul(decimal.NewFromFloat(normalizedDurationRatio))
			discounts[1] = discount
			log = append(log, fmt.Sprintf("Discount after spreading: %s", discount.String()))
			if normalizedDurationRatio > 1 {
				var remainingCycles float64
				if ValidNumberOfBilingCycles < normalizedDurationRatio {
					// cap the normalizedDurationRatio to the number of applicable billing cycles
					remainingCycles = normalizedDurationRatio - ValidNumberOfBilingCycles
					normalizedDurationRatio = ValidNumberOfBilingCycles
				}
				// discount any cycles that are beyond the allowed number of billing cycles
				discountAdjustment := unitPrice.Mul(percentageRate).Mul(decimal.NewFromFloat(remainingCycles))
				discount = discount.Sub(discountAdjustment)
				discounts[1] = discount
				log = append(log, fmt.Sprintf("Discount adjustment for remaining cycles: %s", discountAdjustment.String()))
				log = append(log, fmt.Sprintf("Discount after adjustment: %s", discount.String()))
				if unitPrice.LessThan(discount) && ValidNumberOfBilingCycles <= 1 {
					discount = unitPrice
					discounts[1] = discount
					log = append(log, fmt.Sprintf("Discount capped at unit price: %s", discount.String()))
				} else if unitPrice.LessThan(discount) && ValidNumberOfBilingCycles > 1 {
					// spread the discount over the number of billing cycles that apply
					howManyBillingCycles := discount.DivRound(unitPrice, 2)
					accountedDiscount := 0.0
					for i := 1; i <= int(howManyBillingCycles.IntPart()); i++ {
						discounts[int64(i)] = unitPrice //ensures maximum discount is the unit price
						tmpDiscount, _ := unitPrice.Round(2).Float64()
						accountedDiscount += tmpDiscount
						log = append(log, fmt.Sprintf("Discount capped at unit price: %s for billing cycle %d", discount.String(), i))
					}
					lastBillingCycle := int64(howManyBillingCycles.IntPart()) + 1
					remainingDiscountValue := discount.Sub(decimal.NewFromFloat(accountedDiscount))
					if !remainingDiscountValue.IsZero() {
						discounts[lastBillingCycle] = remainingDiscountValue
						log = append(log, fmt.Sprintf("Discount capped at unit price: %s for billing cycle %d", discount.String(), lastBillingCycle))
					}
				}

			}
		}

	case FixedAmount:
		if phase.Application != Recurring && phase.Application != OneTime {
			return decimal.Zero, errors.New("fixed amount discounts must be recurring or onetime")
		}
		// Apply the fixed amount directly, prorated by the duration ratio if it's less than a full cycle.
		fixedAmountDiscount := decimal.NewFromFloat(phase.DiscountValue)

		if fixedAmountDiscount.GreaterThan(unitPrice) {
			fixedAmountDiscount = unitPrice
			log = append(log, fmt.Sprintf("Discount capped at unit price: %s", fixedAmountDiscount.String()))
		}
		discount = fixedAmountDiscount
		log = append(log, fmt.Sprintf("Discount Type: %s", phase.DiscountType.String()))
		log = append(log, fmt.Sprintf("Fixed Amount Discount: %s", fixedAmountDiscount.String()))

		if phase.Application == Recurring {

			for i := 1; i <= int(phase.ApplicableNumberOfBillingCycles); i++ {
				discounts[int64(i)] = discount
			}
			discount = fixedAmountDiscount.Mul(decimal.NewFromInt(phase.ApplicableNumberOfBillingCycles))
			log = append(log, fmt.Sprintf("Discount after applying for recurring: %s", discount.String()))
		} else { // anything else, we assume a one time fixed amount discount
			discount = fixedAmountDiscount.Mul(decimal.NewFromInt(phase.ApplicableNumberOfBillingCycles))
			log = append(log, fmt.Sprintf("Discount after applying for one time fixed: %s (set the type to: %s)", discount.String(), phase.Application))
			discounts[1] = discount
			log = append(log, fmt.Sprintf("Discount after spreading: %s", discount.String()))
		}

	default:
		return decimal.Zero, errors.New("unsupported discount type")
	}

	// Log final discount value
	log = append(log, fmt.Sprintf("Final Discount: %s", discount.String()))

	// Set logs to c.Logs
	phase.Logs = log
	phase.DiscountsPerBillingCycle = make(map[int64]float64)
	// Set the discounts per billing cycle back to float64 for convenience
	for k, v := range discounts {
		discount, _ := v.Round(2).Float64()
		phase.DiscountsPerBillingCycle[k] = discount
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
