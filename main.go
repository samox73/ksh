package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"
)

const listHeight = 14

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}
	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}
	fmt.Fprint(w, fn(string(i)))
}

type model struct {
	items              list.Model
	namespace          string
	pod                string
	clientset          kubernetes.Clientset
	selectingNamespace bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.items.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.items.SelectedItem().(item)
			if ok {
				if m.selectingNamespace {
					m.namespace = string(i)
					m.items = buildPodList(m.namespace, m.clientset)
					m.selectingNamespace = false
					return m, tea.ClearScreen
				} else {
					m.pod = string(i)
					return m, tea.Quit
				}
			}
		}
	}

	var cmd tea.Cmd
	m.items, cmd = m.items.Update(msg)
	return m, cmd
}

func buildPodList(namespace string, clientset kubernetes.Clientset) list.Model {
	const defaultWidth = 20
	l := list.New(buildPodItems(namespace, clientset), itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select the pod:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func buildNamespaceModel() *model {
	const defaultWidth = 20
	clientset := *getKubernetesClientset()
	items := buildNamespaceItems(clientset)
	l := list.New(items, itemDelegate{}, defaultWidth, listHeight)
	l.Title = "Select the namespace:"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return &model{items: l, clientset: clientset, selectingNamespace: true}
}

func (m *model) View() string {
	if m.pod != "" && m.namespace != "" {
		return quitTextStyle.Render(fmt.Sprintf("Opening a shell into %s/%s", m.namespace, m.pod))
	}
	return "\n" + m.items.View()
}

func buildNamespaceItems(clientset kubernetes.Clientset) []list.Item {
	items := getNamespaces(clientset).Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	out := make([]list.Item, len(items))
	for i, ns := range items {
		out[i] = item(ns.Name)
	}
	return out
}

func buildPodItems(namespace string, clientset kubernetes.Clientset) []list.Item {
	items := getPods(clientset, namespace).Items
	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})
	out := make([]list.Item, len(items))
	for i, pod := range items {
		log.Default().Printf("found pod: %s\n", pod.Name)
		out[i] = item(pod.Name)
	}
	return out
}

func main() {
	m := buildNamespaceModel()

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	if m.namespace != "" && m.pod != "" {
		openShell(&m.clientset, m.namespace, m.pod)
	}
}
