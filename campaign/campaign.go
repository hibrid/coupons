package campaign

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hibrid/coupons/common"
	"github.com/hibrid/coupons/coupon"
	"github.com/hibrid/coupons/strategy"
)

// Simulated database transaction (for illustration)
type DBTransaction struct {
	// Imagine this struct holds database transaction state
}

func (dt *DBTransaction) Begin() error {
	// Begin transaction logic here
	return nil
}

func (dt *DBTransaction) Commit() error {
	// Commit transaction logic here
	return nil
}

func (dt *DBTransaction) Rollback() error {
	// Rollback transaction logic here
	return nil
}

type CampaignType int

const (
	CampaignTypeUnknown      CampaignType = iota // Represents an undefined or default state
	CampaignTypeVanityCode                       // For custom, easy-to-remember promotional codes
	CampaignTypeReferralCode                     // For tracking and incentivizing referral marketing programs
	CampaignTypePromoCode                        // For general promotions offering discounts or special deals
	CampaignTypeLoyaltyCode                      // For rewarding repeat customers as part of loyalty programs
)

// Assuming CampaignTypeLoyaltyCode is the last and highest enum value
const maxCampaignType = CampaignTypeLoyaltyCode

// String converts the CampaignType enum to a string.
func (ct CampaignType) String() string {
	switch ct {
	case CampaignTypeUnknown:
		return "Unknown"
	case CampaignTypeVanityCode:
		return "Vanity Code"
	case CampaignTypeReferralCode:
		return "Referral Code"
	case CampaignTypePromoCode:
		return "Promo Code"
	case CampaignTypeLoyaltyCode:
		return "Loyalty Code"
	default:
		return "Invalid"
	}
}

type CampaignConfig struct {
	ID               uuid.UUID    `json:"id"`
	IsCampaignActive bool         `json:"isActive"`
	StartDate        time.Time    `json:"startDate"`
	EndDate          time.Time    `json:"endDate"`
	CampaignType     CampaignType `json:"campaignType"` // Could be an ENUM in the DB

	//PregeneratedCoupons  []Coupon        `json:"pregeneratedCoupons,omitempty"`
	AllowOnDemandCoupons bool            `json:"allowOnDemandCoupons"`
	PregenerateCoupons   bool            `json:"pregenerateCoupons"`
	EligiblePlans        json.RawMessage `json:"eligiblePlans,omitempty"`

	// coupon enforcement vars
	ExpireAfter       time.Duration
	IsSingleUse       bool // Whether a coupon can be used only once
	UsageLimit        int  // How many times a coupon can be used
	AvailabilityCount int  // How many coupons are available
	RedeemedCount     int  // How many coupons have been redeemed
	LimitPerUser      int  // How many times a user can use a coupon
	SHA256Index       int  // Index to use in the check character

	StrategyConfig json.RawMessage `json:"strategyConfig"`
	Strategy       strategy.DiscountStrategy

	CouponConfig coupon.CouponConfig // Embedded coupon configuration

	Cart *common.Cart // Embedded cart

}

func (c *CampaignConfig) GetCart() *common.Cart {
	// Return cart items for the campaign
	return c.Cart
}

// CreateCampaign validates and creates a new campaign, optionally starting a new transaction
func (c *CampaignConfig) CreateCampaign(dbTrans *DBTransaction) error {
	// Validate campaign for creation
	if err := c.ValidateCreation(); err != nil {
		return err
	}

	// Start a new transaction if none is passed
	if dbTrans == nil {
		dbTrans = &DBTransaction{}
		if err := dbTrans.Begin(); err != nil {
			return err
		}
		defer dbTrans.Rollback() // Rollback if commit is not reached
	}

	// Insert campaign into the database here using dbTrans
	// ...

	// Commit the transaction if it was started within this method
	if err := dbTrans.Commit(); err != nil {
		return err
	}

	return nil
}

func (c *CampaignConfig) GetCampaignDetails() (startDate, endDate string, isActive bool) {
	// Return necessary details from the CampaignConfig
	return c.StartDate.String(), c.EndDate.String(), c.IsCampaignActive
}

func (c *CampaignConfig) GetStartDate() string {
	return c.StartDate.String()
}

func (c *CampaignConfig) GetEndDate() string {
	return c.EndDate.String()
}

// Method to set the discount strategy
func (c *CampaignConfig) SetDiscountStrategy(s strategy.DiscountStrategy) {
	c.Strategy = s
}

// Method to apply the discount strategy
func (c *CampaignConfig) ApplyDiscountStrategy() *common.DiscountResult {
	if c.Strategy == nil || !c.ValidateStrategy() {
		return nil // No strategy set or strategy validation failed
	}
	c.Strategy.ApplyDiscount(c)
	return c.Strategy.GetDiscountResult()
}

// Method to validate the strategy within the context of the campaign
func (c *CampaignConfig) ValidateStrategy() bool {
	if c.Strategy == nil {
		return false // No strategy set, cannot validate
	}
	return c.Strategy.ValidateInputs(c)
}

