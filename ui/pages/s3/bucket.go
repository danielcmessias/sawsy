package s3

import (
	"fmt"
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/help"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/components/tabs"
	"github.com/danielcmessias/sawsy/ui/context"
	"github.com/danielcmessias/sawsy/utils/icons"
)

const breadcrumbHeight = 1

type BucketPageModel struct {
	page.Model
}

type BucketPageContext struct {
	Bucket string
	Prefix string
	Region string
}

func NewBucketPage(ctx *context.ProgramContext) *BucketPageModel {
	return &BucketPageModel{
		Model: page.New(ctx, bucketPageSpec),
	}
}

func (m *BucketPageModel) View() string {
	context := m.Context.(BucketPageContext)
	breadcrumb := fmt.Sprintf("s3 > %s/%s", context.Bucket, context.Prefix)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Tabs.View(),
		m.CurrentPane().View(),
		breadcrumb,
	)
}

func (m *BucketPageModel) SetSize(width int, height int) {
	for _, p := range m.Panes {
		p.SetSize(width, height-tabs.TabsHeight-help.HelpHeight-breadcrumbHeight)
	}
}

func (m *BucketPageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchObjects(client, nil),
		m.fetchBucketPolicy(client),
		m.fetchBucketTags(client),
	)
}

func (m *BucketPageModel) fetchObjects(client data.Client, nextToken *string) tea.Cmd {
	context := m.Context.(BucketPageContext)
	return func() tea.Msg {
		rows, nextToken, _ := client.S3.GetObjects(context.Bucket, context.Region, context.Prefix, nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Objects"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchObjects(client, nextToken)
		}
		return msg
	}
}

func (m *BucketPageModel) fetchBucketPolicy(client data.Client) tea.Cmd {
	context := m.Context.(BucketPageContext)
	return func() tea.Msg {
		policy, err := client.S3.GetBucketPolicy(context.Bucket, context.Region)
		ext := ".json"
		if err != nil {
			policy = err.Error()
			ext = ""
		}

		msg := code.NewCodeContentMsg{
			Page:     m.Spec.Name,
			PaneId:   m.GetPaneId("Bucket Policy"),
			Content:  policy,
			Filepath: ext,
		}
		return msg
	}
}

func (m *BucketPageModel) fetchBucketTags(client data.Client) tea.Cmd {
	context := m.Context.(BucketPageContext)
	return func() tea.Msg {
		rows, _ := client.S3.GetBucketTags(context.Bucket, context.Region)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Tags"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *BucketPageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	context := m.Context.(BucketPageContext)
	key := table.GetMarshalledRow()["Key"]
	sanitizedPrefix := context.Prefix + key
	sanitizedPrefix = strings.Replace(sanitizedPrefix, fmt.Sprintf("%s ", icons.FILE), "", 1)
	sanitizedPrefix = strings.Replace(sanitizedPrefix, fmt.Sprintf("%s ", icons.FOLDER), "", 1)

	var nextPage string
	var nextContext interface{}
	if strings.HasPrefix(key, icons.FILE) {
		nextPage = "s3/object"
		nextContext = ObjectPageContext{
			Bucket: context.Bucket,
			Key:    sanitizedPrefix,
			Region: context.Region,
		}
	} else {
		nextPage = "s3/objects"
		nextContext = BucketPageContext{
			Bucket: context.Bucket,
			Region: context.Region,
			Prefix: sanitizedPrefix,
		}
	}

	changePageCmd := func() tea.Msg {
		return page.ChangePageMsg{
			NewPage:     nextPage,
			FetchData:   true,
			PageContext: nextContext,
		}
	}

	table.ResetCurrentItem()
	return changePageCmd
}
