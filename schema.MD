Campaigns: Stores information about marketing campaigns.

id (Primary Key): Unique identifier for each campaign.
campaign_name: The name or title of the campaign.
start_date: The start date of the campaign.
end_date: The end date of the campaign.
is_active: A flag indicating if the campaign is currently active.
Coupons: Stores information about coupons.

id (Primary Key): Unique identifier for each coupon.
code: The coupon code.
description: A description or label for the coupon.
discount_type: The type of discount (e.g., percentage, fixed amount).
discount_value: The value of the discount.
minimum_purchase: Minimum purchase amount required to use the coupon.
expiration_date: The date when the coupon expires.
is_single_use: A flag indicating if the coupon can be used only once.
usage_limit: The maximum number of times the coupon can be used.
is_active: A flag indicating if the coupon is currently active.
campaign_id (Foreign Key): The ID of the campaign to which the coupon belongs.
SKU: Stores information about products or SKUs.

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