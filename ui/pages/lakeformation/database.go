package lakeformation

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type DatabasePageModel struct {
	page.Model
}

type DatabasePageContext struct {
	DatabaseName string
}

func NewDatabasePage(ctx *context.ProgramContext) *DatabasePageModel {
	return &DatabasePageModel{
		Model: page.New(ctx, databasePageSpec),
	}
}

func (m *DatabasePageModel) FetchData(client *data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchDatabaseDetails(client),
		m.fetchDatabaseTables(client),
		m.fetchDatabaseTags(client),
	)
}

func (m *DatabasePageModel) fetchDatabaseDetails(client *data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.LakeFormation.GetDatabaseDetails(m.Context.(DatabasePageContext).DatabaseName)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Details"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *DatabasePageModel) fetchDatabaseTables(client *data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.LakeFormation.GetDatabaseTables(m.Context.(DatabasePageContext).DatabaseName)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Tables"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *DatabasePageModel) fetchDatabaseTags(client *data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.LakeFormation.GetDatabaseTags(m.Context.(DatabasePageContext).DatabaseName)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("LF-Tags"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *DatabasePageModel) Inspect(client *data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	row := table.GetCurrentRowMarshalled()
	if row == nil {
		return nil
	}

	var nextPage string
	var pageContext interface{}

	switch m.Tabs.CurrentTabId {
	case m.GetPaneId("Tables"):
		nextPage = "lakeformation/table"
		pageContext = TablePageContext{
			TableName:    row["Table"],
			DatabaseName: row["Database"],
		}
	default:
		return nil
	}

	changePageCmd := func() tea.Msg {
		return page.ChangePageMsg{
			NewPage:     nextPage,
			FetchData:   true,
			PageContext: pageContext,
		}
	}
	return changePageCmd
}
