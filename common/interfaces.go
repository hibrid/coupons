package common

// DiscountResult holds information about the application of a discount,
// including original and modified values.
type DiscountResult struct {
	OriginalValues map[string]interface{} // Original input values before the discount was applied
	ModifiedValues map[string]interface{} // Values after the discount application
	Description    string                 // Description of the discount applied
}
