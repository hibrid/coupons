-- Create the Campaigns table
CREATE TABLE Campaigns (
    id INT AUTO_INCREMENT PRIMARY KEY,
    campaign_name VARCHAR(255) NOT NULL,
    coupon_type VARCHAR(255) NOT NULL,
    coupon_vanity_name VARCHAR(255) NOT NULL,
    discount_type VARCHAR(50) NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    minimum_purchase DECIMAL(10, 2) NOT NULL,
    is_single_use BOOLEAN NOT NULL,
    usage_limit INT NOT NULL,
    expire_after_days INT NOT NULL,
    expire_after_hours INT NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN NOT NULL
);

-- Create the Coupons table
CREATE TABLE Coupons (
        id INT AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL,
    'description' VARCHAR(255) NOT NULL,
    coupon_type VARCHAR(255) NOT NULL,
	coupon_vanity_name VARCHAR(255) NOT NULL,
    discount_type VARCHAR(50) NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    minimum_purchase DECIMAL(10, 2) NOT NULL,
    expiration_date DATE NOT NULL,
    is_single_use BOOLEAN NOT NULL,
    usage_limit INT NOT NULL,
    is_active BOOLEAN NOT NULL,
    is_valid BOOLEAN NOT NULL,
    not_valid_reason VARCHAR(255),
    campaign_id INT NOT NULL,

    -- fields for standardized billing periods and subscription-related information
    trial_period_unit VARCHAR(20),
    trial_period_length INT,
    post_trial_pricing DECIMAL(10, 2),
    discount_duration_unit VARCHAR(20),
    discount_duration_length INT,
    fixed_price_duration_unit VARCHAR(20),
    fixed_price_duration_length INT,
    eligible_plans JSON,

    FOREIGN KEY (campaign_id) REFERENCES Campaigns(id),
    INDEX idx_code (code),
    INDEX idx_expiration_date (expiration_date)
);

-- Create the SKU table
CREATE TABLE SKU (
    id INT AUTO_INCREMENT PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_category VARCHAR(255)
);

-- Create the SKU_Coupon_Mapping table
CREATE TABLE SKU_Coupon_Mapping (
    coupon_id INT NOT NULL,
    sku_id INT NOT NULL,
    PRIMARY KEY (coupon_id, sku_id),
    FOREIGN KEY (coupon_id) REFERENCES Coupons(id),
    FOREIGN KEY (sku_id) REFERENCES SKU(id)
);

-- Create the CouponUsage table
CREATE TABLE CouponUsage (
    id INT AUTO_INCREMENT PRIMARY KEY,
    coupon_id INT NOT NULL,
    user_id INT NOT NULL,
    order_id INT NOT NULL,
    usage_date DATETIME NOT NULL,
    is_used BOOLEAN NOT NULL,
    FOREIGN KEY (coupon_id) REFERENCES Coupons(id),
    INDEX idx_usage_date (usage_date)
);

-- Create the Referral table
CREATE TABLE Referral (
    id INT AUTO_INCREMENT PRIMARY KEY,
    referrer_id INT NOT NULL,
    referee_id INT NOT NULL,
    referral_date DATETIME NOT NULL,
    is_rewarded BOOLEAN NOT NULL,
    FOREIGN KEY (referrer_id) REFERENCES Users(id),
    FOREIGN KEY (referee_id) REFERENCES Users(id),
    INDEX idx_referral_date (referral_date)
);

-- Create the Rulesets table
CREATE TABLE Rulesets (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    definition TEXT NOT NULL
);

-- Create the Campaign_Rulesets table to associate campaigns with rulesets
CREATE TABLE Campaign_Rulesets (
    campaign_id INT NOT NULL,
    ruleset_id INT NOT NULL,
    PRIMARY KEY (campaign_id, ruleset_id),
    FOREIGN KEY (campaign_id) REFERENCES Campaigns(id),
    FOREIGN KEY (ruleset_id) REFERENCES Rulesets(id)
);

-- Create the Coupon_Rulesets table to associate coupons with rulesets
CREATE TABLE Coupon_Rulesets (
    coupon_id INT NOT NULL,
    ruleset_id INT NOT NULL,
    PRIMARY KEY (coupon_id, ruleset_id),
    FOREIGN KEY (coupon_id) REFERENCES Coupons(id),
    FOREIGN KEY (ruleset_id) REFERENCES Rulesets(id)
);