package databasepage

import (
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/components/table"
)

var databasesPageSpec = page.PageSpec{
    Name: "lf/databases",
    TableSpecs: []table.TableSpec{
        {
			Name: "Details",
			Icon: "",
			Columns: []table.Column{
				{
					Title: "Field",
				},
				{
					Title: "Value",
				},
			},
		},
        {
			Name: " LF-Tags",
			Icon: "",
			Columns: []table.Column{
				{
					Title: "Key",
				},
				{
					Title: "Value",
				},
			},
		},
    },
}
