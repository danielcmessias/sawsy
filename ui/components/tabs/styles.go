package tabs

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/lfq/ui/styles"
)

var (
    tabsBorderHeight   = 1
    tabsContentHeight  = 2
    TabsHeight         = tabsBorderHeight + tabsContentHeight

    tab = lipgloss.NewStyle().
        Faint(true).
        Padding(0, 2)

    activeTab = tab.
        Copy().
        Faint(false).
        Bold(true).
        Foreground(styles.DefaultTheme.SelectedTab).
        BorderBottom(true).
        BorderStyle(lipgloss.ThickBorder()).
        BorderBottomForeground(styles.DefaultTheme.SelectedTab)

    tabsRow = lipgloss.NewStyle().
        Height(tabsContentHeight).
        PaddingTop(1).
        PaddingBottom(0)


    activeAwsAccount = lipgloss.NewStyle().
        // Background(styles.DefaultTheme.FaintBorder).
        Foreground(styles.DefaultTheme.AWSAccountId)
)
