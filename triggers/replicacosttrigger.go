package triggers

import (
    "github.com/hantaowang/kubehandler/pkg/controller"
    "github.com/hantaowang/kubehandler/pkg/kubefunc"
    "sync/atomic"
    "fmt"
    "math"
)

const costPerPod float64 = 5.0
const maxCostPerService float64 = 8.0

var ReplicasWithinCostAll = controller.Trigger{
    Name: "ReplicasWithinCost",
    Desc: fmt.Sprintf("Each service cannot cost more than $%d", maxCostPerService),
    Satisfied: func(c *controller.Controller) bool {
        for _, s := range c.Services {
            if float64(len(s.Pods)) * costPerPod > maxCostPerService {
                return false
            }
        }
        return true
    },
    Enforce: func(c *controller.Controller) error {
        for _, s := range c.Services {
            if float64(len(s.Pods)) * costPerPod > maxCostPerService {
                dif := maxCostPerService - float64(len(s.Pods)) * costPerPod
                delPod := "-" + string(int(math.Ceil(dif / costPerPod)))
                errRep := kubefunc.ReplicaUpdate(c.Client, s.Name, delPod)
                if errRep != nil {
                    return errRep
                }
            }
        }
        return nil
    },
}
