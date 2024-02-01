package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// GenerateRandomCoupon generates a random coupon for testing
func GenerateRandomCoupon() Coupon {
	// Random coupon code (e.g., "SUMMER25")
	couponCode := generateRandomCouponCode()

	// Random coupon description
	couponDescription := generateRandomCouponDescription()

	// Random discount type (either "percentage" or "fixed")
	discountType := generateRandomDiscountType()

	// Random discount value (between 0% and 100% or 0 and 100)
	discountValue := generateRandomDiscountValue(discountType)

	// Random minimum purchase amount (between 0 and 1000)
	minimumPurchase := generateRandomMinimumPurchase()

	// Random expiration date (within the next 365 days)
	expirationDate := generateRandomExpirationDate()

	// Random single-use flag (true or false)
	isSingleUse := generateRandomIsSingleUse()

	// Random usage limit (between 1 and 100)
	usageLimit := generateRandomUsageLimit()

	// Random isActive flag (true or false)
	isActive := generateRandomIsActive()

	// Random campaign ID (between 1 and 10)
	campaignID := generateRandomCampaignID()

	trialDuration := generateRandomBillingPeriodDuration()
	postTrialPricing := generateRandomPostTrialPricing()
	discountDurationCycles := generateRandomBillingPeriodDuration()
	fixedPriceDuration := generateRandomBillingPeriodDuration()
	eligiblePlans := generateRandomEligiblePlans()

	return Coupon{
		Code:                   couponCode,
		Description:            couponDescription,
		DiscountType:           discountType,
		DiscountValue:          discountValue,
		MinimumPurchase:        minimumPurchase,
		ExpirationDate:         expirationDate,
		IsSingleUse:            isSingleUse,
		UsageLimit:             usageLimit,
		IsActive:               isActive,
		CampaignID:             campaignID,
		TrialDuration:          trialDuration,
		PostTrialPricing:       postTrialPricing,
		DiscountDurationCycles: discountDurationCycles,
		FixedPriceDuration:     fixedPriceDuration,
		EligiblePlans:          eligiblePlans,
	}
}

// Helper functions to generate random values

func generateRandomBillingPeriodDuration() BillingPeriodDuration {
	periods := []BillingPeriod{BillingPeriodHourly, BillingPeriodDaily, BillingPeriodWeekly /* other periods */}
	period := periods[rand.Intn(len(periods))]
	length := rand.Intn(12) + 1 // Random length between 1 and 12
	return BillingPeriodDuration{
		Period: period,
		Length: length,
	}
}

func generateRandomPostTrialPricing() float64 {
	return float64(rand.Intn(200)) // Random pricing between 0 and 200
}

func generateRandomEligiblePlans() json.RawMessage {
	numPlans := rand.Intn(5) + 1 // Random number of plans between 1 and 5
	plans := make([]int, numPlans)
	for i := 0; i < numPlans; i++ {
		plans[i] = rand.Intn(10) + 1 // Random plan IDs between 1 and 10
	}

	// Convert the slice of integers to json.RawMessage
	plansJSON, err := json.Marshal(plans)
	if err != nil {
		// Handle the error appropriately in your application
		panic(err) // For example purposes only
	}
	return json.RawMessage(plansJSON)
}

func generateRandomCouponCode() string {
	// Generate a random coupon code (e.g., "SUMMER25")
	return fmt.Sprintf("COUPON%d", rand.Intn(1000))
}

func generateRandomCouponDescription() string {
	// Generate a random coupon description
	descriptions := []string{"Summer Discount", "Holiday Special", "Flash Sale", "New User Offer"}
	return descriptions[rand.Intn(len(descriptions))]
}

func generateRandomDiscountType() CouponDiscountType {
	// Generate a random discount type from the available types
	types := []CouponDiscountType{
		DiscountTypePercentage,
		DiscountTypeFixedAmount,
		DiscountTypeBuyOneGetOne,
		DiscountTypeFreeShipping,
		DiscountTypeTrialPeriod,
		DiscountTypeRecurringDiscount,
		DiscountTypeFixedPriceSubscription,
	}
	return types[rand.Intn(len(types))]
}

func generateRandomDiscountValue(discountType CouponDiscountType) float64 {
	// Generate a random discount value based on the discount type
	switch discountType {
	case DiscountTypePercentage, DiscountTypeRecurringDiscount:
		// Percentage-based discounts (0% to 100%)
		return float64(rand.Intn(101))
	case DiscountTypeFixedAmount, DiscountTypeFixedPriceSubscription:
		// Fixed amount discounts (e.g., $0 to $100)
		return float64(rand.Intn(101))
	default:
		// For other types, you might return a standard value or based on some other logic
		return 0
	}
}

