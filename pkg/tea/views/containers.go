package views

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/samox73/ksh/pkg/k8s"
	"github.com/samox73/ksh/pkg/tea/components"
	"github.com/samox73/ksh/pkg/tea/styles"
	"github.com/samox73/ksh/pkg/tea/utils"
	"k8s.io/client-go/kubernetes"
)

const containerBanner = `
 ██████╗ ██████╗ ███╗   ██╗████████╗ █████╗ ██╗███╗   ██╗███████╗██████╗ 
██╔════╝██╔═══██╗████╗  ██║╚══██╔══╝██╔══██╗██║████╗  ██║██╔════╝██╔══██╗
██║     ██║   ██║██╔██╗ ██║   ██║   ███████║██║██╔██╗ ██║█████╗  ██████╔╝
██║     ██║   ██║██║╚██╗██║   ██║   ██╔══██║██║██║╚██╗██║██╔══╝  ██╔══██╗
╚██████╗╚██████╔╝██║ ╚████║   ██║   ██║  ██║██║██║ ╚████║███████╗██║  ██║
 ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝╚═╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝`

type ContainersModel struct {
	items     list.Model
	namespace string
	pod       string
	container string
	clientset kubernetes.Clientset
	parent    tea.Model
}

func (m ContainersModel) GetContainer() string                { return m.container }
func (m ContainersModel) GetPod() string                      { return m.pod }
func (m ContainersModel) GetNamespace() string                { return m.namespace }
func (m ContainersModel) GetClientset() *kubernetes.Clientset { return &m.clientset }

func (m ContainersModel) Init() tea.Cmd {
	return nil
}

func (m ContainersModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.items.SetWidth(msg.Width)
		m.items.SetHeight(utils.MinInt(msg.Height-lipgloss.Height(containerBanner), len(m.items.Items())))
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "q":
			return m.parent, nil
		case "enter":
			i, ok := m.items.SelectedItem().(components.Item)
			m.container = i.Name
			if ok {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.items, cmd = m.items.Update(msg)
	return m, cmd
}

func (m ContainersModel) View() string {
	banner := lipgloss.NewStyle().Margin(2, 0, 0, 2).Render(styles.GetBanner(containerBanner))
	context := utils.ViewContext()
	items := m.items.View()
	return lipgloss.JoinVertical(lipgloss.Left, banner, context, items)
}

func buildContainerModel(namespace string, pod string, parent tea.Model) *ContainersModel {
	clientset := *k8s.GetKubernetesClientset()
	containers := k8s.GetContainers(clientset, namespace, pod)
	m := &ContainersModel{
		items:     utils.BuildContainerList(containers),
		clientset: clientset,
		namespace: namespace,
		pod:       pod,
		parent:    parent,
	}
	return m
}
