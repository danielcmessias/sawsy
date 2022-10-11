package table

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/components/listviewport"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/context"

	"github.com/danielcmessias/sawsy/utils"
)

type Model struct {
	pane.Pane

	ctx          *context.ProgramContext
	search       textinput.Model
	Columns      []Column
	rows         []Row
	filteredRows []Row
	filterText   string
	rowsViewport listviewport.Model
	width        int
	height       int
	currColumnId int
	colMaxWidths []int
	noDataLabel  string
}

type Column struct {
	Title string
	// MaxWidth is unused for now
	MaxWidth   *int
	isSelected *bool
}

type Row []string

type TableSpec struct {
	pane.BaseSpec

	Columns []Column
}

func (s TableSpec) NewFromSpec(ctx *context.ProgramContext, spec pane.PaneSpec) pane.Pane {
	tableSpec, ok := spec.(TableSpec)
	if !ok {
		log.Fatal("invalid spec type, expected TableSpec")
	}
	return New(ctx, tableSpec)
}

func New(ctx *context.ProgramContext, spec TableSpec) *Model {
	columns := make([]Column, len(spec.Columns))
	copy(columns, spec.Columns)

	columns[0].isSelected = utils.BoolPtr(true)

	search := textinput.New()
	search.Prompt = "Search: "
	search.Placeholder = "..."
	search.PromptStyle = promptStyle
	search.Width = 40

	colMaxWidths := make([]int, len(columns))
	// Set the max width default to the length of the column header
	for i, c := range columns {
		colMaxWidths[i] = lipgloss.Width(titleCellStyle.Render(c.Title))
	}

	return &Model{
		Pane: pane.New(spec.BaseSpec),

		ctx:          ctx,
		search:       search,
		Columns:      columns,
		filterText:   "",
		rowsViewport: listviewport.NewModel(spec.BaseSpec.Name, 0, 2),
		colMaxWidths: colMaxWidths,
		noDataLabel:  "Loading...",
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (pane.Pane, tea.Cmd, bool) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	m.search, cmd = m.search.Update(msg)
	cmds = append(cmds, cmd)
	m.filter(m.search.Value())

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.ctx.Keys.Down):
			m.rowsViewport.NextItem()
		case key.Matches(msg, m.ctx.Keys.Up):
			m.rowsViewport.PrevItem()
		case key.Matches(msg, m.ctx.Keys.FirstLine):
			m.rowsViewport.FirstItem()
		case key.Matches(msg, m.ctx.Keys.LastLine):
			m.rowsViewport.LastItem()
		case key.Matches(msg, m.ctx.Keys.NextCol):
			m.nextCol()
		case key.Matches(msg, m.ctx.Keys.PrevCol):
			m.prevCol()
		case key.Matches(msg, m.ctx.Keys.StartSearch):
			m.search.Focus()
			m.ctx.LockKeyboardCapture = true
		case key.Matches(msg, m.ctx.Keys.EndSearch):
			m.search.Blur()
			m.ctx.LockKeyboardCapture = false
		}
		m.syncViewPortContent()
	}

	return m, tea.Batch(cmds...), false
}

func (m *Model) View() string {
	header := m.renderHeader()
	body := m.renderBody()

	var search string
	if m.search.Focused() || m.search.Value() != "" {
		search = m.search.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		search,
	)
}

func (m *Model) SetSize(width int, height int) {
	m.width = width
	m.height = height
	m.rowsViewport.SetSize(
		width,
		height-headerHeight-searchHeight,
	)
	m.syncViewPortContent()
}

func (m *Model) ResetCurrentItem() {
	m.rowsViewport.ResetCurrItem()
}

func (m *Model) GetCurrentItem() int {
	return m.rowsViewport.GetCurrItem()
}

func (m *Model) GetCurrentRow() Row {
	return m.filteredRows[m.rowsViewport.GetCurrItem()]
}

func (m *Model) GetMarshalledRow() map[string]string {
	row := m.GetCurrentRow()
	rowMap := make(map[string]string)
	for i, col := range m.Columns {
		rowMap[col.Title] = row[i]
	}
	return rowMap
}

func (m *Model) nextCol() int {
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(false)
	m.currColumnId = (m.currColumnId + 1) % len(m.Columns)
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(true)
	return m.currColumnId
}

func (m *Model) prevCol() int {
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(false)
	m.currColumnId = (m.currColumnId - 1) % len(m.Columns)
	if m.currColumnId < 0 {
		m.currColumnId = len(m.Columns) - 1
	}
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(true)
	return m.currColumnId
}

func (m *Model) syncViewPortContent() {
	headerColumns := m.renderHeaderColumns()
	renderedRows := make([]string, 0, len(m.filteredRows))
	for i := range m.filteredRows {
		renderedRows = append(renderedRows, m.renderRow(i, headerColumns))
	}
	m.rowsViewport.SyncViewPort(
		lipgloss.JoinVertical(lipgloss.Left, renderedRows...),
	)
}

