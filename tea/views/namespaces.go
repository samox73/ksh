package views

import (
	"ksh/pkg"
	"ksh/tea/components"
	"ksh/tea/styles"
	"ksh/tea/utils"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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
	items := m.items.View()
	return lipgloss.JoinVertical(lipgloss.Left, banner, items)
}

func BuildNamespaceModel() *namespacesModel {
	clientset := *pkg.GetKubernetesClientset()
	return &namespacesModel{items: utils.BuildNamespaceList(clientset), clientset: clientset, banner: styles.GetBanner(namespaceBanner)}
}
