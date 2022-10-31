package page

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/chart/histogram"
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/gallery"
	"github.com/danielcmessias/sawsy/ui/components/help"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/components/tabs"
	"github.com/danielcmessias/sawsy/ui/context"
)

type NewRowsMsg struct {
	Page      string
	PaneId    int
	Rows      []table.Row
	NextCmd   tea.Cmd
	Overwrite bool
}

type BatchedNewRowsMsg struct {
	Msgs []NewRowsMsg
}

type UpdateRowMsg struct {
	Page            string
	PaneId          int
	Row             table.Row
	PrimaryKeyIndex int
}

type ChangePageMsg struct {
	NewPage     string // Id of page to switch to
	FetchData   bool   // If true, clears and (re)fetches data on new page
	PageContext interface{}

	InspectedRow map[string]string
}

type PageSpec struct {
	Name      string
	PaneSpecs []pane.PaneSpec
}

// Logical section with a set of tables/tabs
type Model struct {
	ctx     *context.ProgramContext
	Spec    PageSpec
	Context interface{}
	Tabs    tabs.Model
	Panes   []pane.Pane
}

type Page interface {
	Init() tea.Cmd
	View() string
	NextTab() int
	PrevTab() int
	FetchData(client data.Client) tea.Cmd
	ClearData()
	AppendRows(tabId int, rows []table.Row)
	ClearRows(tabId int)
	GetPageContext() interface{}
	SetPageContext(context interface{})
	GetSpec() PageSpec

	GetPaneAt(index int) pane.Pane
	GetCurrentPaneId() int

	Inspect(client data.Client) tea.Cmd

	Update(data.Client, tea.Msg) (cmd tea.Cmd, consumed bool)

	Hide() // Called when the page is no longer visible
	SetSize(width int, height int)
}

func New(ctx *context.ProgramContext, spec PageSpec) Model {
	var tabsList []tabs.Tab
	var panes []pane.Pane
	for _, s := range spec.PaneSpecs {
		tabsList = append(tabsList, tabs.Tab{
			Name: s.GetName(),
			Icon: s.GetIcon(),
		})
		panes = append(panes, s.NewFromSpec(ctx, s))
	}
	return Model{
		ctx:   ctx,
		Spec:  spec,
		Tabs:  tabs.NewModel(ctx, tabsList),
		Panes: panes,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) View() string {
	pane := m.Panes[m.Tabs.CurrentTabId]

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Tabs.View(),
		pane.View(),
	)
}

func (m *Model) SetSize(width int, height int) {
	for _, p := range m.Panes {
		p.SetSize(width, height-tabs.TabsHeight-help.HelpHeight)
	}
}

func (m *Model) CurrentPane() pane.Pane {
	return m.Panes[m.Tabs.CurrentTabId]
}

func (m *Model) NextTab() int {
	m.CurrentPane().Hide()
	return m.Tabs.NextTab()
}

func (m *Model) PrevTab() int {
	m.CurrentPane().Hide()
	return m.Tabs.PrevTab()
}

func (m *Model) AppendRows(tabId int, rows []table.Row) {
	table, ok := m.Panes[tabId].(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}
	table.AppendRows(rows)
}

func (m *Model) ClearRows(tabId int) {
	table, ok := m.Panes[tabId].(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}
	table.ClearRows()
}

func (m *Model) FetchData(client data.Client) tea.Cmd {
	return nil
}

func (m *Model) ClearData() {
	for _, pane := range m.Panes {
		table, ok := pane.(*table.Model)
		if ok {
			table.ClearRows()
		}
	}
}

func (m *Model) GetPageContext() interface{} {
	return m.Context
}

func (m *Model) SetPageContext(context interface{}) {
	m.Context = context
}

func (m *Model) GetSpec() PageSpec {
	return m.Spec
}

func (m *Model) GetPaneId(paneName string) int {
	for i, p := range m.Panes {
		if p.GetSpec().GetName() == paneName {
			return i
		}
	}
	log.Fatalf("Pane %s not found", paneName)
	return -1
}

func (m *Model) Inspect(client data.Client) tea.Cmd {
	// Not implemented
	return nil
}

func (m *Model) Update(client data.Client, msg tea.Msg) (tea.Cmd, bool) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	_, cmd, consumed := m.CurrentPane().Update(msg)
	cmds = append(cmds, cmd)
	if consumed {
		return tea.Batch(cmds...), consumed
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.ctx.Keys.NextTab):
			m.NextTab()
		case key.Matches(msg, m.ctx.Keys.PrevTab):
			m.PrevTab()
		}
	case code.NewCodeContentMsg:
		if msg.Page == m.Spec.Name {
			m.Panes[msg.PaneId].(*code.Model).SetContent(msg.Content, msg.Filepath)
		}
	case histogram.NewDataMsg:
		if msg.Page == m.Spec.Name {
			// Not entirely happy with how this works, but will do for now
			gallery, ok := m.Panes[msg.PaneId].(*gallery.Model)
			if ok {
				gallery.Panes[msg.GalleryPaneId].(*histogram.Model).SetData(msg.Data)
				break
			}
			m.Panes[msg.PaneId].(*histogram.Model).SetData(msg.Data)
		}
	}

	return tea.Batch(cmds...), false
}

func (m *Model) GetPaneAt(index int) pane.Pane {
	return m.Panes[index]
}

func (m *Model) GetCurrentPaneId() int {
	return m.Tabs.CurrentTabId
}

func (m *Model) Hide() {
	pane := m.CurrentPane()
	table, ok := pane.(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}
	table.Hide()
}
