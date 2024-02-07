package campaign

import (
	"encoding/json"
	"errors"
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

	// you can't create a campaign in the past
	if c.StartDate.Before(time.Now()) {
		return errors.New("start date cannot be in the past")
	}

	// you can't create a campaign that end in the past
	if c.EndDate.Before(time.Now()) {
		return errors.New("end date cannot be in the past")
	}

	// if we allow on demand and pregenerate, we need to ensure that the availability count is greater than 0
	if c.AllowOnDemandCoupons && c.PregenerateCoupons && c.AvailabilityCount <= 0 {
		return errors.New("availability count must be greater than 0 if pregenerating coupons and allowing on demand coupons")
	}

	// if we allow on demand coupons, we need and we pregenerate coupons, we need to ensure it lines up with redeemed count
	if c.AllowOnDemandCoupons && c.PregenerateCoupons && c.RedeemedCount+c.AvailabilityCount < c.UsageLimit {
		return errors.New("availability count and redeemed count must be less than usage limit if pregenerating coupons and allowing on demand coupons")
	}

	// validate that pregenerated coupons are available
	if c.PregenerateCoupons && c.AvailabilityCount <= 0 {
		return errors.New("availability count must be greater than 0 if pregenerating coupons")
	}

	// validate that usage limit is greater than or equal to redeemed count
	if c.UsageLimit < c.RedeemedCount {
		return errors.New("usage limit must be greater than or equal to redeemed count")
	}

	if !c.IsActive() {
		return errors.New("campaign is not active")
	}
	// Check if the campaign ID is valid
	if c.ID == uuid.Nil {
		return errors.New("campaign ID cannot be empty")
	}

	// Check if campaign dates are valid
	if !c.StartDate.Before(c.EndDate) {
		return errors.New("start date must be before end date")
	}

	// Check if the availability count is valid
	if c.AvailabilityCount < 0 {
		return errors.New("availability count cannot be negative")
	}

	// Check if CampaignType is within the valid range
	if c.CampaignType < CampaignTypeUnknown || c.CampaignType > maxCampaignType {
		return errors.New("invalid campaign type")
	}

	// If pregenerateCoupons is true, ensure that AvailabilityCount is greater than 0
	if c.PregenerateCoupons && c.AvailabilityCount <= 0 {
		return errors.New("availability count must be greater than 0 if pregenerating coupons")
	}

	// Check for valid usage limits if IsSingleUse is false
	if !c.IsSingleUse && c.UsageLimit <= 1 {
		return errors.New("usage limit must be greater than 1 if coupons are not single use")
	}

	// Check for valid limit per user
	if c.LimitPerUser < 0 {
		return errors.New("limit per user cannot be negative")
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
