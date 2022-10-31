package iam

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type IAMPageModel struct {
	page.Model
}

func NewIAMPage(ctx *context.ProgramContext) *IAMPageModel {
	return &IAMPageModel{
		Model: page.New(ctx, iamPageSpec),
	}
}

func (m *IAMPageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchUsers(client, nil),
		m.fetchRoles(client, nil),
	)
}

func (m *IAMPageModel) fetchUsers(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.IAM.GetUsers(nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Users"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchUsers(client, nextToken)
		}
		return msg
	}
}

func (m *IAMPageModel) fetchRoles(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.IAM.GetRoles(nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Roles"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchRoles(client, nextToken)
		}
		return msg
	}
}

func (m *IAMPageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	row := table.GetCurrentRowMarshalled()

	switch m.GetCurrentPaneId() {
	case m.GetPaneId("Users"):
		return func() tea.Msg {
			return page.ChangePageMsg{
				NewPage:   "iam/user",
				FetchData: true,
				PageContext: UserPageContext{
					UserName: row["Name"],
				},
			}
		}
	case m.GetPaneId("Roles"):
		return func() tea.Msg {
			return page.ChangePageMsg{
				NewPage:   "iam/role",
				FetchData: true,
				PageContext: RolePageContext{
					RoleName: row["Name"],
				},
			}
		}
	default:
		return nil
	}
}
