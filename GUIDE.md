## GUIDE

This is a guide on how Kubehandler is structured. 
`Controller` is where everything happens in Kubewatch. 
It is responsible for keeping the state of the cluster and also
enforcing triggers.

Updates are received by the server in `pkg/server/server.go`. These updates are POST
requests send by `kubewatch` at https://github.com/hantaowang/kubewatch. The server
reads the form of the update and builds an `utils.Event` object based on that data.
This can either be a delete, create, or update event. The event is passed into a channel that
is part of the controller, called the `FuzzyQueue`.

Triggers follow a specific format as defined in `pkg/controller/trigger.go`. Each trigger has a
name, description, `Satisfied` function, and `Enforce` function. The `Satisfied` function is meant to be a 
lightweight function on the controller's state that verfies if the trigger needs to be triggered. The
`Enforce` function is a function that uses tools from `pkg/controller/kubefunc` to operate on the
cluster and attempt to satisfy the trigger. For example, it can change the number of replicas.

When a new event is added to the `FuzzyQueue`, the controller
will remove it from the queue and parse that event. Deletion events removes a specific `utils.Pod, utils.Service, utils.Node` object from
the state. Creation and update events make a request to update the cluster state by setting `UpdateRequest` to 1.
All events make a request to check triggers by setting `TriggerRequest` to 1. These are done atomically.

The controller operates on a switch between 3 cases:
1) There is an item in `FuzzyQueue` -> remove and parse that event, 
update state if required, set `TriggerRequest` or `UpdateRequest` to 1 if required.
2) Every `triggerPeriod` (default 30sec), atomically check if `TriggerRequest` is set to 1. If so,
call `Satisfied()` on every trigger and `Enforce()` on the first unsafisfied trigger. `atmoic.CompareAndSwap` is used on `Lock` to make sure
no two triggers are every enforced at the same time. Then set `TriggerRequest` to 0.
3) Default: Atmoically check if `UpdateRequest` is set to 1. If so, update the cluster state through `client-go` and match
all Pods with Services and Pods with Nodes. Set `UpdateRequest` back to 0.