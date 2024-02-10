#### Overview
The `discount` package includes implementations of specific discount strategies, such as BOGO (Buy One, Get One), that can be applied to promotional campaigns.

#### Types

- **BOGODiscountStrategy**: An implementation of the `DiscountStrategy` interface for applying a Buy One, Get One discount.

#### Use Cases

- Configuring and applying a BOGO discount to a campaign.
- Validating the prerequisites for a BOGO discount based on campaign configurations.

#### Example

```go
package main

import (
    "github.com/hibrid/coupons/campaign"
    "github.com/hibrid/coupons/discount"
)

func main() {
    // Create and configure a BOGO discount strategy
    bogoStrategy := discount.NewBogoDiscountStrategy(&discount.BOGODiscountStrategy{
        DiscountedSku: "SKU123",
        FreeSku:       "SKU456",
        DiscountedQty: 1,
        FreeQty:       1,
    })

    // Initialize a campaign with the BOGO strategy
    campaignConfig := &campaign.CampaignConfig{
        // Campaign configuration...
    }
    campaignConfig.SetDiscountStrategy(bogoStrategy)

    // Validate and apply the strategy
    if campaignConfig.ValidateStrategy() {
        result := campaignConfig.ApplyDiscountStrategy()
        // Process the discount result...
    }
}
```