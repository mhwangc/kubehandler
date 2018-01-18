package triggers

import (
    "fmt"

    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/client-go/kubernetes"
    "github.com/hantaowang/kubehandler/pkg/controller"
    "k8s.io/client-go/util/retry"
)

// ReplicaCount Trigger
var ReplicasWithinCost = controller.Trigger {
    Name: "ReplicasWithinCost",
    Desc: "Trigger that ensures that the total cost of a deployment does not surpass the budgeted amount",
    Satisfied: func(c *controller.Controller, deploymentName string, individualCost int32, totalBudget int32) bool {
        deploymentsClient := c.Clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
        retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
            // RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
            deployment, getErr := deploymentsClient.Get(deploymentName, metav1.GetOptions{})
            return getErr
        })
        if retryErr != nil {
            fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
            return nil
        }
        return (*deployment.Spec.Replicas * individualCost) <= totalBudget
    },
    Enforce: func(c *controller.Controller, deploymentName string, individualCost int32, totalBudget int32) bool {
        deploymentsClient := c.Clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)
        retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
            // RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
            deployment, getErr := deploymentsClient.Get(deploymentName, metav1.GetOptions{})
            if getErr != nil {
                fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
                return getErr
            }
            totalReplicas := *deployment.Spec.Replicas
            totalCost := totalReplicas * individualCost
            if totalCost <= totalBudget {
                return true
            }
            for totalReplicas * individualCost > totalBudget {
                totalReplicas--
            }
            deployment.Spec.Replicas = &totalReplicas
            _, updateErr := deploymentsClient.Update(deployment)
            return updateErr
        })
        if retryErr != nil {
            fmt.Errorf("Failed to enforce ReplicasWithinCost: %v", retryErr)
            return false
        }
        return true
    },
}
