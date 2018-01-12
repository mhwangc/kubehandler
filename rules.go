package main

type Rule struct {
	// Returns if the rule is satisfied
	Satisfied	func(c *Controller) bool

	// Attempts to satisfy the rule, returns if success
	// Also sets c.Lock to false upon completion
	Enforce		func(c *Controller) bool
}


