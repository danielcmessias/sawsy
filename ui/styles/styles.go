package styles

import "github.com/charmbracelet/lipgloss"

var (
	MainTextStyle = lipgloss.NewStyle().
			Foreground(Theme.MainText).
			Bold(true)

	FooterHeight = 3
	FooterStyle  = lipgloss.NewStyle().
			Height(FooterHeight - 1).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(Theme.Border)
)
