package controller

type Rule struct {
	// Returns if the rule is satisfied
	Satisfied	func(c *Controller) bool

	// Attempts to satisfy the rule, returns if success
	// Also atomically sets c.Lock to 0 upon completion
	Enforce		func(c *Controller) bool
}
