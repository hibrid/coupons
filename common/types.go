package common

import "errors"

type BillingPeriod int

const (
	BillingPeriodUnknown BillingPeriod = iota // Default value, represents an undefined billing period
	BillingPeriodDay                          // Represents a daily billing period
	BillingPeriodWeek                         // Represents a weekly billing period
	BillingPeriodMonth                        // Represents a monthly billing period
	BillingPeriodYear                         // Represents an annual billing period
)

const (
	BillingPeriodHourly     BillingPeriod = iota // Represents an hourly billing period
	BillingPeriodDaily                           // Represents a daily billing period
	BillingPeriodWeekly                          // Represents a weekly billing period
	BillingPeriodBiWeekly                        // Represents a bi-weekly billing period
	BillingPeriodThirtyDays                      // Represents a 30-day billing period
	BillingPeriodMonthly                         // Represents a monthly billing period
	BillingPeriodQuarterly                       // Represents a quarterly billing period
	BillingPeriodBiAnnual                        // Represents a bi-annual billing period
	BillingPeriodAnnual                          // Represents an annual billing period
	BillingPeriodBiennial                        // Represents a biennial billing period
	BillingPeriodNoBilling                       // Represents no billing period
)

func (bp BillingPeriod) String() string {
	switch bp {
	case BillingPeriodHourly:
		return "Hourly"
	case BillingPeriodDaily:
		return "Daily"
	case BillingPeriodWeekly:
		return "Weekly"
	case BillingPeriodBiWeekly:
		return "BiWeekly"
	case BillingPeriodThirtyDays:
		return "ThirtyDays"
	case BillingPeriodMonthly:
		return "Monthly"
	case BillingPeriodQuarterly:
		return "Quarterly"
	case BillingPeriodBiAnnual:
		return "BiAnnual"
	case BillingPeriodAnnual:
		return "Annual"
	case BillingPeriodBiennial:
		return "Biennial"
	case BillingPeriodNoBilling:
		return "NoBilling"
	default:
		return "Unknown"
	}
}

var stringToBillingPeriod = map[string]BillingPeriod{
	"Hourly":     BillingPeriodHourly,
	"Daily":      BillingPeriodDaily,
	"Weekly":     BillingPeriodWeekly,
	"BiWeekly":   BillingPeriodBiWeekly,
	"ThirtyDays": BillingPeriodThirtyDays,
	"Monthly":    BillingPeriodMonthly,
	"Quarterly":  BillingPeriodQuarterly,
	"BiAnnual":   BillingPeriodBiAnnual,
	"Annual":     BillingPeriodAnnual,
	"Biennial":   BillingPeriodBiennial,
	"NoBilling":  BillingPeriodNoBilling,
}

func ConvertStringToBillingPeriod(s string) (BillingPeriod, error) {
	if bp, ok := stringToBillingPeriod[s]; ok {
		return bp, nil
	}
	return BillingPeriodUnknown, errors.New("invalid billing period")
}

/*
const (
	BillingPeriodHourly     BillingPeriod = "HOURLY"
	BillingPeriodDaily      BillingPeriod = "DAILY"
	BillingPeriodWeekly     BillingPeriod = "WEEKLY"
	BillingPeriodBiWeekly   BillingPeriod = "BIWEEKLY"
	BillingPeriodThirtyDays BillingPeriod = "THIRTY_DAYS"
	BillingPeriodMonthly    BillingPeriod = "MONTHLY"
	BillingPeriodQuarterly  BillingPeriod = "QUARTERLY"
	BillingPeriodBiAnnual   BillingPeriod = "BIANNUAL"
	BillingPeriodAnnual     BillingPeriod = "ANNUAL"
	BillingPeriodBiennial   BillingPeriod = "BIENNIAL"
	BillingPeriodNoBilling  BillingPeriod = "NO_BILLING_PERIOD"
)
*/

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
