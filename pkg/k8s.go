package pkg

import (
	"context"
	"flag"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/util/homedir"
	"k8s.io/kubectl/pkg/cmd/exec"
	"k8s.io/kubectl/pkg/cmd/util/podcmd"
	"k8s.io/kubectl/pkg/scheme"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var clientset *kubernetes.Clientset

func GetKubernetesClientset() *kubernetes.Clientset {
	if clientset != nil {
		return clientset
	}
	// Set up kubeconfig
	home := homedir.HomeDir()
	kubeconfig := flag.String("kubeconfig", fmt.Sprintf("%s/.kube/config", home), "path to the kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("Error building kubeconfig: %v\n", err)
		os.Exit(1)
	}

	// Create Kubernetes client
	c, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("Error creating Kubernetes client: %v\n", err)
		os.Exit(1)
	}
	clientset = c
	return c
}

func OpenShell(clientset *kubernetes.Clientset, namespace, pod string) {
	for _, cmd := range [][]string{{"bash"}, {"ash"}, {"sh"}} {
		if err := openSpecificShell(clientset, namespace, pod, cmd); err != nil {
			fmt.Printf("Error opening shell: %v\n", err)
		} else {
			return
		}
	}
}

func GetPods(clientset kubernetes.Clientset, namespace string) *corev1.PodList {
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
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

func openSpecificShell(clientset *kubernetes.Clientset, namespace, podName string, command []string) error {
	// the following is mostly stolen from https://github.com/kubernetes/kubectl/blob/master/pkg/cmd/exec/exec.go#L305
	config := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	restconfig, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restconfig.GroupVersion = &schema.GroupVersion{}
	restconfig.NegotiatedSerializer = runtime.NewSimpleNegotiatedSerializer(runtime.SerializerInfo{})

	streams := genericclioptions.IOStreams{In: os.Stdin, Out: os.Stdout, ErrOut: os.Stderr}
	p := exec.ExecOptions{
		Executor: &exec.DefaultRemoteExecutor{},
		Config:   restconfig,
		StreamOptions: exec.StreamOptions{
			IOStreams: streams,
			TTY:       true,
		},
	}
	p.Stdin = true

	if len(podName) != 0 {
		p.Pod, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
		if err != nil {
			return err
		}
	}
	pod := p.Pod

	if pod.Status.Phase == corev1.PodSucceeded || pod.Status.Phase == corev1.PodFailed {
		return fmt.Errorf("cannot exec into a container in a completed pod; current phase is %s", pod.Status.Phase)
	}

	containerName := p.ContainerName
	if len(containerName) == 0 {
		container, err := podcmd.FindOrDefaultContainerByName(pod, containerName, p.Quiet, p.ErrOut)
		if err != nil {
			return err
		}
		containerName = container.Name
	}
	t := p.SetupTTY()
	var sizeQueue remotecommand.TerminalSizeQueue
	if t.Raw {
		// this call spawns a goroutine to monitor/update the terminal size
		sizeQueue = t.MonitorSize(t.GetSize())

		// unset p.Err if it was previously set because both stdout and stderr go over p.Out when tty is
		// true
		p.ErrOut = nil
	}

	fn := func() error {
		req := clientset.CoreV1().RESTClient().Post().
			Resource("pods").
			Name(podName).
			Namespace(namespace).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: containerName,
				Command:   command,
				Stdin:     p.Stdin,
				Stdout:    p.Out != nil,
				Stderr:    p.ErrOut != nil,
				TTY:       t.Raw,
			}, scheme.ParameterCodec)

		return p.Executor.Execute(req.URL(), p.Config, p.In, p.Out, p.ErrOut, t.Raw, sizeQueue)
	}
	if err := t.Safe(fn); err != nil {
		return err
	}
	return nil
}
