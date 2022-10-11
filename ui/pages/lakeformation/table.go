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
	context := m.Context.(TablePageContext)
	return func() tea.Msg {
		output, _ := client.LakeFormation.GetTable(context.TableName, context.DatabaseName)
		var msgs []page.NewRowsMsg
		for i, rows := range output {
			msgs = append(msgs, page.NewRowsMsg{
				Page:      m.Spec.Name,
				PaneId:    i,
				Rows:      rows,
				Overwrite: true,
			})
		}
		return page.BatchedNewRowsMsg{
			Msgs: msgs,
		}
	}
}
