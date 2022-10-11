package gallery

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/context"
)

type Model struct {
	ctx    *context.ProgramContext
	Panes  []pane.Pane
	rows   int
	cols   int
	width  int
	height int

	CurrentPaneId int
	// This is for pagination within the gallery, not the same as the regular Pages
	galleryPageId int

	numberOfPages  int
	fullscreenMode bool

	pane.Pane
}

type GallerySpec struct {
	PaneSpecs []pane.PaneSpec
	Rows      int
	Cols      int

	pane.BaseSpec
}

func (s GallerySpec) NewFromSpec(ctx *context.ProgramContext, spec pane.PaneSpec) pane.Pane {
	gallerySpec, ok := spec.(GallerySpec)
	if !ok {
		log.Fatal("invalid spec type, expected GallerySpec")
	}
	var panes []pane.Pane
	for _, s := range gallerySpec.PaneSpecs {
		panes = append(panes, s.NewFromSpec(ctx, s))
	}
	return New(ctx, panes, gallerySpec)
}

func New(ctx *context.ProgramContext, panes []pane.Pane, spec GallerySpec) *Model {
	// for i := 0; i < 20; i++ {
	// 	spec := histogram.HistogramSpec{
	// 		BaseSpec: pane.BaseSpec{
	// 			Name: fmt.Sprintf("Pane #%d", i),
	// 		},
	// 	}
	// 	panes = append(panes, histogram.New(ctx, spec))
	// }

	return &Model{
		Pane: pane.New(spec.BaseSpec),

		ctx:           ctx,
		Panes:         panes,
		rows:          spec.Rows,
		cols:          spec.Cols,
		numberOfPages: len(panes)/(spec.Rows*spec.Cols) + 1,
	}
}

func (m *Model) Update(msg tea.Msg) (pane.Pane, tea.Cmd, bool) {
	// var cmd tea.Cmd
	var cmds []tea.Cmd
	var consumed bool

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.ctx.Keys.Down):
			m.nextRow()
		case key.Matches(msg, m.ctx.Keys.Up):
			m.prevRow()
		case key.Matches(msg, m.ctx.Keys.NextCol):
			m.nextCol()
		case key.Matches(msg, m.ctx.Keys.PrevCol):
			m.prevCol()
		case key.Matches(msg, m.ctx.Keys.Inspect):
			m.fullscreenMode = !m.fullscreenMode
			consumed = true
			m.updateItemSizes()
		case key.Matches(msg, m.ctx.Keys.PrevPage) && m.fullscreenMode:
			m.fullscreenMode = false
			consumed = true
			m.updateItemSizes()
		}
	}

	return m, tea.Batch(cmds...), consumed
}

func (m *Model) View() string {
	if m.fullscreenMode {
		return m.renderFullscreen()
	} else {
		return m.renderItems()
	}
}

func (m *Model) renderFullscreen() string {
	pane := m.Panes[m.CurrentPaneId]
	return selectedItemStyle.Copy().
		Width(m.width - 2).
		Height(m.height - 2).
		BorderTitle(pane.GetSpec().GetName()).
		Render(pane.View())
}

func (m *Model) renderItems() string {
	var rows []string
	for j := 0; j < m.cols; j++ {
		var cells []string
		for i := 0; i < m.rows; i++ {
			paneId := (m.galleryPageId * m.rows * m.cols) + (j * m.rows) + i
			if paneId >= len(m.Panes) {
				cells = append(cells, "")
				break
			}
			pane := m.Panes[paneId]
			var style lipgloss.Style
			if paneId == m.CurrentPaneId {
				style = selectedItemStyle.Copy()
			} else {
				style = itemStyle.Copy()
			}
			style = style.BorderTitle(pane.GetSpec().GetName())
			cells = append(cells, style.Render(pane.View()))
		}

		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, cells...))
	}
	if m.numberOfPages > 1 {
		rows = append(rows, m.renderPaginator())
	}

	return lipgloss.NewStyle().Height(m.height).Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.updateItemSizes()
}

func (m *Model) renderPaginator() string {
	sl := make([]string, m.numberOfPages)
	for i := 0; i < m.numberOfPages; i++ {
		if i == m.galleryPageId {
			sl[i] = activeDot
		} else {
			sl[i] = inactiveDot
		}
	}
	return paginatorStyle.Render(lipgloss.JoinHorizontal(lipgloss.Center, sl...))
}

func (m *Model) nextRow() {
	m.CurrentPaneId = (m.CurrentPaneId + m.cols) % len(m.Panes)
	m.updateGalleryPageId()
}

func (m *Model) prevRow() {
	m.CurrentPaneId = (m.CurrentPaneId - m.cols) % len(m.Panes)
	if m.CurrentPaneId < 0 {
		m.CurrentPaneId = len(m.Panes) - 1
	}
	m.updateGalleryPageId()
}

func (m *Model) nextCol() {
	m.CurrentPaneId = (m.CurrentPaneId + 1) % len(m.Panes)
	m.updateGalleryPageId()
}

func (m *Model) prevCol() {
	m.CurrentPaneId = (m.CurrentPaneId - 1) % len(m.Panes)
	if m.CurrentPaneId < 0 {
		m.CurrentPaneId = len(m.Panes) - 1
	}
	m.updateGalleryPageId()
}

func (m *Model) updateItemSizes() {
	if m.fullscreenMode {
		for _, p := range m.Panes {
			p.SetSize(m.width-2, m.height-2)
		}
	} else {
		paginatorHeight := 0
		if m.numberOfPages > 1 {
			paginatorHeight = 1
		}
		// Subtract 1 for the border (each side)
		itemWidth := (m.width / m.rows) - 2
		itemHeight := ((m.height - paginatorHeight) / m.cols) - 2

		for _, p := range m.Panes {
			p.SetSize(itemWidth, itemHeight)
		}
	}
}

func (m *Model) updateGalleryPageId() {
	m.galleryPageId = m.CurrentPaneId / (m.rows * m.cols)
}
