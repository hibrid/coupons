#### Overview
The `strategy` package defines the interfaces and structures necessary to implement various discount strategies for promotional campaigns.

#### Interfaces

- **DiscountStrategy**: Defines the methods required for a discount strategy, including validation, application, and result retrieval.
- **CampaignData**: An interface used by discount strategies to access necessary campaign configuration details without direct dependency on the campaign configuration structure.

#### Use Cases

- Implementing various discount strategies, such as BOGO (Buy One, Get One), percentage discounts, and fixed amount discounts.
- Validating input parameters for a discount strategy based on campaign configurations.
- Applying discount strategies to transactions, adjusting prices or quantities accordingly.

#### Example

```go
// Implementing a new discount strategy
package mystrategy

import (
    "github.com/hibrid/coupons/common"
    "github.com/hibrid/coupons/strategy"
)

type MyDiscountStrategy struct {
    // Strategy-specific fields
}

func (m *MyDiscountStrategy) ValidateInputs(data strategy.CampaignData) bool {
    // Validation logic here
    return true
}

func (m *MyDiscountStrategy) ApplyDiscount(data strategy.CampaignData) {
    // Discount application logic here
}

func (m *MyDiscountStrategy) GetDiscountResult() *common.DiscountResult {
    // Return the result of the discount application
    return &common.DiscountResult{
        // Populate the result
    }
}
```