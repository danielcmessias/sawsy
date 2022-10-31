package glue

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type GluePageModel struct {
	page.Model
}

func NewGluePage(ctx *context.ProgramContext) *GluePageModel {
	return &GluePageModel{
		Model: page.New(ctx, gluePageSpec),
	}
}

func (m *GluePageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchJobs(client, nil),
		m.fetchCrawlers(client, nil),
	)
}

func (m *GluePageModel) fetchJobs(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.Glue.GetJobsRows(nextToken)

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Jobs"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchJobs(client, nextToken)
		}
		return msg
	}
}

func (m *GluePageModel) fetchCrawlers(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.Glue.GetCrawlersRows(nextToken)

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Crawlers"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchCrawlers(client, nextToken)
		}
		return msg
	}
}

func (m *GluePageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}
	row := table.GetCurrentRowMarshalled()
	var nextPage string
	var pageContext interface{}

	switch m.Tabs.CurrentTabId {
	case m.GetPaneId("Jobs"):
		nextPage = "glue/job"
		pageContext = JobPageContext{
			JobName: row["Name"],
		}
	case m.GetPaneId("Crawlers"):
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
