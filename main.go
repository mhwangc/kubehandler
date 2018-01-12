package main

import ("net/http"
		"fmt"
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
		go getHandler(w, r)
	} else if r.Method == "POST" {
		go postHandler(w, r)
	}
}

// If GET, then returns to the user the timeline of events seen so far
func getHandler(w http.ResponseWriter, r *http.Request) {
	for _, e := range control.Timeline {
		fmt.Fprintf(w, e.Message + "\n")
	}
}

// If POST, then builds the event and passes it to the controller for processing
func postHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	e := utils.Event{
		Namespace: r.Form["namespace"][0],
		Kind:      r.Form["kind"][0],
		Component: r.Form["component"][0],
		Host:      r.Form["host"][0],
		Reason:    r.Form["reason"][0],
		Status:    r.Form["status"][0],
		Name:      r.Form["name"][0],
		Message:   r.Form["msg"][0],
	}
	control.AddEvent(e)
	fmt.Fprintf(w, "Got POST Reqest")
}



// Runs the controller and starts the server
func main() {
	go control.Run()

	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8000", nil)
}