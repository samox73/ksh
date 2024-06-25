package utils

import (
	"fmt"
	"path/filepath"

	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func ViewContext() string {
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config := clientcmd.GetConfigFromFileOrDie(kubeconfig)
	if config != nil {
		l := fmt.Sprintf("context: %s", config.CurrentContext)
		return lipgloss.NewStyle().Margin(0, 0, 0, 2).Render(l)
	} else {
		return ""
	}
}