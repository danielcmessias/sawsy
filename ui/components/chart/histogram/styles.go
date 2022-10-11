package histogram

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/styles"
)

var (
	barsStyle = lipgloss.NewStyle().
			Foreground(styles.Theme.HighlightRow)

	axisStyle = lipgloss.NewStyle().
			Foreground(styles.Theme.MainText)
)
