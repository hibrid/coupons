package main

import (
	"database/sql"
	"fmt"
)

type Campaign struct {
	ID                 int
	Name               string
	CouponType         CouponType
	CouponVanityName   string
	CouponDiscountType CouponDiscountType
	DiscountValue      float64

	StartDate string
	EndDate   string
	IsActive  bool
}

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
