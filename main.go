package main

import ("github.com/hantaowang/kubehandler/pkg/utils"
		"github.com/hantaowang/kubehandler/pkg/controller"
		"github.com/hantaowang/kubehandler/pkg/state"
		"github.com/hantaowang/kubehandler/pkg/server"
		"github.com/hantaowang/kubehandler/triggers"
)


// Runs the controller and starts the server
func main() {

	// Initialise a controller
	var control = controller.Controller{
		Nodes: make(map[string]*utils.Node),
		Services: make(map[string]*utils.Service),
		Pods: make(map[string]*utils.Pod),
		Client: state.GetClientOutOfCluster(),
		Triggers: triggers.Triggers,
		FuzzyQueue: make(chan *utils.Event, 100),
		Timeline: make([]*utils.Event, 0),
	}

	go control.Run()

	s := server.Server{Control: &control}
	s.Run()

}