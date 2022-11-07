package lambda

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/chart/histogram"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/context"
)

type FunctionPageModel struct {
	page.Model
}

type FunctionPageContext struct {
	FunctionName string
}

func NewFunctionPage(ctx *context.ProgramContext) *FunctionPageModel {
	return &FunctionPageModel{
		Model: page.New(ctx, functionPageSpec),
	}
}

func (m *FunctionPageModel) FetchData(client *data.Client) tea.Cmd {
	cmds := []tea.Cmd{
		m.fetchDetails(client),
	}
	for i, met := range metrics {
		cmds = append(cmds, m.fetchMetric(client, i, met.APIName, met.Statistic, met.Formatter))
	}
	return tea.Batch(cmds...)
}

func (m *FunctionPageModel) fetchDetails(client *data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.Lambda.GetFunctionDetails(m.Context.(FunctionPageContext).FunctionName)

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Details"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *FunctionPageModel) fetchMetric(client *data.Client, galleryPaneId int, metric string, statistic types.Statistic, valueFormatter func(float64) float64) tea.Cmd {
	return func() tea.Msg {
		print("test")
		data, _ := client.Lambda.GetMetric("sherlock-decryptor-production", metric, statistic)

		if valueFormatter != nil {
			for i, d := range data {
				data[i] = valueFormatter(d)
			}
		}

		msg := histogram.NewDataMsg{
			Page:          m.Spec.Name,
			PaneId:        m.GetPaneId("Monitoring"),
			GalleryPaneId: galleryPaneId,
			Data:          data,
		}
		return msg
	}
}
