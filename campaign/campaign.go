package campaign

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

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

	StrategyConfig json.RawMessage `json:"strategyConfig"`
}

// Assuming uuid, json, and time packages are imported correctly above
func (c *CampaignConfig) Validate() error {
	// You can't create a campaign in the past
	if c.StartDate.Before(time.Now()) {
		return &DateError{"start date cannot be in the past"}
	}

	// You can't create a campaign that ends in the past
	if c.EndDate.Before(time.Now()) {
		return &DateError{"end date cannot be in the past"}
	}

	// Check for valid usage limit
	if c.UsageLimit < 0 {
		return &LimitError{"usage limit cannot be negative"}
	}

	// If we allow on-demand and pregenerate, we need to ensure that the availability count is greater than 0
	if c.AllowOnDemandCoupons && c.PregenerateCoupons && c.AvailabilityCount <= 0 {
		return &ValidationError{"availability count must be greater than 0 if pregenerating coupons and allowing on-demand coupons"}
	}

	// If we allow on-demand coupons and pregenerate coupons, we need to ensure it lines up with redeemed count
	if c.AllowOnDemandCoupons && c.PregenerateCoupons && c.RedeemedCount+c.AvailabilityCount >= c.UsageLimit {
		return &ValidationError{"availability count and redeemed count must be less than usage limit if pregenerating coupons and allowing on-demand coupons"}
	}

	// Validate that pregenerated coupons are available
	if c.PregenerateCoupons && c.AvailabilityCount <= 0 {
		return &ValidationError{"availability count must be greater than 0 if pregenerating coupons"}
	}

	// Validate that usage limit is greater than or equal to redeemed count
	if c.UsageLimit < c.RedeemedCount {
		return &LimitError{"usage limit must be greater than or equal to redeemed count"}
	}

	if !c.IsActive() {
		return &ValidationError{"campaign is not active"}
	}
	// Check if the campaign ID is valid
	if c.ID == uuid.Nil {
		return &ValidationError{"campaign ID cannot be empty"}
	}

	// Check if campaign dates are valid
	if !c.StartDate.Before(c.EndDate) {
		return &DateError{"start date must be before end date"}
	}

	// Check if CampaignType is within the valid range
	if c.CampaignType < CampaignTypeUnknown || c.CampaignType > maxCampaignType {
		return &CampaignTypeError{"invalid campaign type"}
	}

	// Check for valid usage limits if IsSingleUse is false
	if !c.IsSingleUse && c.UsageLimit <= 1 {
		return &LimitError{"usage limit must be greater than 1 if coupons are not single use"}
	}

	// Check for valid usage limits if IsSingleUse is true
	if c.IsSingleUse && c.UsageLimit != 1 {
		return &LimitError{"usage limit must be 1 if coupons are single use"}
	}

	// Check for valid redeemed count
	if c.RedeemedCount < 0 {
		return &LimitError{"redeemed count cannot be negative"}
	}

	// Check for valid limit per user
	if c.LimitPerUser < 0 {
		return &LimitError{"limit per user cannot be negative"}
	}

	// Add any additional checks as necessary, for example, validating JSON fields with schema or ensuring values are within expected ranges

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

func (c *CampaignConfig) IsActive() bool {
	return c.IsCampaignActive
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

func (c *CampaignConfig) GetCoupon() error {
	return nil
}

func (c *CampaignConfig) GetCoupons(num int) error {
	return nil
}
