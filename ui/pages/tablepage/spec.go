package tablepage

import (
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/components/table"
)

var databasesPageSpec = page.PageSpec{
    Name: "lf/tables",
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
            Name: "ﴳ Schema",
            Columns: []table.Column{
                {
                    Title: "#",
                },
                {
                    Title: "Column",
                },
                {
                    Title: "Data Type",
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
