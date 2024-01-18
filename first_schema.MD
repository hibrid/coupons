Certainly! Below are the MySQL `CREATE TABLE` statements for the proposed database layout, including optimized indexes and keys:

**Table: Campaigns**

```sql
CREATE TABLE Campaigns (
    campaign_id INT AUTO_INCREMENT PRIMARY KEY,
    campaign_name VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**Table: Coupons**

```sql
CREATE TABLE Coupons (
    coupon_id INT AUTO_INCREMENT PRIMARY KEY,
    coupon_code VARCHAR(20) NOT NULL UNIQUE,
    coupon_description TEXT,
    discount_type ENUM('percentage', 'fixed_amount', 'other') NOT NULL,
    discount_value DECIMAL(10, 2) NOT NULL,
    minimum_purchase_amount DECIMAL(10, 2),
    expiration_date DATE NOT NULL,
    is_single_use BOOLEAN NOT NULL DEFAULT FALSE,
    usage_limit INT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    campaign_id INT,  -- Foreign key reference to Campaigns table
    FOREIGN KEY (campaign_id) REFERENCES Campaigns(campaign_id)
);
```

**Table: SKU_Coupon_Mapping**

```sql
CREATE TABLE SKU_Coupon_Mapping (
    mapping_id INT AUTO_INCREMENT PRIMARY KEY,
    coupon_id INT,
    sku_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coupon_id) REFERENCES Coupons(coupon_id),
    FOREIGN KEY (sku_id) REFERENCES SKU(sku_id),
    UNIQUE KEY (coupon_id, sku_id)  -- Ensures uniqueness of coupon-SKU associations
);
```

**Table: SKU**

```sql
CREATE TABLE SKU (
    sku_id INT AUTO_INCREMENT PRIMARY KEY,
    product_name VARCHAR(255) NOT NULL,
    product_description TEXT,
    product_category VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**Table: RuleSets**

```sql
CREATE TABLE RuleSets (
    ruleset_id INT AUTO_INCREMENT PRIMARY KEY,
    ruleset_name VARCHAR(255) NOT NULL,
    ruleset_definition TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

**Table: Coupon_Ruleset_Mapping**

```sql
CREATE TABLE Coupon_Ruleset_Mapping (
    mapping_id INT AUTO_INCREMENT PRIMARY KEY,
    coupon_id INT,
    ruleset_id INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (coupon_id) REFERENCES Coupons(coupon_id),
    FOREIGN KEY (ruleset_id) REFERENCES RuleSets(ruleset_id),
    UNIQUE KEY (coupon_id, ruleset_id)  -- Ensures uniqueness of coupon-ruleset associations
);
```

These SQL statements define the tables and establish the necessary relationships with foreign keys. Additionally, the `UNIQUE KEY` constraints on specific columns ensure the uniqueness of certain associations, preventing duplication of coupon-SKU and coupon-ruleset mappings. These indexes and keys help optimize data integrity and query performance in your database.