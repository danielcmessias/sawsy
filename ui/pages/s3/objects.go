package s3

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui/components/help"
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/components/tabs"
	"github.com/danielcmessias/lfq/ui/constants"
	"github.com/danielcmessias/lfq/ui/context"
	"github.com/danielcmessias/lfq/utils/icons"
)


type ObjectsPageModel struct {
    page.Model
}

type ObjectsPageMetadata struct {
	Bucket string
	Prefix string
	Region string
}

func NewObjectsPage() ObjectsPageModel {
	return ObjectsPageModel{
		Model: page.New(objectsPageSpec),
	}
}

func (m *ObjectsPageModel) View() string {

	meta := m.Metadata.(ObjectsPageMetadata)
	breadcrumb := fmt.Sprintf("s3 > %s/%s", meta.Bucket, meta.Prefix)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.Tabs.View(m.Ctx),
		m.CurrentTable().View(),
		breadcrumb,
	)
}

func (m *ObjectsPageModel) UpdateProgramContext(ctx *context.ProgramContext) {
	m.Ctx = *ctx
	breadcrumbHeight := 1
	tableDimensions := constants.Dimensions{
		Width:  ctx.ScreenWidth,
		Height: ctx.ScreenHeight - tabs.TabsHeight - breadcrumbHeight - help.HelpHeight,
	}
	for i := range m.Tables {
		m.Tables[i].SetDimensions(tableDimensions)
	}
}

func (m *ObjectsPageModel) FetchDataCmd(client data.Client) tea.Cmd {
    return tea.Batch(
        m.fetchObjects(client, nil),
    )
}

func (m *ObjectsPageModel) fetchObjects(client data.Client, nextToken *string) tea.Cmd {
	meta := m.Metadata.(ObjectsPageMetadata)
    return func() tea.Msg {
        rows, nextToken, _ := client.S3.GetObjectsRows(meta.Bucket, meta.Region, meta.Prefix, nextToken)
        msg := page.NewRowsMsg{
            Page: m.Spec.Name,
            TabId: 0,
            Rows: rows,
        }
        if (nextToken != nil) {
            msg.NextCmd = m.fetchObjects(client, nextToken)
        }
        return msg
    }
}

func (m *ObjectsPageModel) InspectRow(client data.Client) tea.Cmd  {
	meta := m.Metadata.(ObjectsPageMetadata)
	key := m.CurrentTable().GetMarshalledRow()["Key"]
	
	// For now, skip files
	if strings.HasPrefix(key, icons.FILE) {
		return nil
	}

	sanitizedPrefix := meta.Prefix + key
	sanitizedPrefix = strings.Replace(sanitizedPrefix, fmt.Sprintf("%s ", icons.FILE), "", 1)
	sanitizedPrefix = strings.Replace(sanitizedPrefix, fmt.Sprintf("%s ", icons.FOLDER), "", 1)

    changePageCmd := func() tea.Msg {
        return page.ChangePageMsg{
            NewPage: "s3/objects",
            FetchData: true,
            PageMetadata: ObjectsPageMetadata{
                Bucket: meta.Bucket,
                Region: meta.Region,
				Prefix: sanitizedPrefix,
            },
        }
    }
    return changePageCmd
}
