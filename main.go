package main

import ("net/http"
		"fmt"
		"html/template"
		"github.com/hantaowang/kubehandler/pkg/utils"
		"github.com/hantaowang/kubehandler/pkg/controller"
		"github.com/hantaowang/kubehandler/pkg/state"
)

// Initialise a controller
var control = controller.Controller{
	Nodes: make(map[string]*utils.Node),
	Services: make(map[string]*utils.Service),
	Pods: make(map[string]*utils.Pod),
	Client: state.GetClientOutOfCluster(),
}

// Handles requests to :8000 and redirects pased on POST or GET
func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		getHandler(w, r)
	} else if r.Method == "POST" {
		postHandler(w, r)
	}
}

// If GET, then returns to the user the timeline of events seen so far
func getHandler(w http.ResponseWriter, r *http.Request) {

	pods := make([]*utils.Pod, 0, len(control.Pods))
	services := make([]*utils.Service, 0, len(control.Services))
	nodes := make([]*utils.Node, 0, len(control.Nodes))

	for _, p := range control.Pods {pods = append(pods, p)}
	for _, s := range control.Services {services = append(services, s)}
	for _, n := range control.Nodes {nodes = append(nodes, n)}


	podsSorted := make([]*utils.Pod, 0, len(control.Pods))
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

	page := utils.WebPage{
		Timeline: control.Timeline,
		Pods: podsSorted,
		Services: services,
		Nodes: nodes,
	}

	tmpl := template.Must(template.ParseFiles("index.html"))
	tmpl.Execute(w, page)

	fmt.Printf("[%s] Got GET Reqest\n", utils.GetTimeString())

}

// If POST, then builds the event and passes it to the controller for processing
func postHandler(w http.ResponseWriter, r *http.Request) {
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
	go control.AddEvent(e)
	fmt.Printf("[%s] Got POST Reqest\n", utils.GetTimeString())

}

// Returns pong on ping
func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("[%s] Got PING Reqest\n", utils.GetTimeString())
	fmt.Fprint(w, "pong")
}

// Runs the controller and starts the server
func main() {
	go control.Run()

	fmt.Println("Running server on localhost:9000")
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/ping", pingHandler)

	http.ListenAndServe(":9000", nil)
}