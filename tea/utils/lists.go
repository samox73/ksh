package utils

import (
	"ksh/pkg"
	"ksh/tea/components"
	"ksh/tea/styles"
	"log"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"k8s.io/client-go/kubernetes"
)

func getList(items []list.Item) list.Model {
	length := MinInt(len(items)+7, 20)
	l := list.New(items, components.ItemDelegate{}, 60, length)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = styles.TitleStyle
	l.Styles.PaginationStyle = styles.PaginationStyle
	l.Styles.HelpStyle = styles.HelpStyle
	return l
}

func BuildPodList(namespace string, clientset kubernetes.Clientset) list.Model {
	items := buildPodItems(namespace, clientset)
	return getList(items)
}

func BuildNamespaceList(clientset kubernetes.Clientset) list.Model {
	items := buildNamespaceItems(clientset)
	return getList(items)
}

func buildNamespaceItems(clientset kubernetes.Clientset) []list.Item {
	items := pkg.GetNamespaces(clientset).Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	out := make([]list.Item, len(items))
	for i, ns := range items {
		out[i] = components.Item{Name: ns.Name}
	}
	return out
}

func buildPodItems(namespace string, clientset kubernetes.Clientset) []list.Item {
	items := pkg.GetPods(clientset, namespace).Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	out := make([]list.Item, len(items))
	for i, pod := range items {
		log.Default().Printf("found pod: %s\n", pod.Name)
		out[i] = components.Item{Name: pod.Name, Labels: pod.Labels}
	}
	return out
}
