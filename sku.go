package main

import (
	"database/sql"
	"fmt"
	"log"
)

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
