package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"
)

// Define a struct to represent coupon generation configuration
type CouponConfig struct {
	CouponPrefix      string
	CouponSuffix      string
	CouponVanityName  string
	CouponPattern     string
	CouponCount       int
	CouponType        CouponType
	DiscountType      CouponDiscountType
	DiscountValue     float64
	MinimumPurchase   float64
	ExpirationDate    string
	IsSingleUse       bool
	UsageLimit        int
	AvailabilityCount int
	IsActive          bool
	CampaignID        int

	TrialPeriodUnit          BillingPeriod
	TrialPeriodLength        int
	PostTrialPricing         float64
	DiscountDurationUnit     BillingPeriod
	DiscountDurationLength   int
	FixedPriceDurationUnit   BillingPeriod
	FixedPriceDurationLength int
	EligiblePlans            json.RawMessage `json:"eligible_plans"`
}

type Coupon struct {
	ID                int
	Code              string
	CouponType        CouponType
	VanityName        string
	Description       string
	DiscountType      CouponDiscountType
	DiscountValue     float64
	MinimumPurchase   float64
	ExpirationDate    string
	IsSingleUse       bool
	UsageLimit        int
	AvailabilityCount int
	IsActive          bool
	CampaignID        int
	IsValid           bool
	NotValidReason    string

	TrialDuration          BillingPeriodDuration
	PostTrialPricing       float64
	DiscountDurationCycles BillingPeriodDuration
	FixedPriceDuration     BillingPeriodDuration
	EligiblePlans          json.RawMessage `json:"eligible_plans"`
}

func (c *Coupon) String() string {
	return fmt.Sprintf("Coupon %s: %s, Discount: %s, Expiration: %s", c.Code, c.Description, c.DiscountType.String(), c.ExpirationDate)
}

func (c *Coupon) IsVanityCode() bool {
	return c.CouponType.IsVanityCode()
}

func (c *Coupon) IsReferralCode() bool {
	return c.CouponType.IsReferralCode()
}

func (c *Coupon) IsPromoCode() bool {
	return c.CouponType.IsPromoCode()
}

func (c *Coupon) IsLoyaltyCode() bool {
	return c.CouponType.IsLoyaltyCode()
}

func (c *Coupon) IsUnknownCodeType() bool {
	return c.CouponType.IsUnknown()
}

func (c *Coupon) IsKnownCodeType() bool {
	return c.CouponType.IsKnown()
}

func (c *Coupon) IsSubscriptionDiscount() bool {
	return c.DiscountType.IsSubscriptionDiscount()
}

func (c *Coupon) IsRecurringDiscount() bool {
	return c.DiscountType.IsRecurringDiscount()
}

func (c *Coupon) IsTrialPeriod() bool {
	return c.DiscountType.IsTrialPeriod()
}

func (c *Coupon) IsFixedPriceSubscription() bool {
	return c.DiscountType.IsFixedPriceSubscription()
}

func (c *Coupon) IsPercentage() bool {
	return c.DiscountType.IsPercentage()
}

func (c *Coupon) IsFixedAmount() bool {
	return c.DiscountType.IsFixedAmount()
}

func (c *Coupon) IsBuyOneGetOne() bool {
	return c.DiscountType.IsBuyOneGetOne()
}

func (c *Coupon) IsFreeShipping() bool {
	return c.DiscountType.IsFreeShipping()
}

func (c *Coupon) IsUnknownDiscountType() bool {
	return c.DiscountType.IsUnknown()
}

func (c *Coupon) IsKnownDiscountType() bool {
	return c.DiscountType.IsKnown()
}

/*
	AddEligiblePlans will marshal into json the slice of plans

coupon := Coupon{}

coupon.AddEligiblePlans([]string{"planA", "planB", "planC"})
*/
func (c *Coupon) AddEligiblePlans(planIDs []string) error {
	plansJSON, err := json.Marshal(planIDs)
	if err != nil {
		return err
	}
	c.EligiblePlans = plansJSON
	return nil
}

