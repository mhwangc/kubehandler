package triggers

import (
    "github.com/hantaowang/kubehandler/pkg/controller"
    "github.com/hantaowang/kubehandler/pkg/kubefunc"
    "fmt"
)

const minReplicas int = 2

var MinReplicasAll = controller.Trigger{
    Name: "MinReplicasAll",
    Desc: fmt.Sprintf("Each service must have at least %v replicas", minReplicas),
    Satisfied: func(c *controller.Controller) bool {
        for _, s := range c.Services {
            if s.Name != "kube-dns" && s.Name != "kubernetes" && len(s.Pods) < minReplicas {
                return false
            }
        }
        return true
    },
    Enforce: func(c *controller.Controller) error {
        for _, s := range c.Services {
            if s.Name != "kube-dns" && s.Name != "kubernetes" && len(s.Pods) < minReplicas {
                errRep := kubefunc.ReplicaUpdate(c.Client, s.Name, int32(minReplicas - len(s.Pods)))
                if errRep != nil {
                    return errRep
                }
            }
        }
        return nil
    },
}