func (c *CampaignConfig) IsPlanEligible(requestedPlan string) error {
	var eligiblePlans []string
	if err := json.Unmarshal(c.EligiblePlans, &eligiblePlans); err != nil {
		return &ValidationError{"invalid eligible plans data"}
	}

	for _, plan := range eligiblePlans {
		if plan == requestedPlan {
			return nil
		}
	}

	return &ValidationError{"requested plan is not eligible for this campaign"}
}

type StrategyConfigStruct struct {
	// Example fields
	MinSpendAmount float64 `json:"minSpendAmount"`
}

func (c *CampaignConfig) IsStrategyValid() error {
	var strategy StrategyConfigStruct
	if err := json.Unmarshal(c.StrategyConfig, &strategy); err != nil {
		return &ValidationError{"invalid strategy config data"}
	}

	// Example validation
	if strategy.MinSpendAmount <= 0 {
		return &ValidationError{"strategy config min spend amount must be greater than 0"}
	}

	// Implement additional validations based on strategy config rules

	return nil
}

func (c *CampaignConfig) PreviewRedemption(requestedPlan string) error {
	// Perform general redemption validations
	if err := c.ValidateForRedemption(); err != nil {
		return err
	}

	// Validate the requested plan is eligible
	if err := c.IsPlanEligible(requestedPlan); err != nil {
		return err
	}

	// Optionally, check strategy-specific validations
	if err := c.IsStrategyValid(); err != nil {
		return err
	}

	// At this point, all validations for a redemption preview have passed
	// Note: Actual redemption logic affecting state (e.g., decrementing AvailabilityCount) would be skipped

	return nil
}

