package controller

type Trigger struct {
	Name		string
	Desc		string

	// Returns if the trigger is satisfied
	Satisfied	func(c *Controller) bool

	// Attempts to satisfy the trigger, returns any error
	// Also atomically sets c.Lock to 0 upon completion
	Enforce		func(c *Controller) error
}
