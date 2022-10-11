package tabs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/styles"
)

var (
	tabsBorderHeight  = 1
	tabsContentHeight = 2
	TabsHeight        = tabsBorderHeight + tabsContentHeight

	tab = lipgloss.NewStyle().
		Faint(true).
		Padding(0, 2)

	activeTab = tab.
			Copy().
			Faint(false).
			Bold(true).
			Foreground(styles.Theme.HighlightTab).
			BorderBottom(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderBottomForeground(styles.Theme.HighlightTab)

	tabsRow = lipgloss.NewStyle().
		Height(tabsContentHeight).
		PaddingTop(1).
		PaddingBottom(0)

	activeAwsAccount = lipgloss.NewStyle().
				PaddingRight(2).
				Foreground(styles.Theme.PageMetaText)
)
