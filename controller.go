package main

import (
	"time"
	"github.com/hantaowang/kubehandler/utils"
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
	Rules		[]Rule
}

// Adds an event to the queue
func (c *Controller) AddEvent(e utils.Event) {
	c.FuzzyQueue <- e
}

// Adds a rule to follow
func (c *Controller) AddRule(r Rule) {
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

// Parses the event and then updates the state
func (c *Controller) updateEvent(e utils.Event) {
	c.Timeline = append(c.Timeline, e)

	if e.Reason == "delete" {
		if e.Kind == "pod" {
			delete(c.Pods, e.Name)
		} else if e.Kind == "service" {
			delete(c.Services, e.Name)
		} else if e.Kind == "machine" {
			delete(c.Nodes, e.Name)
		}
	} else if e.Reason == "create" {
		if e.Kind == "pod" {
			p := utils.Pod{Name: e.Name}
			c.Pods[p.Name] = &p
		} else if e.Kind == "service" {
			s := utils.Service{Name: e.Name}
			c.Services[s.Name] = &s
		} else if e.Kind == "machine" {
			m := utils.Node{Name: e.Name}
			c.Nodes[m.Name] = &m
		}
	}
}