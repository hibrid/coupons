## Table of Contents

1. [Types and Definitions](#types-and-definitions)
2. [Validations and Required Fields](#validations-and-required-fields)
3. [Computations](#computations)
4. [Use Cases](#use-cases)
    - [Subscriptions](#subscriptions)
    - [Non-Subscriptions](#non-subscriptions)
5. [Examples](#examples)

---

## Types and Definitions

### DiscountType

| Field          | Description                              |
|----------------|------------------------------------------|
| Unknown        | Default value for discount type          |
| Percentage     | Discount represented as a percentage     |
| TimeBased      | Discount applied over a specific time period |
| FixedAmount    | Discount represented as a fixed amount  |

### DiscountApplication

| Field          | Description                              |
|----------------|------------------------------------------|
| UnknownApplication | Default value for discount application |
| Recurring      | Discount applied to recurring charges   |
| Spread         | Discount spread over a specific period  |
| OneTime        | One-time discount applied to the first billing cycle |

### DiscountPhase

| Field                         | Description                                                     |
|-------------------------------|-----------------------------------------------------------------|
| Duration                      | Duration in units of TimeUnit                                   |
| DurationUnit                  | Unit of time for duration                                       |
| DiscountValue                 | Value of the discount (percentage or fixed amount)              |
| DiscountType                  | Type of discount                                                |
| Description                   | Description of the discount phase                                |
| ApplicableNumberOfBillingCycles | Number of billing cycles the discount is applicable to         |
| Application                   | Type of application (Recurring, Spread, OneTime)                 |
| Logs                          | Logs containing information about discount calculations          |
| DiscountsPerBillingCycle      | Map of discounts per billing cycle                               |

### SubscriptionInfo

| Field            | Description                                |
|------------------|--------------------------------------------|
| IsRecurring      | Indicates if the subscription is recurring |
| BillingPeriodUnit| Unit of time for the billing cycle         |
| TrialPeriod      | Duration of the trial period               |
| TrialPeriodUnit  | Unit of time for the trial period          |
| DiscountPhases   | List of discount phases for the subscription|

### CartItem

| Field                           | Description                                     |
|---------------------------------|-------------------------------------------------|
| SkuID                           | SKU ID                                          |
| Quantity                        | Quantity of items                               |
| UnitPrice                       | Price per unit                                  |
| DiscountDescription             | Description of the discount                     |
| DiscountValuePerDiscountedUnit  | Discount amount per discounted unit             |
| NumberOfUnitsDiscounted         | Number of units discounted                      |
| DiscountType                    | Type of discount (percentage or fixed amount)   |
| IsSubscription                  | Indicates if the item is a subscription         |
| Subscription                    | Subscription information                        |

## Validations and Required Fields

- **DiscountPhase**
    - `DiscountValue`: Positive value required, must be greater than 0.
    - `Duration`: Must be a positive integer.
    - `DurationUnit`: Must be a valid time unit.
    - `ApplicableNumberOfBillingCycles`: Must be a positive integer.
    - `Application`: Must be a valid discount application type.

- **CartItem**
    - `Quantity`: Must be a non-negative integer.
    - `UnitPrice`: Must be a non-negative decimal.
    - `DiscountValuePerDiscountedUnit`: Must be a non-negative decimal.
    - `NumberOfUnitsDiscounted`: Must be a non-negative integer.

## Computations

- **calculateTotalDiscountAmount()**: Calculates the total discount amount based on the discount type and duration.
- **calculateDiscountForPhase()**: Calculates the discount for a specific phase based on the discount type and duration.
- **applyRecurringDiscount()**: Applies a recurring discount over multiple billing cycles.
- **applyNonRecurringDiscount()**: Applies a one-time or spread discount over a specific period.
- **applyFixedAmountDiscount()**: Applies a fixed amount discount for one-time or recurring applications.

## Use Cases

#### Subscriptions

1. **Percentage Discount**:
   - **Description**: Applies a percentage discount to the subscription amount for a specified duration.
   - **Example Code**:
     ```go
     // Creating a percentage discount phase for a subscription
     discountPhase := common.DiscountPhase{
         DiscountValue:                   20,                           // 20% discount
         Duration:                        3,                            // Duration of discount in units
         DurationUnit:                    common.TimeUnit(common.Month), // Duration unit (e.g., month)
         ApplicableNumberOfBillingCycles: 3,                            // Number of billing cycles discount applies to
         Application:                     common.Recurring,              // Discount application (Recurring)
     }
     ```

2. **Time-Based Discount**:
   - **Description**: Offers a full discount for a specific duration of the subscription.
   - **Example Code**:
     ```go
     // Creating a time-based discount phase for a subscription
     discountPhase := common.DiscountPhase{
         DiscountValue:                   100,                          // 100% discount
         Duration:                        6,                            // Duration of discount in units
         DurationUnit:                    common.TimeUnit(common.Month), // Duration unit (e.g., month)
         ApplicableNumberOfBillingCycles: 6,                            // Number of billing cycles discount applies to
         Application:                     common.Recurring,              // Discount application (Recurring)
     }
     ```

3. **Fixed Amount Discount**:
   - **Description**: Provides a fixed discount amount for each billing cycle of the subscription.
   - **Example Code**:
     ```go
     // Creating a fixed amount discount phase for a subscription
     discountPhase := common.DiscountPhase{
         DiscountValue:                   10,                           // $10 discount
         ApplicableNumberOfBillingCycles: 5,                            // Number of billing cycles discount applies to
         Application:                     common.Recurring,              // Discount application (Recurring)
     }
     ```

#### Non-Subscriptions

1. **Percentage Discount**:
   - **Description**: Applies a percentage discount to the total amount of non-subscription items.
   - **Example Code**:
     ```go
     // Creating a percentage discount for non-subscription items
     cartItem := common.CartItem{
         DiscountValuePerDiscountedUnit: decimal.NewFromFloat(15), // 15% discount
         NumberOfUnitsDiscounted:        5,                        // Number of units to apply discount to
         DiscountType:                   common.Percentage,        // Discount type (Percentage)
         Application:                    common.OneTime,           // Discount application (One-Time)
     }
     ```

2. **Fixed Amount Discount**:
   - **Description**: Offers a fixed discount amount for non-subscription items.
   - **Example Code**:
     ```go
     // Creating a fixed amount discount for non-subscription items
     cartItem := common.CartItem{
         DiscountValuePerDiscountedUnit: decimal.NewFromFloat(5), // $5 discount
         NumberOfUnitsDiscounted:        3,                       // Number of units to apply discount to
         DiscountType:                   common.FixedAmount,       // Discount type (Fixed Amount)
         Application:                    common.OneTime,          // Discount application (One-Time)
     }
     ```

