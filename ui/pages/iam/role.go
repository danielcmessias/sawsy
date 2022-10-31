package iam

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

type RolePageModel struct {
	page.Model
}

type RolePageContext struct {
	RoleName string
}

func NewRolePage(ctx *context.ProgramContext) *RolePageModel {
	return &RolePageModel{
		Model: page.New(ctx, rolePageSpec),
	}
}

func (m *RolePageModel) FetchData(client data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchRolePolicies(client, nil),
		m.fetchAssumeRolePolicy(client),
		m.fetchRoleTags(client, nil),
	)
}

func (m *RolePageModel) fetchRolePolicies(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.IAM.GetRolePolicies(m.Context.(RolePageContext).RoleName, nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Policies"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchRolePolicies(client, nextToken)
		}
		return msg
	}
}

func (m *RolePageModel) fetchAssumeRolePolicy(client data.Client) tea.Cmd {
	return func() tea.Msg {
		policy, _ := client.IAM.GetAssumeRolePolicy(m.Context.(RolePageContext).RoleName)
		msg := code.NewCodeContentMsg{
			Page:     m.Spec.Name,
			PaneId:   m.GetPaneId("Trust Relationships"),
			Content:  policy,
			Filepath: ".json",
		}
		return msg
	}
}

func (m *RolePageModel) fetchRoleTags(client data.Client, nextToken *string) tea.Cmd {
	return func() tea.Msg {
		rows, nextToken, _ := client.IAM.GetRoleTags(m.Context.(RolePageContext).RoleName, nextToken)
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Tags"),
			Rows:   rows,
		}
		if nextToken != nil {
			msg.NextCmd = m.fetchRoleTags(client, nextToken)
		}
		return msg
	}
}

func (m *RolePageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	row := table.GetCurrentRowMarshalled()
	switch m.GetCurrentPaneId() {
	case m.GetPaneId("Policies"):
		return func() tea.Msg {
			return page.ChangePageMsg{
				NewPage:   "iam/policy",
				FetchData: true,
				PageContext: PolicyPageContext{
					RoleName:   m.Context.(RolePageContext).RoleName,
					PolicyName: row["Name"],
					PolicyArn:  row["ARN"],
				},
			}
		}
	default:
		return nil
	}
}
