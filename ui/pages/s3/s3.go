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

func (m *S3PageModel) FetchData(client *data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchBuckets(client),
	)
}

func (m *S3PageModel) fetchBuckets(client *data.Client) tea.Cmd {
	return func() tea.Msg {
		var nextCmds []tea.Cmd

		rows, err := client.S3.GetBuckets()
		if err != nil {
			log.Fatal(err)
		}

		for _, row := range rows {
			nextCmds = append(nextCmds, m.fetchBucketRegion(client, row))
		}

		msg := page.NewRowsMsg{
			Page:    m.Spec.Name,
			PaneId:  m.GetPaneId("Buckets"),
			Rows:    rows,
			NextCmd: tea.Batch(nextCmds...),
		}
		return msg
	}
}

func (m *S3PageModel) fetchBucketRegion(client *data.Client, row table.Row) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	return func() tea.Msg {
		region, err := client.S3.GetBucketRegion(row[0])
		if err != nil {
			log.Fatal(err)
		}

		marshalledRow := table.MarhsalRow(row)
		marshalledRow["Region"] = region

		return page.UpdateRowMsg{
			Page:            m.Spec.Name,
			PaneId:          m.GetPaneId("Buckets"),
			Row:             table.UnmarshalRow(marshalledRow),
			PrimaryKeyIndex: 0,
		}
	}
}

func (m *S3PageModel) Inspect(client *data.Client) tea.Cmd {
	if m.Tabs.CurrentTabId != m.GetPaneId("Buckets") {
		return nil
	}

	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	row := table.GetCurrentRowMarshalled()

	// If the region hasn't been found yet, we need to wait.
	// An alternative idea might be to just fetch the region in GetObjects if it's missing at that
	// point
	if row["Region"] == data.LOADING_ALIAS {
		return nil
	}

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