// ValidateCreation includes validations specific to campaign creation
func (c *CampaignConfig) ValidateCreation() error {
	// Perform general validations
	if err := c.Validate(); err != nil {
		return err
	}

	// Perform validations specific to campaign creation
	validators := []func() error{
		c.AreDatesValid,                        // Ensure the start date is before the end date
		c.IsEndDateValid,                       // Ensure end date is not in the past
		c.IsStartDateValid,                     // Ensure start date is not in the past for new campaigns
		c.IsCampaignIDValid,                    // Ensure the campaign ID is valid (though usually generated and not a common failure point)
		c.IsUsageLimitValid,                    // Validate the usage limit is not negative
		c.IsLimitPerUserValid,                  // Validate the limit per user is not negative
		c.IsUsageLimitValidForSingleUse,        // Validate the usage limit is appropriate for single/multiple use
		c.IsUsageLimitSufficientForDemand,      // Validate the usage limit is sufficient for on-demand coupons
		c.IsAvailabilityValidForPregeneration,  // Validate the availability count is appropriate for pregenerated coupons
		c.IsUsageLimitGreaterThanRedeemedCount, // Validate the usage limit is greater than the redeemed count
	}

	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateForRedemption includes validations specific to coupon redemption
func (c *CampaignConfig) ValidateForRedemption() error {
	// Perform general validations
	if err := c.Validate(); err != nil {
		return err
	}

	// Perform validations specific to coupon redemption
	if err := c.IsActive(); err != nil {
		return err
	}

	if c.IsExpired() {
		return &ValidationError{"campaign is expired"}
	}

	if !c.IsStarted() {
		return &ValidationError{"campaign has not started"}
	}

	if !c.HasRemainingCoupons() {
		return &ValidationError{"no remaining coupons"}
	}

	// check the eligible plans for the campaign against the requested plan
	c.IsPlanEligibleForRedemption("")

	// check if the coupon has expired

	// check if the coupon has been redeemed. We need to check the redeemed count against the usage limit

	// check if the coupon is single use and has already been redeemed

	// check if the coupon has been redeemed by the user. We need to check the user's redeemed count against the limit per user

	// Depending on requirements, add other redemption-specific validations here
	// For example, checking the redeemed count vs. usage limit, etc.

	return nil
}

// check the eligible plans for the campaign against the requested plan
func (c *CampaignConfig) IsPlanEligibleForRedemption(requestedPlan string) error {
	return c.IsPlanEligible(requestedPlan)
}

// EvaluateCampaignForRedemption checks if a coupon can be redeemed against this campaign, updating records as necessary
func (c *CampaignConfig) EvaluateCampaignForRedemption(dbTrans DBTransaction) error {
	// Validate campaign for coupon redemption
	if err := c.ValidateForRedemption(); err != nil {
		return err
	}

	// Assume a transaction is always passed in for this operation
	// Update campaign records here using dbTrans
	// ...

	return nil
}

func (c *CampaignConfig) IsStartDateValid() error {
	if c.StartDate.Before(time.Now()) {
		return &DateError{"start date cannot be in the past"}
	}
	return nil
}

func (c *CampaignConfig) IsEndDateValid() error {
	if c.EndDate.Before(time.Now()) {
		return &DateError{"end date cannot be in the past"}
	}
	return nil
}

func (c *CampaignConfig) IsUsageLimitValid() error {
	if c.UsageLimit < 0 {
		return &LimitError{"usage limit cannot be negative"}
	}
	return nil
}

func (c *CampaignConfig) IsAvailabilityValidForOnDemandAndPregenerate() error {
	if c.AllowOnDemandCoupons && c.PregenerateCoupons && c.AvailabilityCount <= 0 {
		return &ValidationError{"availability count must be greater than 0 if pregenerating coupons and allowing on-demand coupons"}
	}
	return nil
}

func (c *CampaignConfig) IsUsageLimitSufficientForDemand() error {
	if c.AllowOnDemandCoupons && c.PregenerateCoupons && c.RedeemedCount+c.AvailabilityCount >= c.UsageLimit {
		return &ValidationError{"availability count and redeemed count must be less than usage limit if pregenerating coupons and allowing on-demand coupons"}
	}
	return nil
}

func (c *CampaignConfig) IsAvailabilityValidForPregeneration() error {
	if c.PregenerateCoupons && c.AvailabilityCount <= 0 {
		return &ValidationError{"availability count must be greater than 0 if pregenerating coupons"}
	}
	return nil
}

func (c *CampaignConfig) IsUsageLimitGreaterThanRedeemedCount() error {
	if c.UsageLimit < c.RedeemedCount {
		return &LimitError{"usage limit must be greater than or equal to redeemed count"}
	}
	return nil
}

func (c *CampaignConfig) IsCampaignIDValid() error {
	if c.ID == uuid.Nil {
		return &ValidationError{"campaign ID cannot be empty"}
	}
	return nil
}

func (c *CampaignConfig) AreDatesValid() error {
	if !c.StartDate.Before(c.EndDate) {
		return &DateError{"start date must be before end date"}
	}
	return nil
}

func (c *CampaignConfig) IsCampaignTypeValid() error {
	if c.CampaignType < CampaignTypeUnknown || c.CampaignType > maxCampaignType {
		return &CampaignTypeError{"invalid campaign type"}
	}
	return nil
}

func (c *CampaignConfig) IsUsageLimitValidForSingleUse() error {
	if !c.IsSingleUse && c.UsageLimit <= 1 {
		return &LimitError{"usage limit must be greater than 1 if coupons are not single use"}
	}
	if c.IsSingleUse && c.UsageLimit != 1 {
		return &LimitError{"usage limit must be 1 if coupons are single use"}
	}
	return nil
}

func (c *CampaignConfig) IsRedeemedCountValid() error {
	if c.RedeemedCount < 0 {
		return &LimitError{"redeemed count cannot be negative"}
	}
	return nil
}

func (c *CampaignConfig) IsLimitPerUserValid() error {
	if c.LimitPerUser < 0 {
		return &LimitError{"limit per user cannot be negative"}
	}
	return nil
}

func (c *CampaignConfig) IsActive() error {
	if c.IsCampaignActive {
		return nil
	}
	return &ValidationError{"campaign is not active"}
}

func (c *CampaignConfig) ValidateSHA256Index() error {
	if !c.IsSHA256IndexValid() {
		return &ValidationError{"invalid SHA256 index"}
	}
	return nil
}

func (c *CampaignConfig) IsExpired() bool {
	return time.Now().After(c.EndDate)
}

func (c *CampaignConfig) IsStarted() bool {
	return time.Now().After(c.StartDate)
}

func (c *CampaignConfig) IsPregenerated() bool {
	return c.PregenerateCoupons
}

// IsValidAtTime checks if the campaign is valid at a given time
func (c *CampaignConfig) IsValidAtTime(t time.Time) bool {
	// is t between start and end date
	return t.After(c.StartDate) && t.Before(c.EndDate)
}

func (c *CampaignConfig) IsAvailable() bool {
	return c.AvailabilityCount > 0
}

func (c *CampaignConfig) HasRemainingCoupons() bool {
	return c.AvailabilityCount > 0
}

func (c *CampaignConfig) CanServiceOnDemandRequests() bool {
	return c.AllowOnDemandCoupons
}

func (c *CampaignConfig) IsSHA256IndexValid() bool {
	return c.SHA256Index >= 0 && c.SHA256Index <= 31
}

func (c *CampaignConfig) Validate() error {
	validators := []func() error{
		c.IsCampaignIDValid,                    // common validation
		c.AreDatesValid,                        // common validation
		c.IsCampaignTypeValid,                  // common validation
		c.IsUsageLimitValid,                    // common validatio - Validate the usage limit is not negative
		c.IsUsageLimitGreaterThanRedeemedCount, // common validation
		c.IsUsageLimitValidForSingleUse,        // common validation
		c.IsRedeemedCountValid,                 // common validation
		c.IsLimitPerUserValid,                  // common validation
		c.ValidateSHA256Index,                  // common validation
	}

	for _, validator := range validators {
		if err := validator(); err != nil {
			return err
		}
	}

	return nil
}

func (c *CampaignConfig) Save() error {
	return nil
}

func (c *CampaignConfig) Delete() error {
	return nil
}

func (c *CampaignConfig) Update() error {
	return nil
}

func (c *CampaignConfig) GetCoupon() error {
	return nil
}

func (c *CampaignConfig) GetCoupons(num int) error {
	return nil
}
