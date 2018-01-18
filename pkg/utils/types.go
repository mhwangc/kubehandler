package utils

import ("k8s.io/api/core/v1"
)

// Pod, Node, Service, and Event objects
type Pod struct {
	Name		string
	Service		*Service
	Node		*Node
	HostIP		string
	Namespace	string
	Object		v1.Pod
}

type Node struct {
	Name		string
	Pods		[]*Pod
	HostIP		string
	Role		string
	Object		v1.Node
}

type Service struct {
	Name		string
	Pods		[]*Pod
	Namespace	string
	Object		v1.Service
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
	Time	  string
}