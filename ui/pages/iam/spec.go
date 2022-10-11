package iam

import (
	"github.com/danielcmessias/sawsy/ui/components/code"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/components/table"
	"github.com/danielcmessias/sawsy/utils/icons"
)

var iamPageSpec = page.PageSpec{
	Name: "iam",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Users",
				Icon: icons.USER,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
				{
					Title: "ARN",
				},
				{
					Title: "Last Activity",
				},
				{
					Title: "Created On",
				},
			},
		},
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Roles",
				Icon: icons.USER_CIRCLE,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
				{
					Title: "ARN",
				},
				{
					Title: "Created On",
				},
			},
		},
	},
}

var userPageSpec = page.PageSpec{
	Name: "iam/user",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Policies",
				Icon: icons.LIST,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
				{
					Title: "Type",
				},
				{
					Title: "ARN",
				},
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
					Title: "Value",
				},
			},
		},
	},
}

var rolePageSpec = page.PageSpec{
	Name: "iam/role",
	PaneSpecs: []pane.PaneSpec{
		table.TableSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Policies",
				Icon: icons.LIST,
			},
			Columns: []table.Column{
				{
					Title: "Name",
				},
				{
					Title: "Type",
				},
				{
					Title: "ARN",
				},
			},
		},
		code.CodeSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Trust Relationships",
				Icon: icons.SHIELD,
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
					Title: "Value",
				},
			},
		},
	},
}

var policyPageSpec = page.PageSpec{
	Name: "iam/policy",
	PaneSpecs: []pane.PaneSpec{
		code.CodeSpec{
			BaseSpec: pane.BaseSpec{
				Name: "Permissions",
				Icon: icons.KEY,
			},
		},
	},
}
