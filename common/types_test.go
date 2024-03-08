package common

import (
	"testing"

	"github.com/shopspring/decimal"
)

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
