package kubefunc

import (
    "fmt"
    "strconv"

    apiv1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/util/retry"
)

// Updates the replica count of Deployment METANAME by the value QUANTITY
func ReplicaUpdate(clientset *kubernetes.Clientset, metaname string, quantity string) error {
    // Getting deployments
    deploymentsClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceDefault)

    // Updating deployment
    retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
        // Retrieve the latest version of Deployment before attempting update
        // RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
        result, getErr := deploymentsClient.Get(metaname, metav1.GetOptions{})
        if getErr != nil {
            fmt.Errorf("Failed to get latest version of Deployment: %v", getErr)
            return getErr
        }

        fmt.Printf("Updating replica count of %v by %v\n", metaname, quantity)

        // Parsing quantity to int32
        i, err := strconv.ParseInt(quantity, 10, 32)
        if err != nil {
            fmt.Errorf("Failed to parse int: %v", quantity)
            return err
        }

        // Modify replica count
        oldRep := result.Spec.Replicas
        result.Spec.Replicas = int32Ptr(*oldRep + int32(i))
        if *result.Spec.Replicas < int32(1) {
            result.Spec.Replicas = int32Ptr(1)
        }
        _, updateErr := deploymentsClient.Update(result)
        return updateErr
    })
    if retryErr != nil {
        fmt.Errorf("Update failed: %v", retryErr)
        return retryErr
    }
    fmt.Printf("Updated replica count of Deployment %v\n", metaname)
    return nil
}

func int32Ptr(i int32) *int32 { return &i }
