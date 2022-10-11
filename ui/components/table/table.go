package table

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/lfq/ui/components/listviewport"
	"github.com/danielcmessias/lfq/ui/constants"
	"github.com/danielcmessias/lfq/utils"
)

type Model struct {
	Search       textinput.Model
	Columns      []Column
	rows         []Row
	filteredRows []Row
	filterText   string
	dimensions   constants.Dimensions
	rowsViewport listviewport.Model
	currColumnId int
	
	colMaxWidths []int
}

type Column struct {
	Title       string
	// MaxWidth is unused for now
	MaxWidth    *int
	isSelected  *bool
}

type Row []string

type TableSpec struct {
	Name           string
	Icon           string
	Columns        []Column
	InspectRowPage string
}

func NewModel(columns []Column, rows []Row, itemTypeLabel string) Model {
	columns[0].isSelected = utils.BoolPtr(true)
	dimensions := constants.Dimensions{}

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

	return Model{
		Search:       search,
		Columns:      columns,
		filteredRows: rows,
		filterText:   "",
		rows:         rows,
		dimensions:   dimensions,
		rowsViewport: listviewport.NewModel(dimensions, itemTypeLabel, len(rows), 2),
		colMaxWidths: colMaxWidths,
	}
}

func (m Model) View() string {
	header := m.renderHeader()
	body := m.renderBody()

	var search string
	if (m.Search.Focused() || m.Search.Value() != "") {
		search = m.Search.View()
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		body,
		search,
	)
}

func (m *Model) SetDimensions(dimensions constants.Dimensions) {
	m.dimensions = dimensions
	m.rowsViewport.SetDimensions(constants.Dimensions{
		Width:  m.dimensions.Width,
		Height: m.dimensions.Height - headerHeight - searchHeight,
	})
	m.SyncViewPortContent()
}

func (m *Model) ResetCurrItem() {
	m.rowsViewport.ResetCurrItem()
}

func (m *Model) GetCurrItem() int {
	return m.rowsViewport.GetCurrItem()
}

func (m *Model) GetCurrRow() Row {
	return m.filteredRows[m.rowsViewport.GetCurrItem()]
}

func (m *Model) GetMarshalledRow() map[string]string {
	row := m.GetCurrRow()
	rowMap := make(map[string]string)
	for i, col := range m.Columns {
		rowMap[col.Title] = row[i]
	}
	return rowMap
}

func (m *Model) PrevItem() int {
	currItem := m.rowsViewport.PrevItem()
	m.SyncViewPortContent()

	return currItem
}

func (m *Model) NextItem() int {
	currItemId := m.rowsViewport.NextItem()
	m.SyncViewPortContent()

	return currItemId
}

func (m *Model) FirstItem() int {
	currItemId := m.rowsViewport.FirstItem()
	m.SyncViewPortContent()

	return currItemId
}

func (m *Model) LastItem() int {
	currItemId := m.rowsViewport.LastItem()
	m.SyncViewPortContent()

	return currItemId
}

func (m *Model) NextCol() int {
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(false)
	m.currColumnId = (m.currColumnId + 1) % len(m.Columns)
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(true)
	m.SyncViewPortContent()
	return m.currColumnId
}

func (m *Model) PrevCol() int {
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(false)
	m.currColumnId = (m.currColumnId - 1) % len(m.Columns)
	if m.currColumnId < 0 {
		m.currColumnId = len(m.Columns) - 1
	}
	m.Columns[m.currColumnId].isSelected = utils.BoolPtr(true)
	m.SyncViewPortContent()
	return m.currColumnId
}

func (m *Model) SyncViewPortContent() {
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
}

func (m *Model) ClearRows() {
	m.rows = make([]Row, 0)
	m.filterRows()
}

func (m *Model) OnLineDown() {
	m.rowsViewport.NextItem()
}

func (m *Model) OnLineUp() {
	m.rowsViewport.PrevItem()
}

func (m *Model) Filter(filter string) {
	m.filterText = filter
	m.filterRows()
}

func (m *Model) filterRows() {
	if (m.filterText == "") {
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
	m.SyncViewPortContent()
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
	remainingWidth := m.dimensions.Width
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
		if i != m.currColumnId && remainingWidth - m.colMaxWidths[i] > 0 {
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
		Width(m.dimensions.Width).
		MaxWidth(m.dimensions.Width).
		Render(header)
}

func (m *Model) renderBody() string {
	bodyStyle := lipgloss.NewStyle().
		Height(m.dimensions.Height - headerHeight)
	if len(m.filteredRows) == 0 {
		return bodyStyle.Render("Loading....")
	}

	return m.rowsViewport.View()
}

func (m *Model) renderRow(rowId int, headerColumns []string) string {
	// Rendering is slow, only draw rows that will actually be visible
	rowsPerPage := m.rowsViewport.GetNumRowsPerPage()
	if (rowId < m.rowsViewport.GetCurrItem() - rowsPerPage || rowId >  m.rowsViewport.GetCurrItem() + rowsPerPage) {
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
