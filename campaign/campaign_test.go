package campaign_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hibrid/coupons/campaign"
	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestCampaignConfig_Validate(t *testing.T) {
	// Define a valid campaign ID for reuse
	validCampaignID := uuid.New()

	// Define test cases
	tests := []struct {
		name        string
		config      campaign.CampaignConfig
		wantErr     bool
		wantErrType error // Add this field for checking error type
	}{
		{
			name: "valid configuration",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr: false,
		},
		{
			name: "usage limit cannot be negative",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         -1,
				LimitPerUser:       1,
				IsCampaignActive:   true,
				RedeemedCount:      -1,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"usage limit cannot be negative"},
		},
		{
			name: "redeemed count cannot be negative",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
				RedeemedCount:      -1,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"redeemed count cannot be negative"},
		},
		{
			name: "usage limit must be 1 if coupons are single use",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        true,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"usage limit must be 1 if coupons are single use"},
		},
		{
			name: "usage limit must be greater than 1 if coupons are not single use",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"usage limit must be greater than 1 if coupons are not single use"},
		},
		{
			name: "availability count must be greater than 0 if pregenerating coupons",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  0,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.ValidationError{"availability count must be greater than 0 if pregenerating coupons"},
		},
		{
			name: "usage limit must be greater than or equal to redeemed count",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       1,
				IsCampaignActive:   true,
				RedeemedCount:      2,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"usage limit must be greater than or equal to redeemed count"},
		},
		{
			name: "availability count and redeemed count must be less than usage limit if pregenerating coupons and allowing on-demand coupons",
			config: campaign.CampaignConfig{
				ID:                   validCampaignID,
				StartDate:            time.Now().Add(24 * time.Hour),
				EndDate:              time.Now().Add(48 * time.Hour),
				CampaignType:         campaign.CampaignTypePromoCode,
				PregenerateCoupons:   true,
				AvailabilityCount:    9,
				IsSingleUse:          false,
				UsageLimit:           2,
				LimitPerUser:         1,
				IsCampaignActive:     true,
				AllowOnDemandCoupons: true,
				RedeemedCount:        1,
			},
			wantErr:     true,
			wantErrType: &campaign.ValidationError{"availability count and redeemed count must be less than usage limit if pregenerating coupons and allowing on-demand coupons"},
		},
		{
			name: "availability count must be greater than 0 if pregenerating coupons and allowing on-demand coupons",
			config: campaign.CampaignConfig{
				ID:                   validCampaignID,
				StartDate:            time.Now().Add(24 * time.Hour),
				EndDate:              time.Now().Add(48 * time.Hour),
				CampaignType:         campaign.CampaignTypePromoCode,
				PregenerateCoupons:   true,
				AvailabilityCount:    0,
				IsSingleUse:          false,
				UsageLimit:           2,
				LimitPerUser:         1,
				IsCampaignActive:     true,
				AllowOnDemandCoupons: true,
			},
			wantErr:     true,
			wantErrType: &campaign.ValidationError{"availability count must be greater than 0 if pregenerating coupons and allowing on-demand coupons"},
		},
		{
			name: "start date cannot be in the past",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.DateError{"start date cannot be in the past"},
		},
		{
			name: "campaign is not active",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   false,
			},
			wantErr:     true,
			wantErrType: &campaign.ValidationError{"campaign is not active"},
		},
		{
			name: "negative availability count",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  -1,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.ValidationError{"availability count must be greater than 0 if pregenerating coupons"},
		},
		{
			name: "negative limit per user count",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       -1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"limit per user cannot be negative"},
		},
		{
			name: "invalid campaign ID",
			config: campaign.CampaignConfig{
				ID:                 uuid.Nil,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.ValidationError{"campaign ID cannot be empty"},
		},
		{
			name: "end data cannot be in the past",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(-24 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},

			wantErr:     true,
			wantErrType: &campaign.DateError{"end date cannot be in the past"},
		},
		{
			name: "start date after end date",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(48 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.DateError{"start date must be before end date"},
		},
		{
			name: "invalid campaign type",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignType(-1), // intentionally invalid
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},

			wantErr:     true,
			wantErrType: &campaign.CampaignTypeError{"invalid campaign type"},
		},
		{
			name: "negative limit per user",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignType(-1), // intentionally invalid
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       -1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.CampaignTypeError{"invalid campaign type"},
		},
		{
			name: "invalid usage limit",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignType(1),
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         0,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"usage limit must be greater than 1 if coupons are not single use"},
		},
		{
			name: "invalid availability count",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  0,
				IsSingleUse:        true,
				UsageLimit:         1,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},

			wantErr:     true,
			wantErrType: &campaign.ValidationError{"availability count must be greater than 0 if pregenerating coupons"},
		},
		{
			name: "invalid usage limit for single use coupons",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(48 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       1,
				IsCampaignActive:   true,
			},
			wantErr:     true,
			wantErrType: &campaign.LimitError{"usage limit must be greater than 1 if coupons are not single use"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CampaignConfig.Validate() error = %v, wantErr %v, ", err, tt.wantErr)
				return
			}

			// Check error type and text if an error is expected
			if tt.wantErr {
				if !errors.Is(err, tt.wantErrType) {
					dmp := diffmatchpatch.New()

					diffs := dmp.DiffMain(err.Error(), tt.wantErrType.Error(), false)
					t.Errorf("error type = %T, want %T\n Got  = %s\n Want = %s\n text diff:%s", err, tt.wantErrType, err, tt.wantErrType, diffs)
				}

			}
		})
	}
}
