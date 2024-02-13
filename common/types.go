package common

import "errors"

type TimeUnit int

const (
	TimePeriodUnknown TimeUnit = iota // Default value, represents an undefined billing period
	TimePeriodDay                     // Represents a daily billing period
	TimePeriodWeek                    // Represents a weekly billing period
	TimePeriodMonth                   // Represents a monthly billing period
	TimePeriodYear                    // Represents an annual billing period
)

const (
	TimePeriodHourly     TimeUnit = iota // Represents an hourly billing period
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
