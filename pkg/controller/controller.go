package controller

import (
	"time"
	"github.com/hantaowang/kubehandler/pkg/utils"
	"github.com/hantaowang/kubehandler/pkg/state"
	"k8s.io/client-go/kubernetes"
	"fmt"
	"sync/atomic"
)

// The controller object holds the state, timeline, queue of incoming events,
// and list of triggers to follow.
type Controller struct {
	FuzzyQueue  	chan *utils.Event
	UpdateRequest	int32
	Timeline 		[]*utils.Event
	Nodes  			map[string]*utils.Node
	Pods			map[string]*utils.Pod
	Services		map[string]*utils.Service
	Lock			int32
	Triggers		[]Trigger
	Client			*kubernetes.Clientset
}

// Determines how often the controller checks its triggers
const triggerPeriod  = time.Second * 5

// Adds an event to the queue
func (c *Controller) AddEvent(e *utils.Event) {
	c.FuzzyQueue <- e
}

// Adds a trigger to follow
func (c *Controller) AddTrigger(r Trigger) {
	c.Triggers = append(c.Triggers, r)
}

// Every 500ms, checks that the triggers are followed
// Every time an object is added to the queue, updates the state
func (c *Controller) Run() {
	ticker := time.NewTicker(triggerPeriod)

	for {
		select {
			case e := <- c.FuzzyQueue:
				c.Timeline = append(c.Timeline, e)
				go c.updateEvent(e)
			case <- ticker.C:
				go c.checkTriggers()
			default:
				if atomic.LoadInt32(&c.UpdateRequest) == 1 {
					// Blocking
					c.GetClusterState()
				}
		}
	}
}

// Checks if all triggers are satisfied, and attempts to enforce the
// first unsatisfied trigger (but only if there is no lock).
func (c *Controller) checkTriggers() {
	for _, r := range c.Triggers {
		if !r.Satisfied(c) {
			if atomic.LoadInt32(&c.Lock) == 0 {
				atomic.StoreInt32(&c.Lock, 1)
				go r.Enforce(c)
			}
		}
	}
}

// Given a controller, updates the Pod, Service, and Host attributes to correctly
// reflect the current cluster state.
func (c *Controller) GetClusterState() {
	pods, err1 := state.GetPods(c.Client)
	nodes, err2 := state.GetNodes(c.Client)
	services, err3 := state.GetServices(c.Client)
	if utils.CheckAllErrors(err1, err2, err3) != nil {
		panic(utils.CheckAllErrors(err1, err2, err3))
	}

	pods, services, err1 = state.MatchPodsToServices(c.Client, pods, services)
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

	atomic.StoreInt32(&c.UpdateRequest, 0)
	fmt.Printf("[%s] Cluster State Updated\n", utils.GetTimeString())

}

// Parses the event and then updates the state
// Create, update functions are a little naive...
func (c *Controller) updateEvent(e *utils.Event) {

	if e.Reason == "deleted" {
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
	} else if e.Reason == "created" {
		atomic.StoreInt32(&c.UpdateRequest, 1)
		if e.Kind == "machine" {
			e.Message = "A new machine has been created: " + e.Name
		}
	} else if e.Reason == "updated" {
		atomic.StoreInt32(&c.UpdateRequest, 1)
	}

	c.checkTriggers()
}

