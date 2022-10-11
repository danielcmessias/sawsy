package rds

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/data"
	"github.com/danielcmessias/sawsy/ui/components/chart/histogram"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/context"
)

type InstancePageModel struct {
	page.Model
}

type InstancePageContext struct {
	InstanceId string
}

func NewInstancePage(ctx *context.ProgramContext) *InstancePageModel {
	return &InstancePageModel{
		Model: page.New(ctx, instanceSpecPage),
	}
}

func (m *InstancePageModel) FetchData(client data.Client) tea.Cmd {
	cmds := []tea.Cmd{
		m.fetchDetails(client),
		m.fetchTags(client),
	}

	for i, met := range metrics {
		cmds = append(cmds, m.fetchMetric(client, i, met.APIName, met.Formatter))
	}
	return tea.Batch(cmds...)
}

func (m *InstancePageModel) fetchDetails(client data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.RDS.GetInstanceDetails(m.Context.(InstancePageContext).InstanceId)

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Details"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *InstancePageModel) fetchTags(client data.Client) tea.Cmd {
	return func() tea.Msg {
		rows, _ := client.RDS.GetInstanceTags(m.Context.(InstancePageContext).InstanceId)

		msg := page.NewRowsMsg{
			Page:   m.Spec.Name,
			PaneId: m.GetPaneId("Tags"),
			Rows:   rows,
		}
		return msg
	}
}

func (m *InstancePageModel) fetchMetric(client data.Client, galleryPaneId int, metric string, valueFormatter func(float64) float64) tea.Cmd {
	return func() tea.Msg {
		data, _ := client.RDS.GetMetric(m.Context.(InstancePageContext).InstanceId, metric)

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
