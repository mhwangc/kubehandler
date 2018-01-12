package controller

import "k8s.io/client-go/kubernetes"

// Example rule
var NoMoreThanThreeMachines = Rule{
	Satisfied: func(c *Controller) bool {
		return len(c.Nodes) <= 3
	},
	Enforce: func(c *Controller) bool {
		err := deleteRandomMachine(c.Client)
		c.Lock = false
		if err != nil {
			return false
		}
		return true
	},
}

func deleteRandomMachine(client *kubernetes.Clientset) error {
	// Not implemented
	return nil
}
