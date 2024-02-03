package main

import (
	"testing"
)

func TestCouponTypeString(t *testing.T) {
	tests := []struct {
		couponType CouponType
		expected   string
	}{
		{CouponTypeVanityCode, "VanityCodeType"},
		{CouponTypeReferralCode, "ReferralCodeType"},
		{CouponTypePromoCode, "PromoCodeType"},
		{CouponTypeLoyaltyCode, "LoyaltyCodeType"},
		{CouponTypeUnknown, "UnknownCodeType"},
	}

	for _, test := range tests {
		result := test.couponType.String()
		if result != test.expected {
			t.Errorf("CouponType.String() returned %s, expected %s", result, test.expected)
		}
	}
}

func TestConvertStringToCouponType(t *testing.T) {
	tests := []struct {
		input         string
		expected      CouponType
		errorExpected bool
	}{
		{"VanityCodeType", CouponTypeVanityCode, false},
		{"ReferralCodeType", CouponTypeReferralCode, false},
		{"PromoCodeType", CouponTypePromoCode, false},
		{"LoyaltyCodeType", CouponTypeLoyaltyCode, false},
		{"UnknownCodeType", CouponTypeUnknown, true},
	}

	for _, test := range tests {
		result, err := ConvertStringToCouponType(test.input)
		if err != nil && !test.errorExpected {
			t.Errorf("ConvertStringToCouponType(%s) returned an error: %v", test.input, err)
		}
		if err == nil && test.errorExpected {
			t.Errorf("ConvertStringToCouponType(%s) did not return an error", test.input)
		}
		if result != test.expected {
			t.Errorf("ConvertStringToCouponType(%s) returned %d, expected %d", test.input, result, test.expected)
		}
	}

	// Test an invalid input
	_, err := ConvertStringToCouponType("InvalidType")
	if err == nil {
		t.Errorf("ConvertStringToCouponType('InvalidType') did not return an error")
	}
}

func TestCouponTypeMethods(t *testing.T) {
	tests := []struct {
		couponType     CouponType
		expected       bool
		functionToTest bool
	}{
		{CouponTypeVanityCode, true, CouponTypeVanityCode.IsVanityCode()},
		{CouponTypeReferralCode, true, CouponTypeReferralCode.IsReferralCode()},
		{CouponTypePromoCode, true, CouponTypePromoCode.IsPromoCode()},
		{CouponTypeLoyaltyCode, true, CouponTypeLoyaltyCode.IsLoyaltyCode()},
		{CouponTypeUnknown, true, CouponTypeUnknown.IsUnknown()},
	}

	for _, test := range tests {
		if result := test.functionToTest; result != test.expected {
			t.Errorf("CouponTypeCheck returned %v for %v, expected %v", result, test.couponType, test.expected)
		}

	}
}

func TestCouponDiscountTypeString(t *testing.T) {
	tests := []struct {
		discountType CouponDiscountType
		expected     string
	}{
		{DiscountTypePercentage, "Percentage"},
		{DiscountTypeFixedAmount, "FixedAmount"},
		{DiscountTypeBuyOneGetOne, "BuyOneGetOne"},
		{DiscountTypeFreeShipping, "FreeShipping"},
		{DiscountTypeTrialPeriod, "TrialPeriod"},
		{DiscountTypeRecurringDiscount, "RecurringDiscount"},
		{DiscountTypeFixedPriceSubscription, "FixedPriceSubscription"},
		{DiscountTypeUnknown, "Unknown"},
	}

	for _, test := range tests {
		result := test.discountType.String()
		if result != test.expected {
			t.Errorf("CouponDiscountType.String() returned %s, expected %s", result, test.expected)
		}
	}
}

func TestCouponDiscountTypeMethods(t *testing.T) {
	tests := []struct {
		discountType   CouponDiscountType
		expected       bool
		functionToTest bool
	}{
		{DiscountTypeTrialPeriod, true, DiscountTypeTrialPeriod.IsTrialPeriod()},
		{DiscountTypeRecurringDiscount, true, DiscountTypeRecurringDiscount.IsRecurringDiscount()},
		{DiscountTypeFixedPriceSubscription, true, DiscountTypeFixedPriceSubscription.IsFixedPriceSubscription()},
		{DiscountTypePercentage, false, DiscountTypePercentage.IsTrialPeriod()},
		{DiscountTypeFixedAmount, false, DiscountTypeFixedAmount.IsTrialPeriod()},
		{DiscountTypeBuyOneGetOne, false, DiscountTypeBuyOneGetOne.IsTrialPeriod()},
		{DiscountTypeFreeShipping, false, DiscountTypeFreeShipping.IsTrialPeriod()},
		{DiscountTypeUnknown, false, DiscountTypeUnknown.IsTrialPeriod()},
		{DiscountTypeTrialPeriod, false, DiscountTypeTrialPeriod.IsRecurringDiscount()},
		{DiscountTypeRecurringDiscount, false, DiscountTypeRecurringDiscount.IsFixedPriceSubscription()},
		{DiscountTypeFixedAmount, false, DiscountTypeFixedAmount.IsFixedPriceSubscription()},
		{DiscountTypeFixedPriceSubscription, true, DiscountTypeFixedPriceSubscription.IsRecurringDiscount()},
		{DiscountTypePercentage, true, DiscountTypePercentage.IsPercentage()},
		{DiscountTypeFixedAmount, true, DiscountTypeFixedAmount.IsFixedAmount()},
		{DiscountTypeBuyOneGetOne, true, DiscountTypeBuyOneGetOne.IsBuyOneGetOne()},
		{DiscountTypeFreeShipping, true, DiscountTypeFreeShipping.IsFreeShipping()},
		{DiscountTypeUnknown, true, DiscountTypeUnknown.IsUnknown()},
	}

	for _, test := range tests {
		if result := test.functionToTest; result != test.expected {
			t.Errorf("CouponTypeCheck returned %v for %v, expected %v", result, test.discountType, test.expected)
		}
	}
}
