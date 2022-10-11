package services

import (
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils/icons"
)

var servicesPageSpec = page.PageSpec{
	Name: "services",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Services",
				Icon: icons.APPLICATION,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
			},
		},
	},
}
