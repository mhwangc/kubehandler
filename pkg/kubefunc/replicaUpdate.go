package kubefunc

import (
    "fmt"

    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/util/retry"
    "github.com/hantaowang/kubehandler/pkg/utils"
)

// Updates the replica count of Deployment METANAME by the value QUANTITY
func ReplicaUpdate(clientset *kubernetes.Clientset, metaname string, quantity int32) error {
    // Getting deployments
    deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

    // Updating deployment
    retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
        // Retrieve the latest version of Deployment before attempting update
        // RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
        result, getErr := deploymentsClient.Get(metaname, metav1.GetOptions{})
        if getErr != nil {
            fmt.Errorf("[%s] Failed to get latest version of Deployment: %s", utils.GetTimeString(), getErr)
            return getErr
        }
        fmt.Printf("[%s] Updating replica count of %v by %d\n", utils.GetTimeString(), metaname, quantity)

        // Modify replica count
        oldRep := result.Spec.Replicas
        result.Spec.Replicas = int32Ptr(*oldRep + int32(quantity))
        if *result.Spec.Replicas < int32(1) {
            result.Spec.Replicas = int32Ptr(1)
        }
        _, updateErr := deploymentsClient.Update(result)
        return updateErr
    })
    if retryErr != nil {
        fmt.Errorf("[%s] Update failed: %s\n", utils.GetTimeString(), retryErr)
        return retryErr
    }
    fmt.Printf("[%s] Updated replica count of Deployment %v\n", utils.GetTimeString(), metaname)
    return nil
}

func int32Ptr(i int32) *int32 { return &i }
