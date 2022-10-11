package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/styles"
)

var (
	headerHeight = 3
	searchHeight = 1

	cellStyle = lipgloss.NewStyle().
			Foreground(styles.Theme.MainText).
			PaddingLeft(1).
			PaddingRight(1).
			MaxHeight(1)

	selectedCellStyle = cellStyle.Copy().
				Bold(true).
				Foreground(styles.Theme.HighlightRow)

	titleCellStyle = cellStyle.Copy().
			Bold(true).
			Foreground(styles.Theme.MainText)

	selectedTitleCellStyle = cellStyle.Copy().
				Bold(true).
				Foreground(styles.Theme.HighlightColumn)

	headerStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(styles.Theme.Border).
			BorderBottom(true)

	rowStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(styles.Theme.FaintBorder).
			BorderBottom(true)

	promptStyle = lipgloss.NewStyle().
			Foreground(styles.Theme.SearchPrompt)
)
