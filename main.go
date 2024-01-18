package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

// Define a struct to represent the Campaign data
type Campaign struct {
	ID        int
	Name      string
	StartDate string
	EndDate   string
	IsActive  bool
}

// Define a struct to represent the Coupon data
type Coupon struct {
	ID              int
	Code            string
	Description     string
	DiscountType    string
	DiscountValue   float64
	MinimumPurchase float64
	ExpirationDate  string
	IsSingleUse     bool
	UsageLimit      int
	IsActive        bool
	CampaignID      int
	IsValid         bool
	NotValidReason  string
}

func (mf *Coupon) IsNewCustomer(IsNewCustomer bool) string {
	if IsNewCustomer {
		return fmt.Sprintf("this is a new customer")
	}

	return fmt.Sprintf("this is not a new customer")
}

// Define a struct to represent SKU data
type SKU struct {
	ID                 int
	ProductName        string
	ProductDescription string
	ProductCategory    string
}

// Define a struct to represent SKU-Coupon mappings
type SKUToCouponMapping struct {
	CouponID int
	SKUID    int
}

// Define a struct to represent coupon usage
type CouponUsage struct {
	ID        int
	CouponID  int
	UserID    int
	OrderID   int
	UsageDate string
	IsUsed    bool
}

// Define a struct to represent referral data
type Referral struct {
	ID           int
	ReferrerID   int
	RefereeID    int
	ReferralDate string
	IsRewarded   bool
}

type RedemptionHistoryEntry struct {
	CouponCode     string
	RedemptionDate string
}

type OptionsContext struct {
	Option1 bool
	Option2 string
	Option3 int
}

type Subscription struct {
	ID               int
	UserID           int
	PlanID           int
	Term             string // e.g., monthly, annual, etc.
	MonthlyPrice     float64
	StartDate        time.Time
	EndDate          time.Time
	IsActive         bool
	SubscriptionType string // e.g., basic, premium, etc.
	// Add other relevant subscription details as needed
}

type ChangeContext struct {
	FromSubscription Subscription
	ToSubscription   Subscription
}

type CustomerContext struct {
	IsSubscriber        bool
	PreviousRedemptions []RedemptionHistoryEntry
	IsNewCustomer       bool
	UserID              int
	// Add other relevant customer data here
}

/*
func main() {
	// Database connection parameters
	db, err := sql.Open("mysql", "username:password@tcp(localhost:3306)/your_database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Test the database connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Example: Insert a new campaign
	newCampaign := Campaign{
		Name:      "Summer Sale",
		StartDate: "2023-06-01",
		EndDate:   "2023-06-30",
		IsActive:  true,
	}
	campaignID, err := InsertCampaign(db, newCampaign)
	if err != nil {
		log.Printf("Error inserting campaign: %v", err)
		return
	}

	// Example: Generate coupons in bulk based on configuration
	couponConfig := CouponConfig{
		CouponPrefix:   "SUMMER",
		CouponCount:    10,
		DiscountType:   "percentage",
		DiscountValue:  25.0,
		ExpirationDate: "2023-06-30",
		IsSingleUse:    false,
		UsageLimit:     100,
		IsActive:       true,
		CampaignID:     campaignID,
	}
	generatedCoupons, err := GenerateCoupons(db, couponConfig)
	if err != nil {
		log.Printf("Error generating coupons: %v", err)
		return
	}

	// Example: Insert SKUs
	sku1 := SKU{
		ProductName:        "Product 1",
		ProductDescription: "Description of Product 1",
		ProductCategory:    "Category A",
	}
	sku2 := SKU{
		ProductName:        "Product 2",
		ProductDescription: "Description of Product 2",
		ProductCategory:    "Category B",
	}
	sku1ID := InsertSKU(db, sku1)
	sku2ID := InsertSKU(db, sku2)

	// Example: Map coupons to SKUs
	MapCouponsToSKUs(db, generatedCoupons, []int{sku1ID, sku2ID})

	// Example: Record coupon usage
	user1ID := 1    // Replace with a valid user ID
	order1ID := 101 // Replace with a valid order ID
	RecordCouponUsage(db, generatedCoupons[0].ID, user1ID, order1ID)

	// Example: Send coupon expiration notifications
	SendCouponExpirationNotifications(db, generatedCoupons)

	// Example: Define and apply a ruleset for coupon validation
	ruleset := RuleSet{
		Name: "SummerSaleRules",
		Definition: `
            rule "MinimumPurchaseRule"
                when
                    $coupon : Coupon($minimumPurchase : minimum_purchase_amount)
                    $order : Order(totalAmount >= $minimumPurchase)
                then
                    $coupon.setValid(true);
            end
        `,
	}
	ApplyRuleset(ruleset, generatedCoupons)

	// Example: Implement referral system
	user2ID := 2 // Replace with a valid user ID
	ImplementReferralSystem(db, user1ID, user2ID)
}
*/

