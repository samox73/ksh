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

const namespaceBanner = `
███╗   ██╗ █████╗ ███╗   ███╗███████╗███████╗██████╗  █████╗  ██████╗███████╗
████╗  ██║██╔══██╗████╗ ████║██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██╔██╗ ██║███████║██╔████╔██║█████╗  ███████╗██████╔╝███████║██║     █████╗  
██║╚██╗██║██╔══██║██║╚██╔╝██║██╔══╝  ╚════██║██╔═══╝ ██╔══██║██║     ██╔══╝  
██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗███████║██║     ██║  ██║╚██████╗███████╗
╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝`

type namespacesModel struct {
	items     list.Model
	clientset kubernetes.Clientset
	banner    string
}

func (m namespacesModel) Init() tea.Cmd {
	return nil
}

func (m *namespacesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	i, _ := m.items.SelectedItem().(components.Item)
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.items.SetWidth(msg.Width)
		m.items.SetHeight(msg.Height - lipgloss.Height(m.banner) - len(i.Labels) - 2)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.items.SelectedItem().(components.Item)
			if ok {
				return buildPodModel(i.Name, m), tea.ClearScreen
			}
		}
	}

	var cmd tea.Cmd
	m.items, cmd = m.items.Update(msg)
	return m, cmd
}

func (m *namespacesModel) View() string {
	banner := lipgloss.NewStyle().Margin(2, 0, 0, 2).Render(m.banner)
	context := utils.ViewContext()
	items := m.items.View()
	return lipgloss.JoinVertical(lipgloss.Left, banner, context, items)
}

func BuildNamespaceModel() *namespacesModel {
	clientset := *k8s.GetKubernetesClientset()
	namespaces := k8s.GetNamespaces(clientset).Items
	return &namespacesModel{items: utils.BuildNamespaceList(namespaces), clientset: clientset, banner: styles.GetBanner(namespaceBanner)}
}
