package state

import (
	"flag"
	"path/filepath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/util/homedir"

	"github.com/hantaowang/kubehandler/pkg/utils"
	"fmt"
)

// Creates a kubernetes out of cluster client with client-go
func GetClientOutOfCluster() *kubernetes.Clientset{
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	return clientset
}

// Queries the client and retrieves all v1.Pod{} objects in the cluster.
// Parses them into utils.Pod objects.
func GetPods(clientset *kubernetes.Clientset) ([]*utils.Pod, error) {
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	cluster := make([]*utils.Pod, len(pods.Items))
	for i, pod := range pods.Items {
		cluster[i] = &utils.Pod{
			Name: pod.Name,
			Namespace: pod.Namespace,
			HostIP: pod.Status.HostIP,
			Object: pod}
		containers := make([]*utils.Container, len(pod.Spec.Containers))
		for j, c := range pod.Status.ContainerStatuses {
			containers[j] = &utils.Container{
				Name:  c.Name,
				Image: c.Image,
				ID:    c.ContainerID,
				Pod:   cluster[i],
			}
			fmt.Println(c.ContainerID)
		}
		cluster[i].Containers = containers
	}
	return cluster, nil
}

// Queries the client and retrieves all v1.Service{} objects in the cluster.
// Parses them into utils.Service objects.
func GetServices(clientset *kubernetes.Clientset) ([]*utils.Service, error) {
	services, err := clientset.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	cluster := make([]*utils.Service, len(services.Items))

	for i, ser := range services.Items {
		cluster[i] = &utils.Service{
			Name: ser.Name,
			Namespace: ser.Namespace,
			Object: ser}
	}
	return cluster, nil
}

// Queries the client and retrieves all v1.Node{} objects in the cluster.
// Parses them into utils.Node objects.
func GetNodes(clientset *kubernetes.Clientset) ([]*utils.Node, error) {
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	cluster := make([]*utils.Node, len(nodes.Items))

	for i, n := range nodes.Items {
		cluster[i] = &utils.Node{
			Name: n.Name,
			InternIP: n.Status.Addresses[0].Address,
			ExternIP: n.Status.Addresses[1].Address,
			Object: n}
	}
	return cluster, nil
}

// Queries the client and finds all pods attached to the given service.
// Returns a list of the pod names.
func GetPodListOfService(clientset *kubernetes.Clientset, ser *utils.Service) ([]string, error) {
	podNames := make([]string, 0)
	if ser.Namespace == "default" && ser.Name == "kubernetes" {
		return podNames, nil
	}
	set := labels.Set(ser.Object.Spec.Selector)

	if pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{LabelSelector: set.AsSelector().String()}); err != nil {
		return nil, err
	} else {
		for _, v := range pods.Items {
			podNames = append(podNames, v.Name)
		}
	}

	return podNames, nil
}

// Given a list of pods and a list of services, matches which pods below to which services.
// Modifies the pod and service objects to reflect this information.
func MatchPodsToServices(clientset *kubernetes.Clientset, pods []*utils.Pod, sers []*utils.Service) ([]*utils.Pod, []*utils.Service, error) {
	for _, service := range sers {
		podNames, err := GetPodListOfService(clientset, service)
		if err != nil {
			return pods, sers, err
		}
		for _, name := range podNames {
			for _, pod := range pods {
				if name == pod.Name {
					service.Pods = append(service.Pods, pod)
					pod.Service = service
				}
			}
		}
	}
	return pods, sers, nil
}

// Given a list of pods and a list of nodes, matches which pods are hosted on which nodes.
// Modifies the pod and node objects to reflect this information.
func MatchPodsToNodes(pods []*utils.Pod, nodes []*utils.Node) ([]*utils.Pod, []*utils.Node) {
	for _, n := range nodes {
		n.Role = "worker"
		for _, pod := range pods {
			if pod.HostIP == n.InternIP {
				n.Pods = append(n.Pods, pod)
				pod.Node = n
				if len(pod.Name) >= 14 && pod.Name[0:14] == "kube-apiserver" {
					n.Role = "master"
				}
			}
		}
	}
	return pods, nodes
}

