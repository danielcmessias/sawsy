package page

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui/components/help"
	"github.com/danielcmessias/lfq/ui/components/table"
	"github.com/danielcmessias/lfq/ui/components/tabs"
	"github.com/danielcmessias/lfq/ui/constants"
	"github.com/danielcmessias/lfq/ui/context"
)

const (
	PAGE_HOME     = iota
	PAGE_DATABASE = iota
	PAGE_TABLE    = iota
)

type NewRowsMsg struct {
	Page      string
	TabId     int
	Rows      []table.Row
	NextCmd   tea.Cmd
	Overwrite bool
}

type BatchedNewRowsMsg struct {
	Msgs []NewRowsMsg
}

type ChangePageMsg struct {
	NewPage      string // Id of page to switch to
	FetchData    bool // If true, clears and (re)fetches data on new page
	PageMetadata interface{}

	InspectedRow map[string]string
}

// Logical section with a set of tables/tabs
type Model struct {
	Id     int
	Ctx    context.ProgramContext
	Tabs   tabs.Model
	Tables []table.Model

	Spec PageSpec

	Metadata interface{}
}

type Page interface {
	Init() tea.Cmd
	View() string
	UpdateProgramContext(ctx *context.ProgramContext)
	UpdateSearch(msg tea.Msg) tea.Cmd
	NextTab() int
	PrevTab() int
	NextItem() int
	PrevItem() int
	FirstItem() int
	LastItem() int
	NextCol() int
	PrevCol() int
	StartSearch()
	StopSearch()
	FetchDataCmd(client data.Client) tea.Cmd
	ClearData()
	AppendRows(tabId int, rows []table.Row)
	ClearRows(tabId int)
	InspectFieldCmd(client data.Client) tea.Cmd
	SetPageMetadata(metadata interface{})
	
	
	FromInspectedRow(marshalledRow map[string]string)
	InspectRow(client data.Client) tea.Cmd
}

func PageModelFromSpecs(id int, specs []table.TableSpec) Model {
	var tables []table.Model
	var tabNames []string
	for _, spec := range specs {
		tables = append(tables, table.NewModel(spec.Columns, nil, spec.Name))
		tabNames = append(tabNames, fmt.Sprintf("%s %s", spec.Icon, spec.Name))
	}
	return Model{
		Tabs:   tabs.NewModel(tabNames),
		Tables: tables,
	}
}

func New(spec PageSpec) Model {
	var tables []table.Model
	var tabNames []string
	for _, spec := range spec.TableSpecs {
		tables = append(tables, table.NewModel(spec.Columns, nil, spec.Name))
		tabNames = append(tabNames, fmt.Sprintf("%s %s", spec.Icon, spec.Name))
	}
	return Model{
		Tabs:   tabs.NewModel(tabNames),
		Tables: tables,

		Spec: spec,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Tabs.View(m.Ctx),
		m.CurrentTable().View(),
	)
}

func (m *Model) UpdateProgramContext(ctx *context.ProgramContext) {
	m.Ctx = *ctx
	tableDimensions := constants.Dimensions{
		Width:  ctx.ScreenWidth,
		Height: ctx.ScreenHeight - tabs.TabsHeight - help.HelpHeight,
	}
	for i := range m.Tables {
		m.Tables[i].SetDimensions(tableDimensions)
	}
}

func (m *Model) CurrentTable() *table.Model {
	return &m.Tables[m.Tabs.CurrentTabId]
}

func (m *Model) NextTab() int {
	m.StopSearch()
	return m.Tabs.NextTab()
}

func (m *Model) PrevTab() int {
	m.StopSearch()
	return m.Tabs.PrevTab()
}

func (m *Model) NextItem() int {
	return m.CurrentTable().NextItem()
}

func (m *Model) PrevItem() int {
	return m.CurrentTable().PrevItem()
}

func (m *Model) FirstItem() int {
	return m.CurrentTable().FirstItem()
}

func (m *Model) LastItem() int {
	return m.CurrentTable().LastItem()
}

func (m *Model) NextCol() int {
	return m.CurrentTable().NextCol()
}

func (m *Model) PrevCol() int {
	return m.CurrentTable().PrevCol()
}

func (m *Model) StartSearch() {
	m.CurrentTable().Search.Focus()
}

func (m *Model) StopSearch() {
	m.CurrentTable().Search.Blur()
}

func (m *Model) AppendRows(tabId int, rows []table.Row) {
	m.Tables[tabId].AppendRows(rows)
}

func (m *Model) ClearRows(tabId int) {
	m.Tables[tabId].ClearRows()
}

func (m *Model) FetchDataCmd(client data.Client) tea.Cmd {
	return nil
}

func (m *Model) ClearData() {
	for i := range(m.Tables) {
		m.Tables[i].ClearRows()
	}
}

func (m *Model) UpdateSearch(msg tea.Msg) tea.Cmd {
	var (
		cmds []tea.Cmd
		cmd tea.Cmd
	)
	for i := range(m.Tables) {
		m.Tables[i].Search, cmd = m.Tables[i].Search.Update(msg)
		cmds = append(cmds, cmd)
	}
	m.CurrentTable().Filter(m.CurrentTable().Search.Value())
	return tea.Batch(cmds...)
}

func (m *Model) InspectFieldCmd(client data.Client) tea.Cmd  {
	return nil
}

func (m *Model) SetPageMetadata(metadata interface{}) {
	m.Metadata = metadata
}

func (m *Model) FromInspectedRow(marshalledRow map[string]string) {
	log.Fatal("Not implemented")
}

func (m *Model) InspectRow(client data.Client) tea.Cmd  {
	nextPage := m.Spec.TableSpecs[m.Tabs.CurrentTabId].InspectRowPage
	if (nextPage != "") {
		changePageCmd := func() tea.Msg {
			return ChangePageMsg{
				NewPage: nextPage,
				FetchData: true,
				InspectedRow: m.CurrentTable().GetMarshalledRow(),
			}
		}
		return changePageCmd
	}
	return nil
}
