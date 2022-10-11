package s3

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui/components/page"
)

type BucketsPageModel struct {
    page.Model
}

func NewBucketsPage() BucketsPageModel {
	return BucketsPageModel{
		Model: page.New(bucketsPageSpec),
	}
}

func (m *BucketsPageModel) FetchDataCmd(client data.Client) tea.Cmd {
    return tea.Batch(
        m.fetchBuckets(client, nil),
    )
}

func (m *BucketsPageModel) fetchBuckets(client data.Client, nextToken *string) tea.Cmd {
    return func() tea.Msg {
        rows, nextToken, _ := client.S3.GetBucketsRows(nextToken)
        msg := page.NewRowsMsg{
            Page: m.Spec.Name,
            TabId: 0,
            Rows: rows,
        }
        if (nextToken != nil) {
            msg.NextCmd = m.fetchBuckets(client, nextToken)
        }
        return msg
    }
}

func (m *BucketsPageModel) InspectRow(client data.Client) tea.Cmd  {
    changePageCmd := func() tea.Msg {
        return page.ChangePageMsg{
            NewPage: "s3/objects",
            FetchData: true,
            PageMetadata: ObjectsPageMetadata{
                Bucket: m.CurrentTable().GetMarshalledRow()["Name"],
                Region: m.CurrentTable().GetMarshalledRow()["Region"],
            },
        }
    }
    return changePageCmd
}
