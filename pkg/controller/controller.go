package controller

import (
	"time"
	"github.com/hantaowang/kubehandler/pkg/utils"
	"github.com/hantaowang/kubehandler/pkg/state"
	"github.com/hantaowang/kubehandler/pkg/rules"
)

// Determines how often the controller checks its rules
const rulePeriod  = time.Millisecond * 500

// The controller object holds the state, timeline, queue of incoming events,
// and list of rules to follow.
type Controller struct {
	FuzzyQueue  chan utils.Event
	Timeline 	[]utils.Event
	Nodes  		map[string]*utils.Node
	Pods		map[string]*utils.Pod
	Services	map[string]*utils.Service
	Lock		bool
	Rules		[]rules.Rule
}

// Adds an event to the queue
func (c *Controller) AddEvent(e utils.Event) {
	c.FuzzyQueue <- e
}

// Adds a rule to follow
func (c *Controller) AddRule(r rules.Rule) {
	c.Rules = append(c.Rules, r)
}

// Every 500ms, checks that the rules are followed
// Every time an object is added to the queue, updates the state
func (c *Controller) Run() {
	c.FuzzyQueue = make(chan utils.Event, 100)
	ticker := time.NewTicker(rulePeriod)

	for {
		select {
			case e := <- c.FuzzyQueue:
				go c.updateEvent(e)
			case <- ticker.C:
				go c.checkRules()
		}
	}
}

// Checks if all rules are satisfied, and attempts to enforce the
// first unsatisfied rule (but only if there is no lock).
func (c *Controller) checkRules() {
	for _, r := range c.Rules {
		if !r.Satisfied(c) {
			if !c.Lock {
				c.Lock = true
				go r.Enforce(c)
			}
		}
	}
}

// Given a controller, updates the Pod, Service, and Host attributes to correctly
// reflect the current cluster state.
func (c *Controller) GetClusterState() {
	client := state.GetClientOutOfCluster()
	pods, err1 := state.GetPods(client)
	nodes, err2 := state.GetNodes(client)
	services, err3 := state.GetServices(client)
	if utils.CheckAllErrors(err1, err2, err3) != nil {
		panic(utils.CheckAllErrors(err1, err2, err3))
	}

	pods, services, err1 = state.MatchPodsToServices(client, pods, services)
	pods, nodes = state.MatchPodsToNodes(pods, nodes)
	if utils.CheckAllErrors(err1) != nil {
		panic(utils.CheckAllErrors(err1))
	}

	for _, n := range nodes {
		c.Nodes[n.Name] = n
	}
	for _, s := range services {
		c.Services[s.Name] = s
	}
	for _, p := range pods {
		c.Pods[p.Name] = p
	}

}

// Parses the event and then updates the state
// Create functions are a little naive...
func (c *Controller) updateEvent(e utils.Event) {
	c.Timeline = append(c.Timeline, e)

	if e.Reason == "delete" {
		if e.Kind == "pod" {
			pod := c.Pods[e.Name]
			pod.Service.Pods = utils.DeletePodNameOnce(pod.Service.Pods, e.Name)
			pod.Node.Pods = utils.DeletePodNameOnce(pod.Node.Pods, e.Name)
			delete(c.Pods, e.Name)
		} else if e.Kind == "service" {
			delete(c.Services, e.Name)
		} else if e.Kind == "machine" {
			delete(c.Nodes, e.Name)
		}
	} else if e.Reason == "create" {
		go c.GetClusterState()
	}

	go c.checkRules()
}