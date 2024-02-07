package campaign_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hibrid/coupons/campaign"
)

func TestCampaignConfig_Validate(t *testing.T) {
	// Define a valid campaign ID for reuse
	validCampaignID := uuid.New()

	// Define test cases
	tests := []struct {
		name    string
		config  campaign.CampaignConfig
		wantErr bool
	}{
		{
			name: "valid configuration",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
			},
			wantErr: false,
		},
		{
			name: "negative availability count",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  -1,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       1,
			},
			wantErr: true,
		},
		{
			name: "negative limit per user count",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         2,
				LimitPerUser:       -1,
			},
			wantErr: true,
		},
		{
			name: "invalid campaign ID",
			config: campaign.CampaignConfig{
				ID: uuid.Nil,
			},
			wantErr: true,
		},
		{
			name: "start date after end date",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(24 * time.Hour),
				EndDate:            time.Now().Add(-24 * time.Hour),
				CampaignType:       campaign.CampaignType(-1), // intentionally invalid
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       1,
			},

			wantErr: true,
		},
		{
			name: "invalid campaign type",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignType(-1), // intentionally invalid
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       1,
			},

			wantErr: true,
		},
		{
			name: "negative limit per user",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignType(-1), // intentionally invalid
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       -1,
			},
			wantErr: true,
		},
		{
			name: "invalid usage limit",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignType(-1), // intentionally invalid
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        true,
				UsageLimit:         0,
				LimitPerUser:       1,
			},
			wantErr: true,
		},
		{
			name: "invalid availability count",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  0,
				IsSingleUse:        true,
				UsageLimit:         1,
				LimitPerUser:       1,
			},

			wantErr: true,
		},
		{
			name: "invalid usage limit for single use coupons",
			config: campaign.CampaignConfig{
				ID:                 validCampaignID,
				StartDate:          time.Now().Add(-24 * time.Hour),
				EndDate:            time.Now().Add(24 * time.Hour),
				CampaignType:       campaign.CampaignTypePromoCode,
				PregenerateCoupons: true,
				AvailabilityCount:  10,
				IsSingleUse:        false,
				UsageLimit:         1,
				LimitPerUser:       1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("CampaignConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