func (m *Model) GetRowAt(index int) Row {
	return m.filteredRows[index]
}

func (m *Model) SetRows(rows []Row) {
	m.rows = rows
	m.filterRows()
}

func (m *Model) AppendRows(rows []Row) {
	newRows := append(m.rows, rows...)
	m.SetRows(newRows)

	if len(rows) == 0 {
		m.noDataLabel = "No data"
	}
}

func (m *Model) ClearRows() {
	m.rows = make([]Row, 0)
	m.filterRows()
	m.noDataLabel = "Loading..."
}

func (m *Model) OnLineDown() {
	m.rowsViewport.NextItem()
}

func (m *Model) OnLineUp() {
	m.rowsViewport.PrevItem()
}

func (m *Model) filter(filter string) {
	m.filterText = filter
	m.filterRows()
}

func (m *Model) Hide() {
	m.search.Blur()
}

func (m *Model) filterRows() {
	if m.filterText == "" {
		m.filteredRows = m.rows
	}
	filteredRows := make([]Row, 0)
	for _, r := range m.rows {
		for _, c := range r {
			if strings.Contains(c, m.filterText) {
				filteredRows = append(filteredRows, r)
				break
			}
		}
	}
	m.filteredRows = filteredRows

	for _, row := range m.filteredRows {
		for j, col := range m.Columns {
			w := lipgloss.Width(cellStyle.Copy().Render(row[j]))
			m.colMaxWidths[j] = utils.Max(m.colMaxWidths[j], w)
			if col.MaxWidth != nil {
				m.colMaxWidths[j] = utils.Min(m.colMaxWidths[j], *col.MaxWidth)
			}
		}
	}

	m.rowsViewport.SetNumItems(len(m.filteredRows))
	m.syncViewPortContent()
}

func (m *Model) renderHeaderColumns() []string {
	/* The logic here is basically that in the first pass all columns are assigned the width of
	 * their title. Then the selected column is allowed to grow as much as possible. Finally, all
	 * remaining columns (from left to rigth) are allowed to grow as much as possible. All columns
	 * never grow beyond the largest string in that column, or MaxWidth if set.
	 *
	 * This could be improved with proper horizontal scrolling!
	 */

	renderedColumns := make([]string, len(m.Columns))
	remainingWidth := m.width
	var width int

	for i, column := range m.Columns {
		width = lipgloss.Width(titleCellStyle.Copy().Render(column.Title))
		if i != m.currColumnId {
			renderedColumns[i] = titleCellStyle.
				Copy().
				Width(width).
				MaxWidth(width).
				Render(column.Title)
			remainingWidth -= width
		}
	}

	width = utils.Min(remainingWidth, m.colMaxWidths[m.currColumnId])
	renderedColumns[m.currColumnId] = selectedTitleCellStyle.Copy().
		Width(width).
		MaxWidth(width).
		Render(m.Columns[m.currColumnId].Title)
	remainingWidth -= width

	for i := range m.Columns {
		if i != m.currColumnId && remainingWidth-m.colMaxWidths[i] > 0 {
			renderedColumns[i] = titleCellStyle.Copy().
				Width(m.colMaxWidths[i]).
				MaxWidth(m.colMaxWidths[i]).
				Render(m.Columns[i].Title)
			remainingWidth -= m.colMaxWidths[i]
		}
	}

	return renderedColumns
}

func (m *Model) renderHeader() string {
	headerColumns := m.renderHeaderColumns()
	header := lipgloss.JoinHorizontal(lipgloss.Top, headerColumns...)
	return headerStyle.Copy().
		Width(m.width).
		MaxWidth(m.width).
		Render(header)
}

func (m *Model) renderBody() string {
	bodyStyle := lipgloss.NewStyle().
		Height(m.height - headerHeight)
	if len(m.filteredRows) == 0 {
		return bodyStyle.Copy().PaddingLeft(1).Render(m.noDataLabel)
	}
	return m.rowsViewport.View()
}

func (m *Model) renderRow(rowId int, headerColumns []string) string {
	// Rendering is slow, only draw rows that will actually be visible
	rowsPerPage := m.rowsViewport.GetNumRowsPerPage()
	if rowId < m.rowsViewport.GetCurrItem()-rowsPerPage || rowId > m.rowsViewport.GetCurrItem()+rowsPerPage {
		return "\n"
	}

	var style lipgloss.Style
	if m.rowsViewport.GetCurrItem() == rowId {
		style = selectedCellStyle
	} else {
		style = cellStyle
	}

	renderedColumns := make([]string, len(m.Columns))
	for i, column := range m.filteredRows[rowId] {
		colWidth := lipgloss.Width(headerColumns[i])
		col := style.Copy().Width(colWidth).MaxWidth(colWidth).Render(column)
		renderedColumns = append(renderedColumns, col)
	}

	return rowStyle.Copy().Render(
		lipgloss.JoinHorizontal(lipgloss.Top, renderedColumns...),
	)
}
