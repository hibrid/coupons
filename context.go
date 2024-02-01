package main

import "time"

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
	OtherData        map[string]interface{}
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
	OtherData           map[string]interface{}
	// Add other relevant customer data here
}
