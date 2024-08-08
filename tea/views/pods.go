package views

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samox73/ksh/pkg"
	"github.com/samox73/ksh/tea/components"
	"github.com/samox73/ksh/tea/styles"
	"github.com/samox73/ksh/tea/utils"
	"k8s.io/client-go/kubernetes"
)

const banner = `
██████╗  ██████╗ ██████╗ 
██╔══██╗██╔═══██╗██╔══██╗
██████╔╝██║   ██║██║  ██║
██╔═══╝ ██║   ██║██║  ██║
██║     ╚██████╔╝██████╔╝
╚═╝      ╚═════╝ ╚═════╝`

type PodsModel struct {
	items     list.Model
	namespace string
	pod       string
	clientset kubernetes.Clientset
	parent    tea.Model
}

func (m PodsModel) GetPod() string                      { return m.pod }
func (m PodsModel) GetNamespace() string                { return m.namespace }
func (m PodsModel) GetClientset() *kubernetes.Clientset { return &m.clientset }

func (m PodsModel) Init() tea.Cmd {
	return nil
}

func (m PodsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	i, _ := m.items.SelectedItem().(components.Item)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.items.SetWidth(msg.Width)
		m.items.SetHeight(utils.MinInt(msg.Height - lipgloss.Height(banner) - len(i.Labels), len(m.items.Items())))
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":
			return m.parent, nil
		case "enter":
			i, ok := m.items.SelectedItem().(components.Item)
			m.pod = i.Name
			if ok {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.items, cmd = m.items.Update(msg)
	return m, cmd
}


func (m *PodsModel) viewLabels() string {
	i, ok := m.items.SelectedItem().(components.Item)
	if !ok {
		return ""
	}

	l := ""
	keys := make([]string, 0, len(i.Labels))
	longestKeyLength := 0
	for k := range i.Labels {
		keys = append(keys, k)
		if len(k) > longestKeyLength {
			longestKeyLength = len(k)
		}
	}
	sort.Strings(keys)
	for j, k := range keys {
		l += fmt.Sprintf("%*s: %s", longestKeyLength, k, i.Labels[k])
		if j != len(i.Labels)-1 {
			l += "\n"
		}
	}
	return lipgloss.NewStyle().Margin(0, 0, 0, 2).Border(lipgloss.NormalBorder(), true).Render(l)
}

func (m PodsModel) View() string {
	banner := lipgloss.NewStyle().Margin(2, 0, 0, 2).Render(styles.GetBanner(banner))
	context := utils.ViewContext()
	labels := m.viewLabels()
	items := m.items.View()
	return lipgloss.JoinVertical(lipgloss.Left, banner, context, labels, items)
}

func buildPodModel(namespace string, parent tea.Model) *PodsModel {
	clientset := *pkg.GetKubernetesClientset()
	pods := pkg.GetPods(clientset, namespace).Items
	m := &PodsModel{
		items:     utils.BuildPodList(pods),
		clientset: clientset,
		namespace: namespace,
		parent:    parent,
	}
	return m
}
