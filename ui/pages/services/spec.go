package services

import (
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/components/table"
)

var servicesPageSpec = page.PageSpec{
    Name: "services",
    TableSpecs: []table.TableSpec{
        {
            Name: "Services",
            Icon: "ﬓ",
            Columns: []table.Column{
                {
                    Title: "Name",
                },
            },
        },
    },
}
