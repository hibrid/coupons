package common

import "testing"

func TestCouponDiscountTypeString(t *testing.T) {
	cases := []struct {
		name     string
		dt       CouponDiscountType
		expected string
	}{
		{"Percentage", DiscountTypePercentage, "Percentage"},
		{"FixedAmount", DiscountTypeFixedAmount, "FixedAmount"},
		{"BuyOneGetOne", DiscountTypeBuyOneGetOne, "BuyOneGetOne"},
		{"FreeShipping", DiscountTypeFreeShipping, "FreeShipping"},
		{"TrialPeriod", DiscountTypeTrialPeriod, "TrialPeriod"},
		{"RecurringDiscount", DiscountTypeRecurringDiscount, "RecurringDiscount"},
		{"FixedPriceSubscription", DiscountTypeFixedPriceSubscription, "FixedPriceSubscription"},
		{"Unknown", DiscountTypeUnknown, "Unknown"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.dt.String(); got != tc.expected {
				t.Errorf("String() = %v, want %v for %v", got, tc.expected, tc.name)
			}
		})
	}
}

func TestConvertStringToDiscountType(t *testing.T) {
	// Assuming stringToDiscountType map is defined like this:
	// var stringToDiscountType = map[string]CouponDiscountType{
	// 	"Percentage":              DiscountTypePercentage,
	// 	"FixedAmount":             DiscountTypeFixedAmount,
	// 	// Add all mappings...
	// }

	cases := []struct {
		name           string
		input          string
		expectedType   CouponDiscountType
		expectingError bool
	}{
		{"ValidPercentage", "Percentage", DiscountTypePercentage, false},
		{"ValidFixedAmount", "FixedAmount", DiscountTypeFixedAmount, false},
		// Add tests for all valid cases...
		{"InvalidType", "InvalidType", DiscountTypeUnknown, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotType, err := ConvertStringToDiscountType(tc.input)
			if tc.expectingError {
				if err == nil {
					t.Errorf("ConvertStringToDiscountType(%v) expected an error, got nil", tc.input)
				}
			} else {
				if err != nil || gotType != tc.expectedType {
					t.Errorf("ConvertStringToDiscountType(%v) = %v, %v; want %v, <nil>", tc.input, gotType, err, tc.expectedType)
				}
			}
		})
	}
}

func TestTimeString(t *testing.T) {
	cases := []struct {
		name     string
		bp       TimeUnit
		expected string
	}{
		{"Hourly", TimePeriodHourly, "Hourly"},
		{"Daily", TimePeriodDaily, "Daily"},
		{"Weekly", TimePeriodWeekly, "Weekly"},
		{"BiWeekly", TimePeriodBiWeekly, "BiWeekly"},
		{"ThirtyDays", TimePeriodThirtyDays, "ThirtyDays"},
		{"Monthly", TimePeriodMonthly, "Monthly"},
		{"Quarterly", TimePeriodQuarterly, "Quarterly"},
		{"BiAnnual", TimePeriodBiAnnual, "BiAnnual"},
		{"Annual", TimePeriodAnnual, "Annual"},
		{"Biennial", TimePeriodBiennial, "Biennial"},
		{"NoBilling", TimePeriodNoBilling, "NoBilling"},
		{"Unknown", TimePeriodUnknown, "Unknown"}, // Explicitly testing the default case
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.bp.String(); got != tc.expected {
				t.Errorf("String() = %v, want %v for %v", got, tc.expected, tc.name)
			}
		})
	}
}

func TestHourValue(t *testing.T) {
	cases := []struct {
		name     string
		bp       TimeUnit
		expected int64
	}{
		{"Hourly", TimePeriodHourly, 1},
		{"Daily", TimePeriodDaily, 24},
		{"Weekly", TimePeriodWeekly, 24 * 7},
		{"BiWeekly", TimePeriodBiWeekly, 24 * 14},
		{"ThirtyDays", TimePeriodThirtyDays, 24 * 30},
		{"Monthly", TimePeriodMonthly, 24 * 30}, // Assuming 30 days in a month for simplicity
		{"Quarterly", TimePeriodQuarterly, 24 * 90},
		{"BiAnnual", TimePeriodBiAnnual, 24 * 365 / 2},
		{"Annual", TimePeriodAnnual, 24 * 365},
		{"Biennial", TimePeriodBiennial, 24 * 730},
		{"NoBilling or Unknown", TimePeriodNoBilling, 0}, // Assuming you want to test the default case here too
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.bp.HourValue(); got != tc.expected {
				t.Errorf("HourValue() = %v, want %v for %v", got, tc.expected, tc.name)
			}
		})
	}
}

func TestIsValidWithValidTimePeriods(t *testing.T) {
	for i := TimePeriodHourly; i <= TimePeriodNoBilling; i++ {
		if !i.IsValid() {
			t.Errorf("IsValid() for %v should be true, got false", i)
		}
	}
}

func TestIsValidWithInvalidTimePeriods(t *testing.T) {
	invalidPeriods := []TimeUnit{TimePeriodUnknown - 1, TimePeriodNoBilling + 1}
	for _, bp := range invalidPeriods {
		if bp.IsValid() {
			t.Errorf("IsValid() for %v should be false, got true", bp)
		}
	}
}

func TestHourValueForHourly(t *testing.T) {
	if got := TimePeriodHourly.HourValue(); got != 1 {
		t.Errorf("HourValue() = %v, want %v", got, 1)
	}
}

// Additional tests to be written similarly...

func TestConvertStringToTimePeriodWithValidStrings(t *testing.T) {
	testCases := map[string]TimeUnit{
		"Hourly": TimePeriodHourly,
		"Daily":  TimePeriodDaily,
		// Add the rest...
	}
	for str, expected := range testCases {
		if got, err := ConvertStringToTimePeriod(str); err != nil || got != expected {
			t.Errorf("ConvertStringToTimePeriod(%v) = %v, %v; want %v, <nil>", str, got, err, expected)
		}
	}
}

func TestConvertStringToTimePeriodWithInvalidString(t *testing.T) {
	_, err := ConvertStringToTimePeriod("NotAValidPeriod")
	if err == nil {
		t.Errorf("ConvertStringToTimePeriod(\"NotAValidPeriod\") expected an error, got nil")
	}
}
