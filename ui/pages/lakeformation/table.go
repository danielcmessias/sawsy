package lakeformation

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/context"
)

type TablePageModel struct {
	page.Model
}

type TablePageContext struct {
	TableName    string
	DatabaseName string
}

func NewTablePage(ctx *context.ProgramContext) *TablePageModel {
	return &TablePageModel{
		Model: page.New(ctx, tablePageSpec),
	}
}

func (m *TablePageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchTableDetailsAndSchema(client),
		m.fetchTableTags(client),
	)
}

func (m *TablePageModel) fetchTableDetailsAndSchema(client data.Client) tea.Cmd {
	return func() tea.Msg {
		ctx := m.Context.(TablePageContext)
		detailsRows, schemaRows, _ := client.LakeFormation.GetTableDetailsAndSchema(ctx.TableName, ctx.DatabaseName)

		detailsRowsMsg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Details"),
			Rows:   detailsRows,
		}
		schemaRowsMsg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Schema"),
			Rows:   schemaRows,
		}

		return page.BatchedNewRowsMsg{
			Msgs: []page.NewRowsMsg{detailsRowsMsg, schemaRowsMsg},
		}
	}
}

func (m *TablePageModel) fetchTableTags(client data.Client) tea.Cmd {
	return func() tea.Msg {
		ctx := m.Context.(TablePageContext)
		rows, _ := client.LakeFormation.GetTableTags(ctx.TableName, ctx.DatabaseName)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("LF-Tags"),
			Rows:   rows,
		}
		return msg
	}
}
