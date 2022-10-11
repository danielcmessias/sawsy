package iam

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type UserPageModel struct {
	page.Model
}

type UserPageContext struct {
	UserName string
}

func NewUserPage(ctx *context.ProgramContext) *UserPageModel {
	return &UserPageModel{
		Model: page.New(ctx, userPageSpec),
	}
}

func (m *UserPageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchUserPolicies(client, nil),
		m.fetchUserTags(client, nil),
	)
}

func (m *UserPageModel) fetchUserPolicies(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.IAM.GetUserPolicies(m.Context.(UserPageContext).UserName, nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Policies"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchUserPolicies(client, nextToken)
		}
		return msg
	}
}

func (m *UserPageModel) fetchUserTags(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.IAM.GetUserTags(m.Context.(UserPageContext).UserName, nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Tags"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchUserTags(client, nextToken)
		}
		return msg
	}
}

func (m *UserPageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	switch m.GetCurrentPaneId() {
	case m.GetPaneId("Policies"):
		row := table.GetMarshalledRow()
		return func() tea.Msg {
			return page.ChangePageMsg{
				NewPage:   "iam/policy",
				FetchData: true,
				PageContext: PolicyPageContext{
					UserName:   m.Context.(UserPageContext).UserName,
					PolicyName: row["Name"],
					PolicyArn:  row["ARN"],
				},
			}
		}
	default:
		return nil
	}
}
