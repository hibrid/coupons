#### Overview
The `campaign` package manages campaign configurations for various promotional activities, including discounts and offers. It allows campaigns to be configured, validated, and applied with flexible strategies.

#### Types

- **CampaignConfig**: Represents the configuration of a promotional campaign. It includes start and end dates, eligibility criteria, and a strategy for applying discounts or offers.

#### Interfaces

- **CampaignData**: Abstracts access to campaign configuration details necessary for validating and applying discount strategies.

#### Use Cases

- Creating and configuring new promotional campaigns.
- Validating campaign configurations.
- Applying discount strategies based on campaign rules.

#### Example

```go
package main

import (
    "github.com/hibrid/coupons/campaign"
    "github.com/hibrid/coupons/discount"
    "github.com/hibrid/coupons/strategy"
)

func main() {
    // Initialize a campaign configuration
    campaignConfig := &campaign.CampaignConfig{
        // Set necessary fields...
    }

    // Create a BOGO discount strategy
    bogoStrategy := discount.NewBogoDiscountStrategy(&discount.BOGODiscountStrategy{
        DiscountedSku: "SKU123",
        FreeSku:       "SKU456",
        DiscountedQty: 1,
        FreeQty:       1,
    })

    // Set the discount strategy for the campaign
    campaignConfig.SetDiscountStrategy(bogoStrategy)

    // Validate and apply the discount strategy
    if campaignConfig.ValidateStrategy() {
        result := campaignConfig.ApplyDiscountStrategy()
        // Process the result...
    }
}
```