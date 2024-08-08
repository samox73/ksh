package styles

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

var (
	TitleStyle           = lipgloss.NewStyle().MarginLeft(2)
	ItemStyle            = lipgloss.NewStyle().PaddingLeft(4)
	SelectedItemStyle    = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("#ff895e"))
	PaginationStyle      = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	HelpStyle            = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	QuitTextStyle        = lipgloss.NewStyle().Margin(1, 0, 2, 4)
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

func GetBanner(banner string) string {
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