package s3

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type S3PageModel struct {
	page.Model
}

func NewS3Page(ctx *context.ProgramContext) *S3PageModel {
	return &S3PageModel{
		Model: page.New(ctx, bucketsPageSpec),
	}
}

func (m *S3PageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchBuckets(client, nil),
	)
}

func (m *S3PageModel) fetchBuckets(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.S3.GetBuckets(nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Buckets"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchBuckets(client, nextToken)
		}
		return msg
	}
}

func (m *S3PageModel) Inspect(client data.Client) tea.Cmd {
	if m.Tabs.CurrentTabId != m.GetPaneId("Buckets") {
		return nil
	}

	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}
	row := table.GetMarshalledRow()
	changePageCmd := func() tea.Msg {
		return page.ChangePageMsg{
			NewPage:   "s3/objects",
			FetchData: true,
			PageContext: BucketPageContext{
				Bucket: row["Name"],
				Region: row["Region"],
			},
		}
	}
	return changePageCmd
}
