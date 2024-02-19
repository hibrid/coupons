package common

import (
	"errors"
)

type TimeUnit int

const (
	TimePeriodUnknown    TimeUnit = iota // Default value, represents an undefined billing period
	TimePeriodHourly                     // Represents an hourly billing period
	TimePeriodDaily                      // Represents a daily billing period
	TimePeriodWeekly                     // Represents a weekly billing period
	TimePeriodBiWeekly                   // Represents a bi-weekly billing period
	TimePeriodThirtyDays                 // Represents a 30-day billing period
	TimePeriodMonthly                    // Represents a monthly billing period
	TimePeriodQuarterly                  // Represents a quarterly billing period
	TimePeriodBiAnnual                   // Represents a bi-annual billing period
	TimePeriodAnnual                     // Represents an annual billing period
	TimePeriodBiennial                   // Represents a biennial billing period
	TimePeriodNoBilling                  // Represents no billing period
)

func (bp TimeUnit) IsValid() bool {
	return bp >= TimePeriodHourly && bp <= TimePeriodNoBilling
}

func (bp TimeUnit) HourValue() int64 {
	switch bp {
	case TimePeriodHourly:
		return 1
	case TimePeriodDaily:
		return 24
	case TimePeriodWeekly:
		return 24 * 7
	case TimePeriodBiWeekly:
		return 24 * 14
	case TimePeriodThirtyDays:
		return 24 * 30
	case TimePeriodMonthly:
		return 24 * 30
	case TimePeriodQuarterly:
		return 24 * 90
	case TimePeriodBiAnnual:
		return 24 * 365 / 2
	case TimePeriodAnnual:
		return 24 * 365
	case TimePeriodBiennial:
		return 24 * 730
	default:
		return 0
	}
}

func (bp TimeUnit) String() string {
	switch bp {
	case TimePeriodHourly:
		return "Hourly"
	case TimePeriodDaily:
		return "Daily"
	case TimePeriodWeekly:
		return "Weekly"
	case TimePeriodBiWeekly:
		return "BiWeekly"
	case TimePeriodThirtyDays:
		return "ThirtyDays"
	case TimePeriodMonthly:
		return "Monthly"
	case TimePeriodQuarterly:
		return "Quarterly"
	case TimePeriodBiAnnual:
		return "BiAnnual"
	case TimePeriodAnnual:
		return "Annual"
	case TimePeriodBiennial:
		return "Biennial"
	case TimePeriodNoBilling:
		return "NoBilling"
	default:
		return "Unknown"
	}
}

var stringToTimePeriod = map[string]TimeUnit{
	"Hourly":     TimePeriodHourly,
	"Daily":      TimePeriodDaily,
	"Weekly":     TimePeriodWeekly,
	"BiWeekly":   TimePeriodBiWeekly,
	"ThirtyDays": TimePeriodThirtyDays,
	"Monthly":    TimePeriodMonthly,
	"Quarterly":  TimePeriodQuarterly,
	"BiAnnual":   TimePeriodBiAnnual,
	"Annual":     TimePeriodAnnual,
	"Biennial":   TimePeriodBiennial,
	"NoBilling":  TimePeriodNoBilling,
}

func ConvertStringToTimePeriod(s string) (TimeUnit, error) {
	if bp, ok := stringToTimePeriod[s]; ok {
		return bp, nil
	}
	return TimePeriodUnknown, errors.New("invalid billing period")
}

func NormalizeDuration(phaseDuration int64, phaseFromUnit, billingToUnit TimeUnit) (float64, error) {
	var hoursPerUnit = map[TimeUnit]float64{
		TimePeriodHourly:     1,            // 1 hour
		TimePeriodDaily:      24,           // 24 hours in a day
		TimePeriodWeekly:     24 * 7,       // 168 hours in a week
		TimePeriodBiWeekly:   24 * 14,      // 336 hours in two weeks
		TimePeriodThirtyDays: 24 * 30,      // 720 hours, assuming 30 days per month
		TimePeriodMonthly:    24 * 30,      // 720 hours, assuming 30 days for a simplified month
		TimePeriodQuarterly:  24 * 90,      // 2160 hours, assuming 90 days per quarter
		TimePeriodBiAnnual:   24 * 365 / 2, // 4380 hours, assuming half a year
		TimePeriodAnnual:     24 * 365,     // 8760 hours in a year
		TimePeriodBiennial:   24 * 730,     // 17520 hours, for two years
	}

	fromHours, ok := hoursPerUnit[phaseFromUnit]
	if !ok {
		return 0, errors.New("invalid fromUnit for conversion")
	}
	totalHours := float64(phaseDuration) * fromHours

	toHours, ok := hoursPerUnit[billingToUnit]
	if !ok {
		return 0, errors.New("invalid toUnit for conversion")
	}

	// Calculate the ratio of the trial or discount period to the full billing cycle.
	ratio := totalHours / toHours

	return ratio, nil
}

type CouponDiscountType int

const (
	DiscountTypeUnknown        CouponDiscountType = iota // Default value, represents an undefined discount type
	DiscountTypePercentage                               // Represents a percentage-based discount
	DiscountTypeFixedAmount                              // Represents a fixed amount discount
	DiscountTypeBuyOneGetOne                             // Represents a buy one get one free discount
	DiscountTypeFreeShipping                             // Represents a free shipping discount
	DiscountTypeGetSkuDiscount                           // Represents a SKU-specific discount
	DiscountTypePlanAccess                               // Represents access to a specific plan as a discount
	// Subscription-specific discount types
	DiscountTypeTrialPeriod            // Represents a free or discounted trial period for subscriptions
	DiscountTypeRecurringDiscount      // Represents a recurring discount over a specified number of billing cycles
	DiscountTypeFixedPriceSubscription // Represents a fixed price for a specified duration of the subscription
)

// String method to provide string representation of CouponDiscountType
func (dt CouponDiscountType) String() string {
	switch dt {
	case DiscountTypePercentage:
		return "Percentage"
	case DiscountTypeFixedAmount:
		return "FixedAmount"
	case DiscountTypeBuyOneGetOne:
		return "BuyOneGetOne"
	case DiscountTypeFreeShipping:
		return "FreeShipping"
	case DiscountTypeTrialPeriod:
		return "TrialPeriod"
	case DiscountTypeRecurringDiscount:
		return "RecurringDiscount"
	case DiscountTypeFixedPriceSubscription:
		return "FixedPriceSubscription"
	default:
		return "Unknown"
	}
}

var stringToDiscountType = map[string]CouponDiscountType{
	"Percentage":             DiscountTypePercentage,
	"FixedAmount":            DiscountTypeFixedAmount,
	"BuyOneGetOne":           DiscountTypeBuyOneGetOne,
	"FreeShipping":           DiscountTypeFreeShipping,
	"TrialPeriod":            DiscountTypeTrialPeriod,
	"RecurringDiscount":      DiscountTypeRecurringDiscount,
	"FixedPriceSubscription": DiscountTypeFixedPriceSubscription,
}

func ConvertStringToDiscountType(s string) (CouponDiscountType, error) {
	if dt, ok := stringToDiscountType[s]; ok {
		return dt, nil
	}
	return DiscountTypeUnknown, errors.New("invalid discount type")
}
