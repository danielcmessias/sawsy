package styles

import "github.com/charmbracelet/lipgloss"

type ThemeSpec struct {
	MainText        lipgloss.AdaptiveColor
	FaintText       lipgloss.AdaptiveColor
	PageMetaText    lipgloss.AdaptiveColor
	HighlightTab    lipgloss.AdaptiveColor
	HighlightRow    lipgloss.AdaptiveColor
	HighlightColumn lipgloss.AdaptiveColor
	Border          lipgloss.AdaptiveColor
	FaintBorder     lipgloss.AdaptiveColor
	SearchPrompt    lipgloss.AdaptiveColor
}

var dracula = ThemeSpec{
	MainText:        lipgloss.AdaptiveColor{Light: "#f8f8f2", Dark: "#f8f8f2"},
	FaintText:       lipgloss.AdaptiveColor{Light: "#6272a4", Dark: "#6272a4"},
	PageMetaText:    lipgloss.AdaptiveColor{Light: "#f1fa8c", Dark: "#f1fa8c"},
	HighlightTab:    lipgloss.AdaptiveColor{Light: "#bd93f9", Dark: "#bd93f9"},
	HighlightRow:    lipgloss.AdaptiveColor{Light: "#ff79c6", Dark: "#ff79c6"},
	HighlightColumn: lipgloss.AdaptiveColor{Light: "#50fa7b", Dark: "#50fa7b"},
	Border:          lipgloss.AdaptiveColor{Light: "#44475a", Dark: "#44475a"},
	FaintBorder:     lipgloss.AdaptiveColor{Light: "#2b2b40", Dark: "#2b2b40"},
	SearchPrompt:    lipgloss.AdaptiveColor{Light: "#50fa7b", Dark: "#50fa7b"},
}

// Light is latte, dark is ???
var catppuccin = ThemeSpec{
	MainText:        lipgloss.AdaptiveColor{Light: "#4c4f69", Dark: "#0fff00"},
	FaintText:       lipgloss.AdaptiveColor{Light: "#6c6f85", Dark: "#6272a4"},
	PageMetaText:    lipgloss.AdaptiveColor{Light: "#df8e1d", Dark: "#f1fa8c"},
	HighlightTab:    lipgloss.AdaptiveColor{Light: "#d20f39	", Dark: "#bd93f9"},
	HighlightRow:    lipgloss.AdaptiveColor{Light: "#e64553	", Dark: "#ff79c6"},
	HighlightColumn: lipgloss.AdaptiveColor{Light: "#209fb5", Dark: "#50fa7b"},
	Border:          lipgloss.AdaptiveColor{Light: "#dce0e8", Dark: "#44475a"},
	FaintBorder:     lipgloss.AdaptiveColor{Light: "#e6e9ef", Dark: "#2b2b40"},
	SearchPrompt:    lipgloss.AdaptiveColor{Light: "#7287fd", Dark: "#50fa7b"},
}

var (
	SingleRuneWidth    = 4
	MainContentPadding = 1
)

var Theme = func() ThemeSpec {
	return dracula
	// return catppuccin
}()
