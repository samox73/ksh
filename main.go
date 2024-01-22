package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"k8s.io/client-go/kubernetes"
)

const listHeight = 10

const namespaceBanner = `
███╗   ██╗ █████╗ ███╗   ███╗███████╗███████╗██████╗  █████╗  ██████╗███████╗
████╗  ██║██╔══██╗████╗ ████║██╔════╝██╔════╝██╔══██╗██╔══██╗██╔════╝██╔════╝
██╔██╗ ██║███████║██╔████╔██║█████╗  ███████╗██████╔╝███████║██║     █████╗  
██║╚██╗██║██╔══██║██║╚██╔╝██║██╔══╝  ╚════██║██╔═══╝ ██╔══██║██║     ██╔══╝  
██║ ╚████║██║  ██║██║ ╚═╝ ██║███████╗███████║██║     ██║  ██║╚██████╗███████╗
╚═╝  ╚═══╝╚═╝  ╚═╝╚═╝     ╚═╝╚══════╝╚══════╝╚═╝     ╚═╝  ╚═╝ ╚═════╝╚══════╝`

const podBanner = `
██████╗  ██████╗ ██████╗ 
██╔══██╗██╔═══██╗██╔══██╗
██████╔╝██║   ██║██║  ██║
██╔═══╝ ██║   ██║██║  ██║
██║     ╚██████╔╝██████╔╝
╚═╝      ╚═════╝ ╚═════╝`

var (
	titleStyle           = lipgloss.NewStyle().MarginLeft(2)
	itemStyle            = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle    = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ff895e"))
	paginationStyle      = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle            = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle        = lipgloss.NewStyle().Margin(1, 0, 2, 4)
	LogoForegroundStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5f00")).Background(lipgloss.Color("#ff5f00")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#e65400")).Background(lipgloss.Color("#e65400")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#cc4b00")).Background(lipgloss.Color("#cc4b00")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#b34100")).Background(lipgloss.Color("#b34100")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#993800")).Background(lipgloss.Color("#993800")),
		lipgloss.NewStyle(),
	}
	LogoBackgroundStyles = []lipgloss.Style{
		lipgloss.NewStyle().Foreground(lipgloss.Color("255")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("252")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("249")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("246")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("243")),
		lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
	}
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

func getBanner(banner string) string {
	trimmedBanner := strings.TrimSpace(banner)
	var finalBanner strings.Builder

	for i, s := range strings.Split(trimmedBanner, "\n") {
		if i > 0 {
			finalBanner.WriteRune('\n')
		}

		foreground := LogoForegroundStyles[i]
		background := LogoBackgroundStyles[i]

		for _, c := range s {
			if c == '█' {
				finalBanner.WriteString(foreground.Render("█"))
			} else if c != ' ' {
				finalBanner.WriteString(background.Render(string(c)))
			} else {
				finalBanner.WriteRune(c)
			}
		}
	}
	return finalBanner.String()
}

type model struct {
	items              list.Model
	namespace          string
	pod                string
	clientset          kubernetes.Clientset
	selectingNamespace bool
	banner             string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.items.SetWidth(msg.Width)
		m.items.SetHeight(msg.Height - lipgloss.Height(m.banner) - 1)
		return m, nil
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.items.SelectedItem().(item)
			if ok {
				if m.selectingNamespace {
					m.banner = getBanner(podBanner)
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

func getList(items []list.Item) list.Model {
	length := minInt(len(items)+7, 20)
	l := list.New(items, itemDelegate{}, 60, length)
	l.SetShowStatusBar(false)
	l.SetShowTitle(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	return l
}

func buildPodList(namespace string, clientset kubernetes.Clientset) list.Model {
	items := buildPodItems(namespace, clientset)
	return getList(items)
}

func buildNamespaceModel() *model {
	clientset := *getKubernetesClientset()
	items := buildNamespaceItems(clientset)
	return &model{items: getList(items), clientset: clientset, selectingNamespace: true, banner: getBanner(namespaceBanner)}
}

func (m *model) View() string {
	banner := lipgloss.NewStyle().Margin(2, 0, 0, 2).Render(m.banner)
	if m.pod != "" && m.namespace != "" {
		return quitTextStyle.Render(fmt.Sprintf("Opening a shell into %s/%s", m.namespace, m.pod))
	}
	items := m.items.View()
	return lipgloss.JoinVertical(lipgloss.Left, banner, items)
}

func minInt(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
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
