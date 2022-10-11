package listviewport

import (
	"fmt"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/utils"
)

type Model struct {
	viewport       viewport.Model
	TopBoundId     int
	BottomBoundId  int
	currId         int
	ListItemHeight int
	NumItems       int
	ItemTypeLabel  string
}

func NewModel(itemTypeLabel string, numItems, listItemHeight int) Model {
	model := Model{
		NumItems:       numItems,
		ListItemHeight: listItemHeight,
		currId:         0,
		viewport:       viewport.Model{},
		TopBoundId:     1,
		ItemTypeLabel:  itemTypeLabel,
	}
	model.BottomBoundId = utils.Min(model.NumItems-1, model.GetNumRowsPerPage()-1)
	return model
}

func (m *Model) SetNumItems(numItems int) {
	m.NumItems = numItems
	m.BottomBoundId = utils.Min(m.NumItems-1, m.GetNumRowsPerPage()-1)
}

func (m *Model) SyncViewPort(content string) {
	m.viewport.SetContent(content)
}

func (m *Model) GetNumRowsPerPage() int {
	return (m.viewport.Height / m.ListItemHeight)
}

func (m *Model) ResetCurrItem() {
	m.currId = 0
}

func (m *Model) GetCurrItem() int {
	return m.currId
}

func (m *Model) NextItem() int {
	atBottomOfViewport := m.currId >= m.BottomBoundId
	if atBottomOfViewport && !m.viewport.AtBottom() {
		m.TopBoundId += 1
		m.BottomBoundId += 1
		m.viewport.LineDown(m.ListItemHeight)
	}

	newId := utils.Min(m.currId+1, m.NumItems-1)
	newId = utils.Max(newId, 0)
	m.currId = newId

	return m.currId
}

func (m *Model) PrevItem() int {
	atTopOfViewport := m.currId <= m.TopBoundId
	if atTopOfViewport && !m.viewport.AtTop() {
		m.TopBoundId -= 1
		m.BottomBoundId -= 1
		m.viewport.LineUp(m.ListItemHeight)
	}

	m.currId = utils.Max(m.currId-1, 0)
	return m.currId
}

func (m *Model) FirstItem() int {
	m.currId = 0
	m.viewport.GotoTop()
	return m.currId
}

func (m *Model) LastItem() int {
	m.currId = m.NumItems - 1
	m.viewport.GotoBottom()
	return m.currId
}

func (m *Model) SetSize(width int, height int) {
	m.viewport.Width = width
	m.viewport.Height = height - pagerHeight
}

func (m *Model) View() string {
	pagerContent := ""
	if m.NumItems > 0 {
		pagerContent = fmt.Sprintf(
			"%v %v/%v",
			m.ItemTypeLabel,
			m.currId+1,
			m.NumItems,
		)
	}
	viewport := m.viewport.View()

	pager := pagerStyle.Copy().Render(pagerContent)

	return lipgloss.NewStyle().
		Width(m.viewport.Width).
		MaxWidth(m.viewport.Width).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			viewport,
			pager,
		))
}
