package table

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/lfq/ui/styles"
)

var (
    headerHeight = 3
    searchHeight = 1
    
    cellStyle = lipgloss.NewStyle().
        PaddingLeft(1).
        PaddingRight(1).
        MaxHeight(1)

    selectedCellStyle = cellStyle.Copy().
        Bold(true).
        Foreground(styles.DefaultTheme.BrightMainText)

    titleCellStyle = cellStyle.Copy().
        Bold(true).
        Foreground(styles.DefaultTheme.MainText)

    selectedTitleCellStyle = cellStyle.Copy().
        Bold(true).
        Foreground(styles.DefaultTheme.SelectedColHeader)

    singleRuneTitleCellStyle = titleCellStyle.Copy().Width(styles.SingleRuneWidth)

    headerStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(styles.DefaultTheme.Border).
        BorderBottom(true)

    rowStyle = lipgloss.NewStyle().
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(styles.DefaultTheme.FaintBorder).
        BorderBottom(true)

    promptStyle = lipgloss.NewStyle().
        Foreground(styles.DefaultTheme.SearchPrompt)
)
