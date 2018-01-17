package rules

import (
	"k8s.io/client-go/kubernetes"
	"github.com/hantaowang/kubehandler/pkg/controller"
)

// Example rule
var NoMoreThanThreeMachines = controller.Rule{
	Satisfied: func(c *controller.Controller) bool {
		return len(c.Nodes) <= 3
	},
	Enforce: func(c *controller.Controller) bool {
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
