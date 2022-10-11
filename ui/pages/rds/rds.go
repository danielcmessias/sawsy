package rds

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type RDSPageModel struct {
	page.Model
}

func NewRDSPage(ctx *context.ProgramContext) *RDSPageModel {
	return &RDSPageModel{
		Model: page.New(ctx, rdsPageSpec),
	}
}

func (m *RDSPageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchDatabases(client, nil),
	)
}

func (m *RDSPageModel) fetchDatabases(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.RDS.GetDBInstances(nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Databases"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchDatabases(client, nextToken)
		}
		return msg
	}
}

func (m *RDSPageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	row := table.GetMarshalledRow()
	var nextPage string
	var pageContext interface{}

	switch m.Tabs.CurrentTabId {
	case m.GetPaneId("Databases"):
		nextPage = "rds/instance"
		pageContext = InstancePageContext{
			InstanceId: row["Identifier"],
		}
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
