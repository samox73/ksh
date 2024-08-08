package utils

import (
	"sort"

	"github.com/charmbracelet/bubbles/list"
	"github.com/samox73/ksh/pkg/tea/components"
	"github.com/samox73/ksh/pkg/tea/styles"
	corev1 "k8s.io/api/core/v1"
)

func listFromItems(items []list.Item) list.Model {
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

func BuildNamespaceList(namespaces []corev1.Namespace) list.Model {
	items := buildNamespaceItems(namespaces)
	return listFromItems(items)
}

func buildNamespaceItems(namespaces []corev1.Namespace) []list.Item {
	sort.Slice(namespaces, func(i, j int) bool {
		return namespaces[i].Name < namespaces[j].Name
	})
	out := make([]list.Item, len(namespaces))
	for i, ns := range namespaces {
		out[i] = components.Item{Name: ns.Name}
	}
	return out
}

func BuildPodList(pods []corev1.Pod) list.Model {
	items := buildPodItems(pods)
	return listFromItems(items)
}

func buildPodItems(pods []corev1.Pod) []list.Item {
	sort.Slice(pods, func(i, j int) bool {
		return pods[i].Name < pods[j].Name
	})
	out := make([]list.Item, len(pods))
	for i, pod := range pods {
		out[i] = components.Item{Name: pod.Name, Labels: pod.Labels}
	}
	return out
}

func BuildContainerList(pods []corev1.Container) list.Model {
	items := buildContainerItems(pods)
	return listFromItems(items)
}

func buildContainerItems(pods []corev1.Container) []list.Item {
	sort.Slice(pods, func(i, j int) bool {
		return pods[i].Name < pods[j].Name
	})
	out := make([]list.Item, len(pods))
	for i, pod := range pods {
		out[i] = components.Item{Name: pod.Name}
	}
	return out
}
