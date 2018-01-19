package triggers

import (
	"k8s.io/client-go/kubernetes"
	"github.com/hantaowang/kubehandler/pkg/controller"
)

// Example trigger
var NoMoreThanThreeMachines = controller.Trigger{
	Name: "NoMoreThanThreeMachines",
	Desc: "Trigger that there cannot be more than 3 machines (workers and master) active at any given time",
	Satisfied: func(c *controller.Controller) bool {
		return len(c.Nodes) <= 3
	},
	Enforce: func(c *controller.Controller) error {
		err := deleteRandomMachine(c.Client)
		if err != nil {
			return err
		}
		return nil
	},
}

func deleteRandomMachine(client *kubernetes.Clientset) error {
	// Not implemented
	return nil
}
