package gallery

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/styles"
)

var (
	titleStyle = lipgloss.NewStyle().
			Foreground(styles.Theme.MainText).
			Align(lipgloss.Center)

	selectedTitleStyle = titleStyle.Copy().
				Foreground(styles.Theme.HighlightColumn).
				Bold(true)

	itemStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.Theme.Border).
			BorderTitleStyle(titleStyle)

	selectedItemStyle = itemStyle.Copy().
				BorderForeground(styles.Theme.HighlightColumn).
				BorderTitleStyle(selectedTitleStyle)

	paginatorStyle = lipgloss.NewStyle().
			Align(lipgloss.Center)
	activeDot   = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).PaddingLeft(1).PaddingRight(1).Render("⬤")
	inactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).PaddingLeft(1).PaddingRight(1).Render("⬤")
)
