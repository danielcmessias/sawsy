package glue

import (
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils/icons"
)

var gluePageSpec = page.PageSpec{
	Name: "glue",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Jobs",
				Icon: icons.TASKS,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
				{
					Title: "Type",
				},
				{
					Title: "Last Modified",
				},
				{
					Title: "Glue Version",
				},
				{
					Title: "Worker type",
				},
				{
					Title: "# Workers",
				},
			},
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Crawlers",
				Icon: icons.BUG,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
				{
					Title: "Schedule",
				},
				{
					Title: "Status",
				},
				{
					Title: "Last runtime",
				},
				{
					Title: "Median runtime",
				},
			},
		},
	},
}

var jobPageSpec = page.PageSpec{
	Name: "glue/job",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Details",
				Icon: icons.INFO,
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
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Runs",
				Icon: icons.TASKS,
			},
			Columns: []table.Column{
				{
					Title: "Start time",
				},
				{
					Title: "Run status",
				},
				{
					Title: "Attempt",
				},
				{
					Title: "Execution time",
				},
			},
		},
		code.CodeSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Script",
				Icon: icons.FILE_CODE,
			},
		},
	},
}
