package glue

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/context"
)

type JobPageModel struct {
	page.Model
}

type JobPageContext struct {
	JobName string
}

func NewJobsPage(ctx *context.ProgramContext) *JobPageModel {
	return &JobPageModel{
		Model: page.New(ctx, jobPageSpec),
	}
}

func (m *JobPageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchDetails(client),
		m.fetchRuns(client),
		m.fetchScript(client),
	)
}

func (m *JobPageModel) fetchDetails(client data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.Glue.GetJobDetails(m.Context.(JobPageContext).JobName)

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Details"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *JobPageModel) fetchRuns(client data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.Glue.GetJobRuns(m.Context.(JobPageContext).JobName)

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Runs"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *JobPageModel) fetchScript(client data.Client) tea.Cmd {
	return func() tea.Msg {
		script, location, _ := client.Glue.GetJobScript(m.Context.(JobPageContext).JobName)

		msg := code.NewCodeContentMsg{
			Page:     m.Spec.Name,
			PaneId:   m.GetPaneId("Script"),
			Content:  script,
			Filepath: location,
		}
		return msg
	}
}
