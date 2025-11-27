package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/samox73/ksh/pkg/k8s"
	"github.com/samox73/ksh/pkg/tea/views"
)

func main() {
	init := views.BuildNamespaceModel()

	model, err := tea.NewProgram(init, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	// Try to get the final selection from the container model
	result, ok := model.(views.ContainersModel)
	if !ok {
		// User quit before completing selection
		os.Exit(0)
	}

	namespace := result.GetNamespace()
	pod := result.GetPod()
	container := result.GetContainer()

	if namespace != "" && pod != "" && container != "" {
		fmt.Printf("Opening shell to %s/%s/%s\n", namespace, pod, container)
		k8s.OpenShell(result.GetClientset(), namespace, pod, container)
	} else {
		// Selection incomplete
		os.Exit(0)
	}
}
