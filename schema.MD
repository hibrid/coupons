Campaigns: Stores information about marketing campaigns.

id (Primary Key): Unique identifier for each campaign.
campaign_name: The name or title of the campaign.
start_date: The start date of the campaign.
end_date: The end date of the campaign.
is_active: A flag indicating if the campaign is currently active.
Coupons: Stores information about coupons.

id (Primary Key): Unique identifier for each coupon.
code: The coupon code, used by customers to apply the coupon.
description: A description or label for the coupon.
discount_type: The type of discount (e.g., percentage, fixed amount, trial period, recurring discount, fixed price subscription).
discount_value: The value of the discount, which could be a percentage or a fixed amount, depending on the discount type.
minimum_purchase: Minimum purchase amount required to use the coupon.
expiration_date: The date when the coupon expires.
is_single_use: A flag indicating if the coupon can be used only once.
usage_limit: The maximum number of times the coupon can be used.
is_active: A flag indicating if the coupon is currently active.
is_valid: A flag indicating if the coupon is valid. This can be used to enable or disable coupons without deleting them.
not_valid_reason: A textual description explaining why the coupon is not valid, if applicable.
campaign_id (Foreign Key): The ID of the campaign to which the coupon belongs.
trial_period_unit and trial_period_length: These fields define the unit (e.g., "MONTHLY", "ANNUAL") and length (integer) of the trial period for subscription-based coupons.
post_trial_pricing: This field holds the pricing after the trial period ends. It is applicable if the coupon offers a trial period for a subscription.
discount_duration_unit and discount_duration_length: These fields represent the duration unit and length for recurring discounts. This is relevant for coupons that provide a discount over multiple billing cycles.
fixed_price_duration_unit and fixed_price_duration_length: These fields are used to store the duration unit and length of a fixed price subscription period. This is applicable for coupons offering a subscription at a fixed price for a certain duration.
eligible_plans: A JSON-encoded array containing an ordered list of subscription plan identifiers that are eligible for the coupon. The array preserves the order specified by the user, determining which plans are applicable for the coupon and in what sequence.

id (Primary Key): Unique identifier for each SKU.
product_name: The name or title of the product.
product_description: A description of the product.
product_category: The category or type of the product.
SKU_Coupon_Mapping: Maps coupons to SKUs, indicating which coupons are valid for which products.

coupon_id (Foreign Key): The ID of the coupon.
sku_id (Foreign Key): The ID of the SKU.
CouponUsage: Records coupon usage for tracking and analytics.

id (Primary Key): Unique identifier for each coupon usage record.
coupon_id (Foreign Key): The ID of the coupon that was used.
user_id: The ID of the user who used the coupon.
order_id: The ID of the order associated with coupon usage.
usage_date: The date when the coupon was used.
is_used: A flag indicating if the coupon was used.
Referral: Stores information about referral relationships.

id (Primary Key): Unique identifier for each referral record.
referrer_id: The ID of the referrer (the user who made the referral).
referee_id: The ID of the referee (the user who was referred).
referral_date: The date when the referral occurred.
is_rewarded: A flag indicating if the referrer received a reward for the referral.


Rulesets: Create a table to store rulesets that define coupon validation rules.

id (Primary Key): Unique identifier for each ruleset.
name: A descriptive name for the ruleset.
definition: The definition of the ruleset, which contains the validation rules in a format compatible with your chosen rules engine (e.g., grule).
Campaign_Rulesets: Create a table to associate campaigns with rulesets. Each campaign can have one or more associated rulesets.

campaign_id (Foreign Key): The ID of the campaign.
ruleset_id (Foreign Key): The ID of the ruleset.
Coupon_Rulesets: Create a table to associate coupons with rulesets. Each coupon can have one or more associated rulesets.

coupon_id (Foreign Key): The ID of the coupon.
ruleset_id (Foreign Key): The ID of the ruleset.