package styles

import "github.com/charmbracelet/lipgloss"

type Theme struct {
    MainText           lipgloss.AdaptiveColor
    BrightMainText     lipgloss.AdaptiveColor
    Border             lipgloss.AdaptiveColor
    FaintBorder        lipgloss.AdaptiveColor
    FaintText          lipgloss.AdaptiveColor
    SelectedBackground lipgloss.AdaptiveColor
    SelectedTab        lipgloss.AdaptiveColor
    SelectedColHeader  lipgloss.AdaptiveColor
    SearchPrompt       lipgloss.AdaptiveColor
    AWSAccountId       lipgloss.AdaptiveColor
}

var DefaultTheme = Theme{
    MainText:           lipgloss.AdaptiveColor{Light: draculaForeground, Dark: draculaForeground},
    BrightMainText:     lipgloss.AdaptiveColor{Light: draculaPink, Dark: draculaPink},
    Border:             lipgloss.AdaptiveColor{Light: draculaCurrentLine, Dark: draculaCurrentLine},
    FaintBorder:        lipgloss.AdaptiveColor{Light: "#2b2b40", Dark: "#2b2b40"},
    FaintText:          lipgloss.AdaptiveColor{Light: draculaComment, Dark: draculaComment},
    SelectedBackground: lipgloss.AdaptiveColor{Light: draculaCurrentLine, Dark: draculaCurrentLine},
    SelectedTab:        lipgloss.AdaptiveColor{Light: draculaPurple, Dark: draculaPurple},
    SelectedColHeader:  lipgloss.AdaptiveColor{Light: draculaGreen, Dark: draculaGreen},
    SearchPrompt:       lipgloss.AdaptiveColor{Light: draculaGreen, Dark: draculaGreen},
    AWSAccountId:       lipgloss.AdaptiveColor{Light: draculaYellow, Dark: draculaYellow},
}

var (
    SingleRuneWidth    = 4
    MainContentPadding = 1
)
