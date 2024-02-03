package main

import (
	"fmt"
	"log"

	"github.com/hyperjumptech/grule-rule-engine/ast"
	"github.com/hyperjumptech/grule-rule-engine/builder"
	"github.com/hyperjumptech/grule-rule-engine/engine"
	"github.com/hyperjumptech/grule-rule-engine/pkg"
)

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
				fmt.Printf("Coupon %s is not valid: %s, %v \n", coupon.Code, coupon.NotValidReason, coupon)
			}
		}

	}

}
