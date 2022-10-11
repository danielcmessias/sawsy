package tablepage

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui/components/page"
)

type TablePageModel struct {
    page.Model
    tableName    string
    databaseName string
}

type TablePage interface {
    page.Page
}

type PageMetadata struct {
    TableName    string
    DatabaseName string
}

func NewTablePage() TablePageModel {
	return TablePageModel{
		Model: page.New(databasesPageSpec),
	}
}


func (m *TablePageModel) FetchDataCmd(client data.Client) tea.Cmd {
	fetchDetailsCmd := func() tea.Msg {
		output, _ := client.FetchTableDetailRows(m.tableName, m.databaseName)
		var msgs []page.NewRowsMsg
		for i, rows := range output {
			msgs = append(msgs, page.NewRowsMsg{
				Page: m.Spec.Name,
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

func (m *TablePageModel) SetPageMetadata(metadata interface{}) {
    m.tableName = metadata.(PageMetadata).TableName
    m.databaseName = metadata.(PageMetadata).DatabaseName
}
