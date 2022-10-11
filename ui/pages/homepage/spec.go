package homepage

import (
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/components/table"
	"github.com/danielcmessias/lfq/utils"
)

var lakeFormationPageSpec = page.PageSpec{
    Name: "lf/home",
    TableSpecs: []table.TableSpec{
        {
            Name: "Databases",
            Icon: "",
            Columns: []table.Column{
                {
                    Title: "Database",
                },
                {
                    Title: "CatalogId",
                },
                {
                    Title: "Description",
                },
                {
                    Title: "S3 Path",
                },
                {
                    Title: "Shared Resource",
                },
                {
                    Title: "Shared CatalogId",
                },
            },
            InspectRowPage: "lf/databases",
        },
        {
            Name: "Tables",
            Icon: "",
            Columns: []table.Column{
                {
                    Title: "Table",
                    MaxWidth: utils.IntPtr(128),
                },
                {
                    Title: "Database",
                },
                {
                    Title: "CatalogId",
                },
                {
                    Title: "Description",
                },
                {
                    Title: "S3 Path",
                },
                {
                    Title: "Shared Resource",
                },
                {
                    Title: "Shared CatalogId",
                },
            },
            InspectRowPage: "lf/tables",
        },
        {
            Name: "LF-Tags",
            Icon: "",
            Columns: []table.Column{
                {
                    Title: "Key",
                },
                {
                    Title: "Values",
                },
                {
                    Title: "Owner account ID",
                },
            },
        },
        {
            Name: "LF-Tag Perms",
            Icon: "ﰠ",
            Columns: []table.Column{
                {
                    Title: "Principal",
                },
                {
                    Title: "Key",
                },
                {
                    Title: "Values",
                },
                {
                    Title: "Permissions",
                },
                {
                    Title: "Grantable",
                }, 
            },
        },
        {
            Name: "LF Locations",
            Icon: "",
            Columns: []table.Column{
                {
                    Title: "S3 Path",
                },
                {
                    Title: "Last Modified",
                },
            },
        },
    },
}