func generateRandomMinimumPurchase() float64 {
	// Generate a random minimum purchase amount (between 0 and 1000)
	return float64(rand.Intn(1001))
}

func generateRandomExpirationDate() string {
	// Generate a random expiration date within the next 365 days
	now := time.Now()
	randomExpiration := now.Add(time.Duration(rand.Intn(365)) * 24 * time.Hour)
	return randomExpiration.Format("2006-01-02")
}

func generateRandomIsSingleUse() bool {
	// Generate a random single-use flag (true or false)
	return rand.Intn(2) == 1
}

func generateRandomUsageLimit() int {
	// Generate a random usage limit (between 1 and 100)
	return rand.Intn(100) + 1
}

func generateRandomIsActive() bool {
	// Generate a random isActive flag (true or false)
	return rand.Intn(2) == 1
}

func generateRandomCampaignID() int {
	// Generate a random campaign ID (between 1 and 10)
	return rand.Intn(10) + 1
}

func GenerateRandomChangeContext() ChangeContext {
	// Generate two random subscriptions as the "from" and "to" subscriptions
	fromSubscription := GenerateRandomSubscription()
	toSubscription := GenerateRandomSubscription()

	return ChangeContext{
		FromSubscription: fromSubscription,
		ToSubscription:   toSubscription,
	}
}

func GenerateRandomOptionsContext() OptionsContext {
	return OptionsContext{
		Option1: rand.Intn(2) == 1,                       // Random boolean value (true or false)
		Option2: fmt.Sprintf("Option%d", rand.Intn(3)+1), // Random string value (Option1, Option2, Option3)
		Option3: rand.Intn(100),                          // Random integer value (0 to 99)
	}
}

func GenerateRandomCustomerContext() CustomerContext {
	return CustomerContext{
		UserID:              rand.Intn(1000),                   // Random user ID
		IsNewCustomer:       rand.Intn(2) == 1,                 // Random new customer flag
		IsSubscriber:        rand.Intn(2) == 1,                 // Random subscriber flag
		PreviousRedemptions: generateRandomRedemptionHistory(), // Generate random redemption history
		// Add other relevant customer data initialization here
	}
}

func generateRandomRedemptionHistory() []RedemptionHistoryEntry {
	numEntries := rand.Intn(11) // Random number of redemption history entries (0 to 10)
	var history []RedemptionHistoryEntry
	for i := 0; i < numEntries; i++ {
		entry := RedemptionHistoryEntry{
			CouponCode:     fmt.Sprintf("COUPON%d", rand.Intn(1000)), // Random coupon code
			RedemptionDate: generateRandomRedemptionDate(),           // Random redemption date
		}
		history = append(history, entry)
	}
	return history
}

func generateRandomRedemptionDate() string {
	// Generate a random redemption date within the last year
	now := time.Now()
	maxDate := now.AddDate(0, 0, -365)
	randomDate := maxDate.Add(time.Duration(rand.Intn(365)) * 24 * time.Hour)
	return randomDate.Format("2006-01-02")
}

func GenerateRandomSubscription() Subscription {
	startDate := generateRandomDate()
	endDate := startDate.AddDate(1, 0, 0) // Set the end date to be 1 year from the start date

	return Subscription{
		ID:               rand.Intn(1000),             // Random subscription ID
		UserID:           rand.Intn(1000),             // Random user ID
		PlanID:           rand.Intn(100),              // Random plan ID
		Term:             getRandomTerm(),             // Random term (e.g., monthly, annual)
		MonthlyPrice:     rand.Float64() * 100,        // Random monthly price
		StartDate:        startDate,                   // Random start date
		EndDate:          endDate,                     // End date is 1 year from start date
		IsActive:         rand.Intn(2) == 1,           // Random active flag
		SubscriptionType: getRandomSubscriptionType(), // Random subscription type (e.g., basic, premium)
		// Add other relevant subscription details initialization here
	}
}

func generateRandomDate() time.Time {
	// Generate a random date within the last year
	now := time.Now()
	maxDate := now.AddDate(0, 0, -365)
	randomDate := maxDate.Add(time.Duration(rand.Intn(365)) * 24 * time.Hour)
	return randomDate
}

func getRandomTerm() string {
	terms := []string{"monthly", "annual", "quarterly"} // Add more terms as needed
	return terms[rand.Intn(len(terms))]
}

func getRandomSubscriptionType() string {
	types := []string{"basic", "premium", "pro"} // Add more subscription types as needed
	return types[rand.Intn(len(types))]
}
