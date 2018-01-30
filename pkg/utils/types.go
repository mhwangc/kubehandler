package utils

import ("k8s.io/api/core/v1"
)

// Pod, Container Node, Service, and Event objects
type Pod struct {
	Name		string
	Service		*Service
	Node		*Node
	Containers	[]*Container
	HostIP		string
	Namespace	string
	Object		v1.Pod
	Type		string
}

type Container struct {
	Name		string
	Pod			*Pod
	Image		string
	ID			string
}

type Node struct {
	Name		string
	Pods		[]*Pod
	InternIP	string
	ExternIP	string
	Role		string
	Object		v1.Node
	Type		string
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
