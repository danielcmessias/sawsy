package s3

import (
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/components/table"
	"github.com/danielcmessias/lfq/utils/icons"
)

var bucketsPageSpec = page.PageSpec{
    Name: "s3/buckets",
    TableSpecs: []table.TableSpec{
        {
            Name: "Buckets",
            Icon: icons.BUCKET,
            Columns: []table.Column{
                {
                    Title: "Name",
                },
                {
                    Title: "Region",
                },
                {
                    Title: "Creation Date",
                },
            },
        },
    },
}

var objectsPageSpec = page.PageSpec{
    Name: "s3/objects",
    TableSpecs: []table.TableSpec{
        {
            Name: "Objects",
            Icon: icons.FILES,
            Columns: []table.Column{
                {
                    Title: "Key",
                },
                {
                    Title: "Last Modified",
                },
                {
                    Title: "Size",
                },
            },
        },
    },
}
