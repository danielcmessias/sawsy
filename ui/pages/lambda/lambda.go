package lambda

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type LambdaPageModel struct {
	page.Model
}

func NewLambdaPage(ctx *context.ProgramContext) *LambdaPageModel {
	return &LambdaPageModel{
		Model: page.New(ctx, lambdaPageSpec),
	}
}

func (m *LambdaPageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchFunctions(client, nil),
	)
}

func (m *LambdaPageModel) fetchFunctions(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.Lambda.GetFunctions(nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Functions"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchFunctions(client, nextToken)
		}
		return msg
	}
}

func (m *LambdaPageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}
	row := table.GetMarshalledRow()
	var nextPage string
	var pageContext interface{}

	switch m.Tabs.CurrentTabId {
	case m.GetPaneId("Functions"):
		nextPage = "lambda/function"
		pageContext = FunctionPageContext{
			FunctionName: row["Name"],
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
