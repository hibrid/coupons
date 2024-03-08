package discount

import (
	"errors"

	"github.com/hibrid/coupons/common"
)

type DiscountConfig struct {
	Type            common.CouponDiscountType `json:"type"`
	Value           float64                   `json:"value"`
	MinimumPurchase float64                   `json:"minimumPurchase,omitempty"`
	DiscountedSku   string                    `json:"discountedSku,omitempty"`
	FreeSku         string                    `json:"freeSku,omitempty"`
	DiscountedQty   int                       `json:"discountedQty,omitempty"`
	FreeQty         int                       `json:"freeQty,omitempty"`
}

type DiscountType int

const (
	Unknown DiscountType = iota
	Percentage
	TimeBased
	FixedAmount
)

func (dt DiscountType) String() string {
	switch dt {
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
	Duration                        int64               `json:"duration"`                        // The length of time the discount applies, in units defined by DurationUnit
	DurationUnit                    common.TimeUnit     `json:"durationUnit"`                    // The unit of time for the duration (e.g., days, weeks, months)
	DiscountValue                   float64             `json:"discountValue"`                   // The value of the discount, which could represent a percentage off, a fixed amount, or other values depending on DiscountType
	DiscountType                    DiscountType        `json:"discountType"`                    // The nature of the discount (percentage, fixed amount, or time-based)
	Description                     string              `json:"description"`                     // A human-readable description of the discount phase
	ApplicableNumberOfBillingCycles int64               `json:"applicableNumberOfBillingCycles"` // For subscription-based discounts, the number of billing cycles the discount applies to
	Application                     DiscountApplication `json:"application"`                     // The scope of the discount's application (e.g., recurring, one-time)
	Logs                            []string            `json:"logs"`                            // Logs related to the discount application for auditing or debugging
	DiscountsPerBillingCycle        map[int64]float64   `json:"discountPerBillingCycle"`         // Specific discount amounts applied per billing cycle, if applicable
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

func (dp *DiscountPhase) GetDurationUnit() common.TimeUnit {
	return dp.DurationUnit
}

func (dp *DiscountPhase) SetDurationUnit(unit common.TimeUnit) {
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

func (dp *DiscountPhase) SetDuration(duration int64, durationUnit common.TimeUnit) {
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
	if dp.DurationUnit == common.TimePeriodUnknown {
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
	return dp.ValidatePercentageDiscount()
}

func (dp *DiscountPhase) ValidateFixedAmountDiscount() error {
	if dp.DiscountValue <= 0 {
		return errors.New("fixed amount discount must be positive")
	}
	if dp.DurationUnit != common.TimePeriodUnknown {
		return errors.New("invalid duration unit for fixed amount discount. This field should be TimePeriodUnknown")
	}
	if dp.Duration != 0 {
		return errors.New("invalid duration for fixed amount discount. This field should be 0")
	}
	if dp.Application != OneTime && dp.Application != Recurring {
		return errors.New("invalid application for fixed amount discount. Should be one-time or recurring")
	}
	if dp.ApplicableNumberOfBillingCycles <= 0 {
		return errors.New("invalid number of applicable billing cycles for fixed amount discount")
	}
	if dp.ApplicableNumberOfBillingCycles > 1 && dp.Application != Recurring {
		return errors.New("fixed amount discount can only be applied to one billing cycle for one-time application")
	}
	return nil
}

// SetDiscountType sets the discount type and resets the fields based on the discount type.
func (dp *DiscountPhase) SetDiscountType(discountType DiscountType) error {
	var err error
	// validate the fields based on the discount type
	// Percentage, TimeBased, or FixedAmount
	switch discountType {
	case Percentage, TimeBased:
		if discountType == TimeBased { // if time based, set the discount value to 100 to represent a full discount for the time period
			dp.DiscountValue = 100
		}
		err = dp.ValidatePercentageDiscount()
	case FixedAmount:
		err = dp.ValidateFixedAmountDiscount()
	}
	dp.DiscountType = discountType
	return err
}
