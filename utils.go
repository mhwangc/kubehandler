package main

import (
	"fmt"
	"strings"
)

// Pod, Machine, Service, and Event objects

type Pod struct {
	Name	string
	Service	Service
}

type Machine struct {
	Name	string
	Pods	[]Pod
}

type Service struct {
	Name	string
}

type Event struct {
	Namespace string
	Kind      string
	Component string
	Host      string
	Reason    string
	Status    string
	Name      string
	Message	  string
}

// Legacy sort function for when the FuzzyQueue was actually fuzzy
// Might be used if the replicationcontroller cannot be turned off
func sort(passed []Event) []Event {
	delete := make([]Event, 0)
	create := make([]Event, 0)
	update := make([]Event, 0)
	for _, e := range passed {
		if e.Reason == "deleted" {
			delete = append(delete, e)
		} else if e.Reason == "created" {
			create = append(create, e)
		} else if e.Reason == "updated" {
			update = append(update, e)
		} else {
			fmt.Printf("Unexpected event : %s\n", e)
		}
	}
	combined := append(delete, create...)
	combined = append(combined, update...)
	return combined
}

// Given a pod name, determines which service it is a part of
// This is currently a HACK (and probably doesn't work)
func getServiceName(podName string) string {
	broken := strings.Split(podName, "-")
	broken = broken[:len(broken)-2]
	return strings.Join(broken, "-")
}