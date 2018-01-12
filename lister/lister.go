package lister

import (
	"flag"
	"path/filepath"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/util/homedir"

	"github.com/hantaowang/kubehandler/utils"
)

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
	}
	return cluster, nil
}

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

func GetNodes(clientset *kubernetes.Clientset) ([]*utils.Node, error) {
	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	cluster := make([]*utils.Node, len(nodes.Items))

	for i, n := range nodes.Items {
		cluster[i] = &utils.Node{
			Name: n.Name,
			HostIP: n.Status.Addresses[0].Address,
			Object: n}
	}
	return cluster, nil
}

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

func MatchPodstoServices(clientset *kubernetes.Clientset, pods []*utils.Pod, sers []*utils.Service) ([]*utils.Pod, []*utils.Service, error) {
	for _, service := range sers {
		podNames, err := GetPodListOfService(clientset, service)
		if err != nil {
			return nil, nil, err
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

func MatchPodstoNodes(pods []*utils.Pod, nodes []*utils.Node) ([]*utils.Pod, []*utils.Node) {
	for _, n := range nodes {
		for _, pod := range pods {
			if pod.HostIP == n.HostIP {
				n.Pods = append(n.Pods, pod)
				pod.Node = n
			}
		}
	}
	return pods, nodes
}

