package s3

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/context"
)

type ObjectPageModel struct {
	page.Model
}

type ObjectPageContext struct {
	Bucket string
	Key    string
	Region string
}

func NewObjectPage(ctx *context.ProgramContext) *ObjectPageModel {
	return &ObjectPageModel{
		Model: page.New(ctx, objectPageSpec),
	}
}

func (m *ObjectPageModel) FetchData(client *data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchProperties(client),
	)
}

func (m *ObjectPageModel) fetchProperties(client *data.Client) tea.Cmd {
	context := m.Context.(ObjectPageContext)
	return func() tea.Msg {
		rows, err := client.S3.GetObjectProperties(context.Bucket, context.Key, context.Region)
		if err != nil {
			log.Fatal(err)
		}

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Properties"),
			Rows:   rows,
		}
		return msg
	}
}
