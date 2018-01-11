package main

import (
	"time"
)

// Determines how often the controller checks its rules
const rulePeriod  = time.Millisecond * 500

// The controller object holds the state, timeline, queue of incoming events,
// and list of rules to follow.
type Controller struct {
	FuzzyQueue  chan Event
	Timeline 	[]Event
	Machines  	map[string]Machine
	Pods		map[string]Pod
	Services	map[string]Service
	Lock		bool
	Rules		[]Rule
}

// Adds an event to the queue
func (c *Controller) AddEvent(e Event) {
	c.FuzzyQueue <- e
}

// Adds a rule to follow
func (c *Controller) AddRule(r Rule) {
	c.Rules = append(c.Rules, r)
}

// Every 500ms, checks that the rules are followed
// Every time an object is added to the queue, updates the state
func (c *Controller) Run() {
	c.FuzzyQueue = make(chan Event, 100)
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

// Checks if all rules are satisfied, and attempts to satisfy the
// first unsatisfied rule (but only if there is no lock).
func (c *Controller) checkRules() {
	for _, r := range c.Rules {
		if !r.Valid(c) {
			if !c.Lock {
				c.Lock = true
				go r.Satisfy(c)
			}
		}
	}
}

// Parses the event and then updates the state
func (c *Controller) updateEvent(e Event) {
	c.Timeline = append(c.Timeline, e)

	if e.Reason == "delete" {
		if e.Kind == "pod" {
			delete(c.Pods, e.Name)
		} else if e.Kind == "service" {
			delete(c.Services, e.Name)
		} else if e.Kind == "machine" {
			delete(c.Machines, e.Name)
		}
	} else if e.Reason == "create" {
		if e.Kind == "pod" {
			p := Pod{Name: e.Name}
			c.Pods[p.Name] = p
		} else if e.Kind == "service" {
			s := Service{Name: e.Name}
			c.Services[s.Name] = s
		} else if e.Kind == "machine" {
			m := Machine{Name: e.Name}
			c.Machines[m.Name] = m
		}
	}
}