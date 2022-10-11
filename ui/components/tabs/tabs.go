package tabs

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/context"
	"github.com/danielcmessias/sawsy/utils/icons"
)

type Tab struct {
	Name string
	Icon string
}

type Model struct {
	ctx          *context.ProgramContext
	Tabs         []Tab
	CurrentTabId int
}

func NewModel(ctx *context.ProgramContext, tabs []Tab) Model {
	return Model{
		ctx:          ctx,
		Tabs:         tabs,
		CurrentTabId: 0,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	return m, nil
}

func (m Model) View() string {
	var tabs []string
	for i, t := range m.Tabs {
		tabName := t.Name
		if m.ctx.Config.Theme.ShowIcons {
			tabName = fmt.Sprintf("%s %s", t.Icon, tabName)
		}
		if i == m.CurrentTabId {
			tabs = append(tabs, activeTab.Render(tabName))
		} else {
			tabs = append(tabs, tab.Render(tabName))
		}
	}

	accountId := activeAwsAccount.Render(
		fmt.Sprintf(
			"%s (%s) | %s",
			icons.AWS, m.ctx.AwsAccountId,
			m.ctx.AwsService))

	tabsWidth := m.ctx.ScreenWidth - lipgloss.Width(accountId)

	renderedTabs := lipgloss.NewStyle().
		Width(tabsWidth).
		MaxWidth(tabsWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, tabs...))

	return tabsRow.Copy().
		Width(m.ctx.ScreenWidth).
		MaxWidth(m.ctx.ScreenWidth).
		Render(lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs, accountId))
}

func (m *Model) NextTab() int {
	m.CurrentTabId = (m.CurrentTabId + 1) % len(m.Tabs)
	return m.CurrentTabId
}

func (m *Model) PrevTab() int {
	m.CurrentTabId = (m.CurrentTabId - 1) % len(m.Tabs)
	if m.CurrentTabId < 0 {
		m.CurrentTabId = len(m.Tabs) - 1
	}
	return m.CurrentTabId
}
