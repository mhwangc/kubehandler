package rules

import ("github.com/hantaowang/kubehandler/pkg/controller")

type Rule struct {
	// Returns if the rule is satisfied
	Satisfied	func(c *controller.Controller) bool

	// Attempts to satisfy the rule, returns if success
	// Also sets c.Lock to false upon completion
	Enforce		func(c *controller.Controller) bool
}


