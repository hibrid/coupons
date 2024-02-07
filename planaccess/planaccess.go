package planaccess

type Plan struct {
	ID   string
	Name string
}

type PlanAccessDiscountConfig struct {
	Plans []Plan // Identifier for the plan being accessed
}
