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

	result, ok := model.(views.ContainersModel)
	if !ok {
		fmt.Println("resulting model is invalid")
	}
	namespace := result.GetNamespace()
	pod := result.GetPod()
	container := result.GetContainer()
	if pod != "" && namespace != "" && container != "" {
		fmt.Printf("Opening shell to %s/%s/%s", namespace, pod, container)
		k8s.OpenShell(result.GetClientset(), namespace, pod, container)
	} else {
		fmt.Println("invalid values")
	}
}
