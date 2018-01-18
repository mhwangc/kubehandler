package triggers

import "github.com/hantaowang/kubehandler/pkg/controller"

// Initializes list of Rules to be enforced
var Triggers = []controller.Trigger{
	NoMoreThanThreeMachines,
}
