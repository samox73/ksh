package k8s

import (
	"context"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetPods(clientset kubernetes.Clientset, namespaceName string) *corev1.PodList {
	pods, err := clientset.CoreV1().Pods(namespaceName).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %v\n", err)
		os.Exit(1)
	}
	return pods
}

func GetNamespaces(clientset kubernetes.Clientset) *corev1.NamespaceList {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		fmt.Printf("Error listing pods: %v\n", err)
		os.Exit(1)
	}
	return namespaces
}

func GetContainers(clientset kubernetes.Clientset, namespaceName string, podName string) []corev1.Container {
	pod, err := clientset.CoreV1().Pods(namespaceName).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		os.Exit(1)
	}
	return pod.Spec.Containers
}
