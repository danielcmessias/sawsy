package s3

import (
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils/icons"
)

var bucketsPageSpec = page.PageSpec{
	Name: "s3",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Buckets",
				Icon: icons.BUCKET,
			},
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

var bucketPageSpec = page.PageSpec{
	Name: "s3/objects",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Objects",
				Icon: icons.FILES,
			},
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
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Properties",
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
		code.CodeSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Bucket Policy",
				Icon: icons.KEY,
			},
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Tags",
				Icon: icons.TAG,
			},
			Columns: []table.Column{
				{
					Title: "Key",
				},
				{
					Title: "Values",
				},
			},
		},
	},
}

var objectPageSpec = page.PageSpec{
	Name: "s3/object",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Properties",
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
				Name: "Permissions",
				Icon: icons.KEY,
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
