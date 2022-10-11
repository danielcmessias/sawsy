package tabs

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/lfq/ui/context"
)

type Model struct {
    TabNames     []string
    CurrentTabId int
}

func NewModel(tabNames []string) Model {
    return Model{
        TabNames: tabNames,
        CurrentTabId: 0,
    }
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
    return m, nil
}

func (m Model) View(ctx context.ProgramContext) string {
    var tabs []string

    for i, name := range m.TabNames {
        if (i == m.CurrentTabId) {
            tabs = append(tabs, activeTab.Render(name))
        } else {
            tabs = append(tabs, tab.Render(name))
        }
    }

    accountId := activeAwsAccount.Render(fmt.Sprintf(" (%s)", ctx.AwsAccountId))
    
    tabsWidth := ctx.ScreenWidth - lipgloss.Width(accountId)

    renderedTabs := lipgloss.NewStyle().
        Width(tabsWidth).
        MaxWidth(tabsWidth).
        Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))


    return tabsRow.Copy().
        Width(ctx.ScreenWidth).
        MaxWidth(ctx.ScreenWidth).
        Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs, accountId))

}

func (m *Model) NextTab() int {
    m.CurrentTabId = (m.CurrentTabId + 1) % len(m.TabNames)
    return m.CurrentTabId
}

func (m *Model) PrevTab() int {
    m.CurrentTabId = (m.CurrentTabId - 1) % len(m.TabNames)
    if m.CurrentTabId < 0 {
        m.CurrentTabId = len(m.TabNames) - 1
    }
    return m.CurrentTabId
}
