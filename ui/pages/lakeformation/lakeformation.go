package lakeformation

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type LakeFormationPageModel struct {
	page.Model
}

func NewLakeFormationPage(ctx *context.ProgramContext) *LakeFormationPageModel {
	return &LakeFormationPageModel{
		Model: page.New(ctx, lakeFormationPageSpec),
	}
}

func (m *LakeFormationPageModel) FetchData(client *data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchDatabasesCmd(client, nil),
		m.fetchTablesCmd(client, nil),
		m.fetchLFTagsCmd(client, nil),
		m.fetchLFTagPermissionsCmd(client, nil),
		m.fetchDataLakeLocationsCmd(client, nil),
	)
}

func (m *LakeFormationPageModel) fetchDatabasesCmd(client *data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, err := client.LakeFormation.GetDatabases(nextToken)
		if err != nil {
			log.Fatal(err)
		}

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Databases"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchDatabasesCmd(client, nextToken)
		}
		return msg
	}
}

func (m *LakeFormationPageModel) fetchTablesCmd(client *data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, err := client.LakeFormation.GetTables(nextToken)
		if err != nil {
			log.Fatal(err)
		}

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Tables"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchTablesCmd(client, nextToken)
		}
		return msg
	}
}

func (m *LakeFormationPageModel) fetchLFTagsCmd(client *data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, err := client.LakeFormation.GetLFTags(nextToken)
		if err != nil {
			log.Fatal(err)
		}

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("LF-Tags"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchLFTagsCmd(client, nextToken)
		}
		return msg
	}
}

func (m *LakeFormationPageModel) fetchLFTagPermissionsCmd(client *data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, err := client.LakeFormation.GetLFTagPermissions(nextToken)
		if err != nil {
			log.Fatal(err)
		}

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("LF-Tag Perms"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchLFTagPermissionsCmd(client, nextToken)
		}
		return msg
	}
}

func (m *LakeFormationPageModel) fetchDataLakeLocationsCmd(client *data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, err := client.LakeFormation.GetDataLakeLocations(nextToken)
		if err != nil {
			log.Fatal(err)
		}

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("LF Locations"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchDataLakeLocationsCmd(client, nextToken)
		}
		return msg
	}
}

func (m *LakeFormationPageModel) Inspect(client *data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	row := table.GetCurrentRowMarshalled()
	var nextPage string
	var pageContext interface{}

	switch m.Tabs.CurrentTabId {
	case m.GetPaneId("Databases"):
		nextPage = "lakeformation/database"
		pageContext = DatabasePageContext{
			DatabaseName: row["Database"],
		}
	case m.GetPaneId("Tables"):
		nextPage = "lakeformation/table"
		pageContext = TablePageContext{
			TableName:    row["Table"],
			DatabaseName: row["Database"],
		}
	case m.GetPaneId("LF-Tags"):
		// Not implemented
		return nil
	case m.GetPaneId("LF-Tag Perms"):
		// Not implemented
		return nil
	case m.GetPaneId("LF Locations"):
		// Not implemented
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
