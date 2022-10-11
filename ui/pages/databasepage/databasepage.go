package databasepage

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui/components/page"
)

type DatabasePageModel struct {
	databaseName string
	page.Model
}

type DatabasePage interface {
	page.Page
}

type PageMetadata struct {
	DatabaseName string
}

func NewDatabasesPage() DatabasePageModel {
	return DatabasePageModel{
		Model: page.New(databasesPageSpec),
	}
}

func (m *DatabasePageModel) FetchDataCmd(client data.Client) tea.Cmd {
	fetchDetailsCmd := func() tea.Msg {
		output, _ := client.FetchDatabaseDetailRows(m.databaseName)
		var msgs []page.NewRowsMsg
		for i, rows := range output {
			msgs = append(msgs, page.NewRowsMsg{
				Page:      m.Spec.Name,
				TabId:     i,
				Rows:      rows,
				Overwrite: true,
			})
		}
		return page.BatchedNewRowsMsg{
			Msgs: msgs,
		}
	}
	return tea.Batch(fetchDetailsCmd)	
}

func (m *DatabasePageModel) SetPageMetadata(metadata interface{}) {
	m.databaseName = metadata.(PageMetadata).DatabaseName
}