func main() {
	// Example: Define and apply a ruleset for coupon validation
	ruleset := RuleSet{
		Name:    "SummerSaleRules",
		Version: "1.0.0",
		Definition: `
		rule SummerSaleRules "Check New Subscriber" salience 5 {
			when
			CustomerContext.IsSubscriber != true && CustomerContext.IsNewCustomer == true
			Then
				Coupon.IsValid = true;
				Retract("SummerSaleRules");
		}
        `,
	}
	rulesets := []RuleSet{ruleset}

	for i := 0; i < 10; i++ {
		coupon := GenerateRandomCoupon()
		generatedCoupons := []Coupon{coupon}
		customerContext := GenerateRandomCustomerContext()
		changeContext := GenerateRandomChangeContext()
		optionsContext := GenerateRandomOptionsContext()
		ApplyRuleset(rulesets, generatedCoupons, customerContext, changeContext, optionsContext)
	}

}

// Insert a new campaign into the Campaigns table and return the campaign ID
func InsertCampaign(db *sql.DB, campaign Campaign) (int, error) {
	stmt, err := db.Prepare("INSERT INTO Campaigns (campaign_name, start_date, end_date, is_active) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(campaign.Name, campaign.StartDate, campaign.EndDate, campaign.IsActive)
	if err != nil {
		return 0, err
	}

	campaignID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	fmt.Println("Campaign inserted successfully")
	return int(campaignID), nil
}

// Define a struct to represent coupon generation configuration
type CouponConfig struct {
	CouponPrefix   string
	CouponCount    int
	DiscountType   string
	DiscountValue  float64
	ExpirationDate string
	IsSingleUse    bool
	UsageLimit     int
	IsActive       bool
	CampaignID     int
}

// Generate coupons in bulk based on configuration and return the generated coupons
func GenerateCoupons(db *sql.DB, config CouponConfig) ([]Coupon, error) {
	var generatedCoupons []Coupon

	// Prepare the SQL statement
	stmt, err := db.Prepare("INSERT INTO Coupons (coupon_code, coupon_description, discount_type, discount_value, minimum_purchase_amount, expiration_date, is_single_use, usage_limit, is_active, campaign_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	for i := 1; i <= config.CouponCount; i++ {
		couponCode := fmt.Sprintf("%s%d", config.CouponPrefix, i)
		newCoupon := Coupon{
			Code:           couponCode,
			Description:    fmt.Sprintf("%s Coupon %d", config.CouponPrefix, i),
			DiscountType:   config.DiscountType,
			DiscountValue:  config.DiscountValue,
			ExpirationDate: config.ExpirationDate,
			IsSingleUse:    config.IsSingleUse,
			UsageLimit:     config.UsageLimit,
			IsActive:       config.IsActive,
			CampaignID:     config.CampaignID,
		}

		// Execute the prepared statement with each coupon's data
		_, err := stmt.Exec(newCoupon.Code, newCoupon.Description, newCoupon.DiscountType, newCoupon.DiscountValue, newCoupon.MinimumPurchase, newCoupon.ExpirationDate, newCoupon.IsSingleUse, newCoupon.UsageLimit, newCoupon.IsActive, newCoupon.CampaignID)
		if err != nil {
			return nil, err
		}

		generatedCoupons = append(generatedCoupons, newCoupon)
	}

	fmt.Printf("%d Coupons generated successfully\n", config.CouponCount)
	return generatedCoupons, nil
}

/*
// Example usage
	daysUntilExpiration := 10
	coupons, err := findCouponsExpiringInDays(db, daysUntilExpiration)
	if err != nil {
		fmt.Println("Error finding coupons:", err)
		return
	}

	for _, coupon := range coupons {
		fmt.Printf("Coupon expiring in %d days: %+v\n", daysUntilExpiration, coupon)
	}
*/

func FindCouponsExpiringInDays(db *sql.DB, days int) ([]Coupon, error) {
	// Calculate the date range
	today := time.Now()
	endDate := today.AddDate(0, 0, days)

	// Prepare the SQL query
	query := `SELECT id, code, description, discount_type, discount_value, minimum_purchase, expiration_date, is_single_use, usage_limit, is_active, campaign_id FROM Coupons WHERE expiration_date BETWEEN ? AND ?`

	// Execute the query
	rows, err := db.Query(query, today.Format("2006-01-02"), endDate.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the coupons
	var coupons []Coupon
	for rows.Next() {
		var coupon Coupon
		err := rows.Scan(&coupon.ID, &coupon.Code, &coupon.Description, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MinimumPurchase, &coupon.ExpirationDate, &coupon.IsSingleUse, &coupon.UsageLimit, &coupon.IsActive, &coupon.CampaignID)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return coupons, nil
}

/*
// Example usage of findExpiredCoupons
	startDate := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC)
	expiredCoupons, err := findExpiredCoupons(db, startDate, endDate)
	if err != nil {
		fmt.Println("Error finding expired coupons:", err)
		return
	}

	// Do something with expiredCoupons
	for _, coupon := range expiredCoupons {
		fmt.Printf("Expired Coupon: %+v\n", coupon)
	}
*/

func FindExpiredCoupons(db *sql.DB, start, end time.Time) ([]Coupon, error) {
	// Prepare the SQL query
	query := `SELECT id, code, description, discount_type, discount_value, minimum_purchase, expiration_date, is_single_use, usage_limit, is_active, campaign_id FROM Coupons WHERE expiration_date BETWEEN ? AND ?`

	// Execute the query
	rows, err := db.Query(query, start.Format("2006-01-02"), end.Format("2006-01-02"))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the rows and populate the coupons
	var coupons []Coupon
	for rows.Next() {
		var coupon Coupon
		err := rows.Scan(&coupon.ID, &coupon.Code, &coupon.Description, &coupon.DiscountType, &coupon.DiscountValue, &coupon.MinimumPurchase, &coupon.ExpirationDate, &coupon.IsSingleUse, &coupon.UsageLimit, &coupon.IsActive, &coupon.CampaignID)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, coupon)
	}

	// Check for errors encountered during iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return coupons, nil
}

// Insert a new SKU into the SKU table and return the SKU ID
func InsertSKU(db *sql.DB, sku SKU) int {
	result, err := db.Exec("INSERT INTO SKU (product_name, product_description, product_category) "+
		"VALUES (?, ?, ?)",
		sku.ProductName, sku.ProductDescription, sku.ProductCategory)
	if err != nil {
		log.Fatal(err)
	}
	skuID, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("SKU inserted successfully with ID: %d\n", skuID)
	return int(skuID)
}

// Map coupons to SKUs
func MapCouponsToSKUs(db *sql.DB, coupons []Coupon, skuIDs []int) {
	for _, coupon := range coupons {
		for _, skuID := range skuIDs {
			InsertSKUToCouponMapping(db, coupon.ID, skuID)
		}
	}
	fmt.Println("Coupons mapped to SKUs successfully")
}

// Insert a mapping between a coupon and a SKU
func InsertSKUToCouponMapping(db *sql.DB, couponID, skuID int) {
	_, err := db.Exec("INSERT INTO SKU_Coupon_Mapping (coupon_id, sku_id) VALUES (?, ?)",
		couponID, skuID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SKU-Coupon mapping inserted successfully")
}

// Record coupon usage
func RecordCouponUsage(db *sql.DB, couponID, userID, orderID int) {
	_, err := db.Exec("INSERT INTO CouponUsage (coupon_id, user_id, order_id, usage_date, is_used) "+
		"VALUES (?, ?, ?, NOW(), true)",
		couponID, userID, orderID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Coupon usage recorded successfully for Coupon ID: %d\n", couponID)
}

// Send coupon expiration notifications
func SendCouponExpirationNotifications(db *sql.DB, coupons []Coupon) {
	for _, coupon := range coupons {
		// Check if the coupon is about to expire within a specified threshold (e.g., 7 days)
		expirationDate, _ := time.Parse("2006-01-02", coupon.ExpirationDate)
		expirationThreshold := time.Now().AddDate(0, 0, 7) // 7 days from now

		if expirationDate.Before(expirationThreshold) {
			// Send notification to the user (implementation required)
			userID := 1 // Replace with the actual user ID
			SendExpirationNotification(db, coupon, userID)
		}
	}
}

// Send an expiration notification to the user
func SendExpirationNotification(db *sql.DB, coupon Coupon, userID int) {
	// Implement notification sending logic (e.g., email or push notification)
	// Example:
	fmt.Printf("Sending expiration notification to User ID: %d for Coupon ID: %d\n", userID, coupon.ID)
}

// Define a struct to represent a ruleset
type RuleSet struct {
	Name       string
	Definition string
	Version    string
}

/*
// Example: Apply ruleset to coupons
customerContext := CustomerContext{
    IsSubscriber: true, // Set customer subscription status
    PreviousRedemptions: []RedemptionHistoryEntry{
        {CouponCode: "SUMMER25"},
        // Add other previous redemption history entries as needed
    },
    // Add other customer data
}

changeContext := ChangeContext{
    FromSubscription: Subscription{
        // Define previous subscription details
    },
    ToSubscription: Subscription{
        // Define new subscription details
    },
}

optionsContext := OptionsContext{
    // Retrieve options from the database
    // For example: DiscountPercentage: 25, MinimumPurchaseAmount: 50, etc.
}

applyRuleset(db, ruleset, coupons, customerContext, changeContext, optionsContext)

*/

// Apply a ruleset to coupons for validation
func ApplyRuleset(rulesets []RuleSet, coupons []Coupon, customerContext CustomerContext, changeContext ChangeContext, optionsContext OptionsContext) {
	// Create a new knowledge base for Grule
	knowledgeLibrary := ast.NewKnowledgeLibrary()

	ruleBuilder := builder.NewRuleBuilder(knowledgeLibrary)

	for _, ruleset := range rulesets {
		bs := pkg.NewBytesResource([]byte(ruleset.Definition))
		err := ruleBuilder.BuildRuleFromResource(ruleset.Name, ruleset.Version, bs)
		if err != nil {
			panic(err)
		}

		// Load the ruleset definition into the knowledge base

	}

	for i, coupon := range coupons {

		// Create a knowledge context for each coupon
		ctx := ast.NewDataContext()
		ctx.Add("Coupon", &coupon)
		ctx.Add("CustomerContext", &customerContext)
		ctx.Add("ChangeContext", &changeContext)
		ctx.Add("OptionsContext", &optionsContext)

		// Execute the ruleset with the context

		for _, ruleset := range rulesets {
			knowledgeBase, err := knowledgeLibrary.NewKnowledgeBaseInstance(ruleset.Name, ruleset.Version)
			if err != nil {
				log.Fatal(err)
			}
			engine := engine.NewGruleEngine()
			fmt.Printf("Applying ruleset: %s v%s\n", knowledgeBase.Name, knowledgeBase.Version)
			err = engine.Execute(ctx, knowledgeBase)
			if err != nil {
				log.Printf("Error applying ruleset to coupon %d: %s", i+1, err)
			}
			fmt.Printf("Ruleset %s Version %s applied to coupons\n", knowledgeBase.Name, knowledgeBase.Version)
			// Check if the coupon is valid after ruleset execution
			fmt.Printf("%+v\n", customerContext)
			if coupon.IsValid {
				fmt.Printf("Coupon %s is valid\n", coupon.Code)
			} else {
				fmt.Printf("Coupon %s is not valid: %s \n", coupon.Code, coupon.NotValidReason, coupon)
			}
		}

	}

}

// Implement referral system
func ImplementReferralSystem(db *sql.DB, referrerID, refereeID int) {
	// Check if the referee (new user) made a purchase
	// You would need to implement this logic based on your application flow

	if RefereeMadePurchase(db, refereeID) {
		// Record the referral
		RecordReferral(db, referrerID, refereeID)
	}
}

// Check if the referee (new user) made a purchase
func RefereeMadePurchase(db *sql.DB, refereeID int) bool {
	// Implement logic to check if the referee made a purchase
	// Example: Check the Order table for orders placed by the referee
	// and return true if a purchase is found, otherwise return false
	// ...
	return true // Replace with actual logic
}

// Record the referral
func RecordReferral(db *sql.DB, referrerID, refereeID int) {
	// Record the referral in the Referral table
	_, err := db.Exec("INSERT INTO Referral (referrer_id, referee_id, referral_date, is_rewarded) "+
		"VALUES (?, ?, NOW(), false)",
		referrerID, refereeID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Referral recorded successfully for Referrer ID: %d and Referee ID: %d\n", referrerID, refereeID)
}

// Retrieve coupons associated with a campaign by Campaign ID
func GetCouponsByCampaignID(db *sql.DB, campaignID int) []Coupon {
	rows, err := db.Query("SELECT * FROM Coupons WHERE campaign_id = ?", campaignID)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var coupons []Coupon
	for rows.Next() {
		var coupon Coupon
		if err := rows.Scan(&coupon.ID, &coupon.Code, &coupon.Description, &coupon.DiscountType,
			&coupon.DiscountValue, &coupon.MinimumPurchase, &coupon.ExpirationDate,
			&coupon.IsSingleUse, &coupon.UsageLimit, &coupon.IsActive, &coupon.CampaignID); err != nil {
			log.Fatal(err)
		}
		coupons = append(coupons, coupon)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return coupons
}
