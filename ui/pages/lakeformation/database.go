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
	return func() tea.Msg {
		output, _ := client.LakeFormation.GetDatabase(m.Context.(DatabasePageContext).DatabaseName)
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
