package lambda

import (
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/danielcmessias/sawsy/ui/components/chart/histogram"
	"github.com/danielcmessias/sawsy/ui/components/gallery"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils/icons"
)

var lambdaPageSpec = page.PageSpec{
	Name: "lambda",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Functions",
				Icon: icons.LAMBDA,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
				{
					Title: "Runtime",
				},
				{
					Title: "Last Modified",
				},
				{
					Title: "Description",
				},
				{
					Title: "ARN",
				},
			},
		},
	},
}

type metric struct {
	APIName   string
	HumanName string
	Statistic types.Statistic
	Formatter func(float64) float64
}

var metrics = []metric{
	{"Invocations", "Invocations", types.StatisticSum, nil},
	{"Duration", "Duration", types.StatisticAverage, nil},
}

var metricSpecs = func() []pane.PaneSpec {
	var arr []pane.PaneSpec = make([]pane.PaneSpec, len(metrics))
	for i, m := range metrics {
		arr[i] = histogram.HistogramSpec{
			BaseSpec: pane.BaseSpec{
				Name: m.HumanName,
			},
		}
	}
	return arr
}()

var functionPageSpec = page.PageSpec{
	Name: "lambda/function",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Details",
				Icon: icons.LAMBDA,
			},
			Columns: []table.Column{
				{
					Title: "Key",
				},
				{
					Title: "Value",
				},
			},
		},
		gallery.GallerySpec{
			BaseSpec: pane.BaseSpec{
				Name: "Monitoring",
				Icon: icons.CHART,
			},
			Rows:      3,
			Cols:      3,
			PaneSpecs: metricSpecs,
		},
	},
}
