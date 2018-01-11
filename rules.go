package main

type Rule struct {
	// Returns if the rule is satisfied
	Valid   func(c *Controller) bool

	// Attempts to satisfy the rule, returns if success
	// Also sets c.Lock to false upon completion
	Satisfy func(c *Controller) bool
}


