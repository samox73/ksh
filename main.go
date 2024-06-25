package main

import (
	"fmt"
	"ksh/pkg"
	"ksh/tea/views"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	m := views.BuildNamespaceModel()

	model, err := tea.NewProgram(m, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	pods, ok := model.(views.PodsModel)
	if ok {
		namespace := pods.GetNamespace()
		pod := pods.GetPod()
		if pod != "" && namespace != "" {
			fmt.Printf("Opening shell to %s/%s", namespace, pod)
			pkg.OpenShell(pods.GetClientset(), namespace, pod)
		}
	}
}
