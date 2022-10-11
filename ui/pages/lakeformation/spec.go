package lakeformation

import (
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils"
	"github.com/danielcmessias/sawsy/utils/icons"
)

var lakeFormationPageSpec = page.PageSpec{
	Name: "lakeformation",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Databases",
				Icon: icons.DATABASE,
			},
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
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Tables",
				Icon: icons.TABLE,
			},
			Columns: []table.Column{
				{
					Title:    "Table",
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
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "LF-Tags",
				Icon: icons.TAG,
			},
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
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "LF-Tag Perms",
				Icon: icons.KEY,
			},
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
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "LF Locations",
				Icon: icons.LOCATION,
			},
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

var databasePageSpec = page.PageSpec{
	Name: "lakeformation/databases",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Details",
				Icon: icons.INFO,
			},
			Columns: []table.Column{
				{
					Title: "Field",
				},
				{
					Title: "Value",
				},
			},
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "LF-Tags",
				Icon: icons.TAG,
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
	},
}

var tablePageSpec = page.PageSpec{
	Name: "lakeformation/tables",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Details",
				Icon: icons.INFO,
			},
			Columns: []table.Column{
				{
					Title: "Field",
				},
				{
					Title: "Value",
				},
			},
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Schema",
				Icon: icons.SCHEMA,
			},
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
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "LF-Tags",
				Icon: icons.TAG,
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
	},
}
