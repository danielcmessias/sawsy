package lakeformation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
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

func (m *DatabasePageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchDatabaseDetails(client),
		m.fetchDatabaseTags(client),
	)
}

func (m *DatabasePageModel) fetchDatabaseDetails(client data.Client) tea.Cmd {
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

func (m *DatabasePageModel) fetchDatabaseTags(client data.Client) tea.Cmd {
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
