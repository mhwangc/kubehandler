package server

import (
	"github.com/hantaowang/kubehandler/pkg/utils"
	"fmt"
	"net/http"
	"html/template"
	"github.com/hantaowang/kubehandler/pkg/controller"
)

type Server struct {
	Control	*controller.Controller
}

type WebPage struct {
	Pods 		[]*utils.Pod
	Services	[]*utils.Service
	Nodes		[]*utils.Node
	Timeline	[]*utils.Event
	Triggers	[]controller.Trigger
}

// Handles requests to :8000 and redirects pased on POST or GET
func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		s.getHandler(w, r)
	} else if r.Method == "POST" {
		s.postHandler(w, r)
	}
}

// If GET, then returns to the user the timeline of events seen so far
func (s *Server) getHandler(w http.ResponseWriter, r *http.Request) {

	pods := make([]*utils.Pod, 0, len(s.Control.Pods))
	services := make([]*utils.Service, 0, len(s.Control.Services))
	nodes := make([]*utils.Node, 0, len(s.Control.Nodes))
	triggers := make([]controller.Trigger, 0, len(s.Control.Triggers))

	for _, p := range s.Control.Pods {pods = append(pods, p)}
	for _, s := range s.Control.Services {services = append(services, s)}
	for _, n := range s.Control.Nodes {nodes = append(nodes, n)}
	for _, t := range s.Control.Triggers {triggers = append(triggers, t) }


	podsSorted := make([]*utils.Pod, 0, len(s.Control.Pods))
	for _, s := range services {
		for _, p := range pods {
			if p.Service != nil && s.Name == p.Service.Name {
				podsSorted = append(podsSorted, p)
			}
		}
	}
	for _, p := range pods {
		if p.Service == nil {
			podsSorted = append(podsSorted, p)
		}
	}

	page := WebPage{
		Timeline: s.Control.Timeline,
		Pods: podsSorted,
		Services: services,
		Nodes: nodes,
		Triggers: triggers,
	}

	tmpl := template.Must(template.ParseFiles("pkg/server/index.html"))
	tmpl.Execute(w, page)

	fmt.Printf("[%s] Got GET Reqest\n", utils.GetTimeString())

}

// If POST, then builds the event and passes it to the controller for processing
func (s *Server) postHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	e := utils.Event{
		Namespace: r.Form.Get("namespace"),
		Kind:      r.Form["kind"][0],
		Component: r.Form["component"][0],
		Host:      r.Form["host"][0],
		Reason:    r.Form["reason"][0],
		Status:    r.Form["status"][0],
		Name:      r.Form["name"][0],
		Message:   r.Form["event"][0],
		Time:      utils.GetTimeString(),
	}
	go s.Control.AddEvent(&e)
	fmt.Printf("[%s] Got POST Reqest\n", utils.GetTimeString())

}

// Returns pong on ping
func (s *Server) pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] Got PING Reqest\n", utils.GetTimeString())
	fmt.Fprint(w, "pong")
}

// Runs the server
func (s *Server) Run() {
	fmt.Println("Running server on localhost:9000")
	http.HandleFunc("/", s.indexHandler)
	http.HandleFunc("/ping", s.pingHandler)

	http.ListenAndServe(":9000", nil)
}