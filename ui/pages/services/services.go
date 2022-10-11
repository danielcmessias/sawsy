package services

import (
	"log"
	"sort"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/ui/context"
)

var services = map[string]string{
	"Glue":           "glue",
	"IAM":            "iam",
	"Lake Formation": "lakeformation",
	"Lambda":         "lambda",
	"RDS":            "rds",
	"S3":             "s3",
}

type ServicesPageModel struct {
	page.Model
}

func NewServicesPage(ctx *context.ProgramContext) *ServicesPageModel {
	return &ServicesPageModel{
		Model: page.New(ctx, servicesPageSpec),
	}
}

func (m *ServicesPageModel) FetchData(client data.Client) tea.Cmd {
	var rows []table.Row
	var l []string
	for s := range services {
		l = append(l, s)
	}
	sort.Strings(l)
	for _, s := range l {
		rows = append(rows, table.Row{s})
	}

	return func() tea.Msg {
		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: 0,
			Rows:   rows,
		}
		return msg
	}
}

func (m *ServicesPageModel) Inspect(client data.Client) tea.Cmd {
	table, ok := m.CurrentPane().(*table.Model)
	if !ok {
		log.Fatal("This pane is not a table")
	}

	nextPage := services[table.GetCurrentRow()[0]]

	return func() tea.Msg {
		return page.ChangePageMsg{
			NewPage:   nextPage,
			FetchData: true,
		}
	}
}