/*
	GetEligiblePlans returns an unmarshalled slice of plans

eligiblePlans, _ := coupon.GetEligiblePlans()

	for _, plan := range eligiblePlans {
	    println("Eligible Plan:", plan)
	}
*/
func (c *Coupon) GetEligiblePlans() ([]string, error) {
	var plans []string
	err := json.Unmarshal(c.EligiblePlans, &plans)
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (mf *Coupon) IsNewCustomer(IsNewCustomer bool) string {
	if IsNewCustomer {
		return fmt.Sprintf("this is a new customer")
	}

	return fmt.Sprintf("this is not a new customer")
}

type CouponType int

const (
	CouponTypeUnknown      CouponType = iota // Default value
	CouponTypeVanityCode                     // Vanity code is used for vanity URLs
	CouponTypeReferralCode                   // Referral code is used for referral programs
	CouponTypePromoCode                      // Promo code is used for general promotions
	CouponTypeLoyaltyCode                    // Loyalty code is used for customer loyalty programs
)

func (ct CouponType) String() string {
	switch ct {
	case CouponTypeVanityCode:
		return "VanityCodeType"
	case CouponTypeReferralCode:
		return "ReferralCodeType"
	case CouponTypePromoCode:
		return "PromoCodeType"
	case CouponTypeLoyaltyCode:
		return "LoyaltyCodeType"
	default:
		return "UnknownCodeType"
	}
}

var stringToCouponType = map[string]CouponType{
	"VanityCodeType":   CouponTypeVanityCode,
	"ReferralCodeType": CouponTypeReferralCode,
	"PromoCodeType":    CouponTypePromoCode,
	"LoyaltyCodeType":  CouponTypeLoyaltyCode,
}

func ConvertStringToCouponType(s string) (CouponType, error) {
	if ct, ok := stringToCouponType[s]; ok {
		return ct, nil
	}
	return CouponTypeUnknown, errors.New("invalid coupon type")
}

func (ct CouponType) IsVanityCode() bool {
	fmt.Println("IsVanityCode")
	return ct == CouponTypeVanityCode
}

func (ct CouponType) IsReferralCode() bool {
	return ct == CouponTypeReferralCode
}

func (ct CouponType) IsPromoCode() bool {
	return ct == CouponTypePromoCode
}

func (ct CouponType) IsLoyaltyCode() bool {
	return ct == CouponTypeLoyaltyCode
}

func (ct CouponType) IsUnknown() bool {
	return ct == CouponTypeUnknown
}

func (ct CouponType) IsKnown() bool {
	return ct != CouponTypeUnknown
}

type CouponDiscountType int

const (
	DiscountTypeUnknown      CouponDiscountType = iota // Default value
	DiscountTypePercentage                             // Percentage-based discount
	DiscountTypeFixedAmount                            // Fixed amount discount
	DiscountTypeBuyOneGetOne                           // Buy one get one free
	DiscountTypeFreeShipping                           // Free shipping

	// Subscription-specific discount types
	DiscountTypeTrialPeriod            // Free or discounted trial period for subscriptions
	DiscountTypeRecurringDiscount      // Recurring discount over a specified number of billing cycles
	DiscountTypeFixedPriceSubscription // Fixed price for a specified duration of the subscription
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

func (dt CouponDiscountType) IsSubscriptionDiscount() bool {
	return dt == DiscountTypeTrialPeriod || dt == DiscountTypeRecurringDiscount || dt == DiscountTypeFixedPriceSubscription
}

func (dt CouponDiscountType) IsRecurringDiscount() bool {
	return dt == DiscountTypeRecurringDiscount || dt == DiscountTypeFixedPriceSubscription
}

func (dt CouponDiscountType) IsTrialPeriod() bool {
	return dt == DiscountTypeTrialPeriod
}

func (dt CouponDiscountType) IsFixedPriceSubscription() bool {
	return dt == DiscountTypeFixedPriceSubscription
}

func (dt CouponDiscountType) IsPercentage() bool {
	return dt == DiscountTypePercentage
}

func (dt CouponDiscountType) IsFixedAmount() bool {
	return dt == DiscountTypeFixedAmount
}

func (dt CouponDiscountType) IsBuyOneGetOne() bool {
	return dt == DiscountTypeBuyOneGetOne
}

func (dt CouponDiscountType) IsFreeShipping() bool {
	return dt == DiscountTypeFreeShipping
}

func (dt CouponDiscountType) IsUnknown() bool {
	return dt == DiscountTypeUnknown
}

func (dt CouponDiscountType) IsKnown() bool {
	return dt != DiscountTypeUnknown
}

/*
ConvertStringToDiscountType converts a string to CouponDiscountType

		 // Example query (replace with your actual query)
	    rows, err := db.Query("SELECT discount_type FROM Coupons")
	    if err != nil {
	        panic(err)
	    }
	    defer rows.Close()

	    for rows.Next() {
	        var discountTypeStr string
	        if err := rows.Scan(&discountTypeStr); err != nil {
	            panic(err)
	        }
	        discountType, err := ConvertStringToDiscountType(discountTypeStr)
	        if err != nil {
	            fmt.Printf("Error: %s\n", err)
	        } else {
	            fmt.Printf("Discount Type: %s\n", discountType.String())
	        }
	    }

	    if err := rows.Err(); err != nil {
	        panic(err)
	    }}
		    }
*/
func ConvertStringToDiscountType(s string) (CouponDiscountType, error) {
	if dt, ok := stringToDiscountType[s]; ok {
		return dt, nil
	}
	return DiscountTypeUnknown, errors.New("invalid discount type")
}

func validateBillingPeriodUnit(period BillingPeriod) error {
	if period == "" {
		return errors.New("billing period is required")
	}
	return nil
}

func validatePeriodDuration(fieldName string, period BillingPeriodDuration) error {
	if err := validateBillingPeriodUnit(period.Period); err != nil {
		// add the field name to the error message
		return errors.New("invalid " + fieldName + " period: " + err.Error())
	}
	if period.Period != "" && period.Length <= 0 {
		return errors.New("invalid " + fieldName + " length")
	}
	return nil
}

func validateDiscountTypeConfig(discountType CouponDiscountType, config CouponConfig) error {
	if discountType.IsTrialPeriod() && config.TrialPeriodUnit == "" {
		return errors.New("trial period unit is required for TrialPeriod discount type")
	}

	if discountType.IsRecurringDiscount() && config.DiscountDurationUnit == "" {
		return errors.New("discount duration unit is required for RecurringDiscount discount type")
	}

	if discountType.IsFixedPriceSubscription() && config.FixedPriceDurationUnit == "" {
		return errors.New("fixed price duration unit is required for FixedPriceSubscription discount type")
	}
	return nil
}

func validateCouponConfig(config CouponConfig) error {
	if !config.CouponType.IsKnown() {
		return errors.New("invalid CouponType")
	}

	if !config.DiscountType.IsKnown() {
		return errors.New("invalid DiscountType")
	}

	if config.DiscountType == DiscountTypePercentage || config.DiscountType == DiscountTypeFixedAmount {
		if config.DiscountValue < 0 || config.DiscountValue > 100 {
			return errors.New("DiscountValue must be between 0 and 100")
		}
	}

	_, err := time.Parse("2006-01-02", config.ExpirationDate)
	if err != nil {
		return errors.New("invalid ExpirationDate format")
	}

	if err := validateDiscountTypeConfig(config.DiscountType, config); err != nil {
		return err
	}

	return nil
}

// Generate coupons in bulk based on configuration and return the generated coupons
func GenerateCoupons(db *sql.DB, config CouponConfig) ([]Coupon, error) {

	if config.CouponType.IsVanityCode() && config.CouponVanityName == "" {
		return nil, errors.New("vanity code name is required for Vanity code type")
	}

	validateCouponConfig(config)

	var generatedCoupons []Coupon

	stmt, err := db.Prepare(`
        INSERT INTO Coupons (
            coupon_code, 
			coupon_description, 
			coupon_type,
			coupon_vanity_name,
			discount_type, 
			discount_value, 
			minimum_purchase_amount,
            expiration_date, 
			is_single_use, 
			usage_limit, 
			availability_count,
			is_active, 
			campaign_id,
            trial_period_unit, 
			trial_period_length, 
			post_trial_pricing,
            discount_duration_unit, 
			discount_duration_length, 
			fixed_price_duration_unit,
            fixed_price_duration_length, 
			eligible_plans
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for i := 1; i <= config.CouponCount; i++ {
		couponCode := fmt.Sprintf("%s%d", config.CouponPrefix, i)
		newCoupon := Coupon{
			Code:              couponCode,
			Description:       fmt.Sprintf("%s Coupon %d", config.CouponPrefix, i),
			DiscountType:      config.DiscountType,
			DiscountValue:     config.DiscountValue,
			CouponType:        config.CouponType,
			VanityName:        config.CouponVanityName,
			MinimumPurchase:   config.MinimumPurchase,
			ExpirationDate:    config.ExpirationDate,
			IsSingleUse:       config.IsSingleUse,
			UsageLimit:        config.UsageLimit,
			AvailabilityCount: config.AvailabilityCount,
			IsActive:          config.IsActive,
			CampaignID:        config.CampaignID,
			IsValid:           true, // Set to true initially, can be changed based on rules later
			NotValidReason:    "",

			TrialDuration: BillingPeriodDuration{
				Period: config.TrialPeriodUnit,
				Length: config.TrialPeriodLength,
			},
			PostTrialPricing: config.PostTrialPricing,
			DiscountDurationCycles: BillingPeriodDuration{
				Period: config.DiscountDurationUnit,
				Length: config.DiscountDurationLength,
			},
			FixedPriceDuration: BillingPeriodDuration{
				Period: config.FixedPriceDurationUnit,
				Length: config.FixedPriceDurationLength,
			},
			EligiblePlans: json.RawMessage{},
		}

		eligiblePlansJSON, err := json.Marshal(newCoupon.EligiblePlans)
		if err != nil {
			return nil, err
		}
		_, err = stmt.Exec(
			newCoupon.Code, newCoupon.Description, newCoupon.CouponType.String(), newCoupon.DiscountType.String(), newCoupon.DiscountValue, newCoupon.MinimumPurchase,
			newCoupon.ExpirationDate, newCoupon.IsSingleUse, newCoupon.UsageLimit, newCoupon.AvailabilityCount, newCoupon.IsActive, newCoupon.CampaignID,
			newCoupon.TrialDuration.Period, newCoupon.TrialDuration.Length, newCoupon.PostTrialPricing,
			newCoupon.DiscountDurationCycles.Period, newCoupon.DiscountDurationCycles.Length, newCoupon.FixedPriceDuration.Period,
			newCoupon.FixedPriceDuration.Length, eligiblePlansJSON,
		)
		if err != nil {
			return nil, err
		}

		generatedCoupons = append(generatedCoupons, newCoupon)
	}

	fmt.Printf("%d Coupons generated successfully\n", config.CouponCount)
	return generatedCoupons, nil
}

/*
// Example usage
	daysUntilExpiration := 10
	coupons, err := findCouponsExpiringInDays(db, daysUntilExpiration)
	if err != nil {
		fmt.Println("Error finding coupons:", err)
		return
	}

	for _, coupon := range coupons {
		fmt.Printf("Coupon expiring in %d days: %+v\n", daysUntilExpiration, coupon)
	}
*/

func FindCouponsExpiringInDays(db *sql.DB, days int) ([]Coupon, error) {
	// Calculate the date range
	today := time.Now()
	endDate := today.AddDate(0, 0, days)

	// Prepare the SQL query
	query := `SELECT id, code, description, discount_type, discount_value, minimum_purchase,
	expiration_date, is_single_use, usage_limit, is_active, is_valid, not_valid_reason,
	campaign_id, trial_period_unit, trial_period_length, post_trial_pricing,
	discount_duration_unit, discount_duration_length, fixed_price_duration_unit,
	fixed_price_duration_length, eligible_plans FROM Coupons WHERE expiration_date BETWEEN ? AND ?`

	// Execute the query
	rows, err := db.Query(query, today.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the coupons
	var coupons []Coupon
	for rows.Next() {
		var coupon Coupon
		err := rows.Scan(&coupon.ID, &coupon.Code, &coupon.Description, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MinimumPurchase, &coupon.ExpirationDate, &coupon.IsSingleUse, &coupon.UsageLimit, &coupon.IsActive, &coupon.CampaignID)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return coupons, nil
}

/*
// Example usage of findExpiredCoupons
	startDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	expiredCoupons, err := findExpiredCoupons(db, startDate, endDate)
	if err != nil {
		fmt.Println("Error finding expired coupons:", err)
		return
	}

	// Do something with expiredCoupons
	for _, coupon := range expiredCoupons {
		fmt.Printf("Expired Coupon: %+v\n", coupon)
	}
*/

func FindExpiredCoupons(db *sql.DB, start, end time.Time) ([]Coupon, error) {
	// Prepare the SQL query
	query := `SELECT id, code, description, discount_type, discount_value, minimum_purchase,
	expiration_date, is_single_use, usage_limit, is_active, is_valid, not_valid_reason,
	campaign_id, trial_period_unit, trial_period_length, post_trial_pricing,
	discount_duration_unit, discount_duration_length, fixed_price_duration_unit,
	fixed_price_duration_length, eligible_plans 
	FROM Coupons WHERE expiration_date BETWEEN ? AND ?`

	// Execute the query
	rows, err := db.Query(query, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the coupons
	var coupons []Coupon
	for rows.Next() {
		var coupon Coupon
		err := rows.Scan(&coupon.ID, &coupon.Code, &coupon.Description, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MinimumPurchase, &coupon.ExpirationDate, &coupon.IsSingleUse, &coupon.UsageLimit, &coupon.IsActive, &coupon.CampaignID)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return coupons, nil
}

// Record coupon usage
func RecordCouponUsage(db *sql.DB, couponID, userID, orderID int) {
	_, err := db.Exec("INSERT INTO CouponUsage (coupon_id, user_id, order_id, usage_date, is_used) "+
		"VALUES (?, ?, ?, NOW(), true)",
		couponID, userID, orderID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Coupon usage recorded successfully for Coupon ID: %d\n", couponID)
}

// Retrieve coupons associated with a campaign by Campaign ID
func GetCouponsByCampaignID(db *sql.DB, campaignID int) ([]Coupon, error) {
	rows, err := db.Query(`
	SELECT id, code, description, discount_type, discount_value, minimum_purchase,
	expiration_date, is_single_use, usage_limit, is_active, is_valid, not_valid_reason,
	campaign_id, trial_period_unit, trial_period_length, post_trial_pricing,
	discount_duration_unit, discount_duration_length, fixed_price_duration_unit,
	fixed_price_duration_length, eligible_plans
	FROM Coupons WHERE campaign_id = ?
	`, campaignID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []Coupon
	for rows.Next() {
		var coupon Coupon
		var trialPeriodUnit, discountDurationUnit, fixedPriceDurationUnit string
		var eligiblePlansJSON json.RawMessage

		err := rows.Scan(
			&coupon.ID, &coupon.Code, &coupon.Description, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MinimumPurchase,
			&coupon.ExpirationDate, &coupon.IsSingleUse, &coupon.UsageLimit, &coupon.IsActive, &coupon.IsValid, &coupon.NotValidReason,
			&coupon.CampaignID, &trialPeriodUnit, &coupon.TrialDuration.Length, &coupon.PostTrialPricing,
			&discountDurationUnit, &coupon.DiscountDurationCycles.Length, &fixedPriceDurationUnit,
			&coupon.FixedPriceDuration.Length, &eligiblePlansJSON,
		)
		if err != nil {
			return nil, err
		}
		// Convert string representation back to BillingPeriod
		coupon.TrialDuration.Period = BillingPeriod(trialPeriodUnit)
		coupon.DiscountDurationCycles.Period = BillingPeriod(discountDurationUnit)
		coupon.FixedPriceDuration.Period = BillingPeriod(fixedPriceDurationUnit)

		// Deserialize eligible_plans
		err = json.Unmarshal(eligiblePlansJSON, &coupon.EligiblePlans)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return coupons, nil
}
