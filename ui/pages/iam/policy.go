package iam

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/context"
)

type PolicyPageModel struct {
	page.Model
}

type PolicyPageContext struct {
	// Managed Policies require only the Policy ARN
	PolicyArn string
	// Inlines have the PolicyName + (UserName or RoleName)
	PolicyName string
	UserName   string
	RoleName   string
}

func NewPolicyPage(ctx *context.ProgramContext) *PolicyPageModel {
	return &PolicyPageModel{
		Model: page.New(ctx, policyPageSpec),
	}
}

func (m *PolicyPageModel) FetchData(client *data.Client) tea.Cmd {
	return tea.Batch(
		m.fetchPolicyPermissions(client),
	)
}

func (m *PolicyPageModel) fetchPolicyPermissions(client *data.Client) tea.Cmd {
	return func() tea.Msg {
		context := m.Context.(PolicyPageContext)

		var policy string
		var err error

		if context.PolicyArn != "" {
			policy, err = client.IAM.GetManagedPolicy(context.PolicyArn)
		} else {
			if context.UserName != "" {
				policy, err = client.IAM.GetInlineUserPolicy(context.UserName, context.PolicyName)
			} else {
				policy, err = client.IAM.GetInlineRolePolicy(context.RoleName, context.PolicyName)
			}
		}
		if err != nil {
			log.Fatal(err)
		}

		msg := code.NewCodeContentMsg{
			Page:     m.Spec.Name,
			PaneId:   m.GetPaneId("Permissions"),
			Content:  policy,
			Filepath: ".json",
		}
		return msg
	}
}
