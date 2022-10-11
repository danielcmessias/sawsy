package homepage

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/pages/databasepage"
	"github.com/danielcmessias/lfq/ui/pages/tablepage"
)

const (
    TAB_DATABASES          = iota
    TAB_TABLES             = iota
    TAB_LF_TAGS            = iota
    TAB_LF_TAG_PERMISSIONS = iota
    TAB_LOCATIONS          = iota
)

type LakeFormationPageModel struct {
    page.Model
}

type HomePage interface {
    page.Page
}

// func NewModel(id int) Model {

//     databasesTbl := table.NewModel(
//         databaseColSpec,
//         nil,
//         "Databases",
//     )
//     tablesTbl := table.NewModel(
//         tableColSpec,
//         nil,
//         "Tables",
//     )
//     lfTagsTbl := table.NewModel(
//         lfTagColSpec,
//         nil,
//         "LF-Tags",
//     )
//     lfTagsPermsTbl := table.NewModel(
//         lfTagPermColSpec,
//         nil,
//         "LF-Tag Permissions",
//     )

//     locationsTbl := table.NewModel(
//         locationsColSpec,
//         nil,
//         "LF-Tags",
//     )

//     return Model{
//         Model: page.Model{
//             Id:     id,
//             Tabs:   tabs.NewModel([]string{" Databases", " Tables", " LF-Tags", "ﰠ LF-Tag Perms", " LF Locations"}),
//             Tables: []table.Model{databasesTbl, tablesTbl, lfTagsTbl, lfTagsPermsTbl, locationsTbl},
//         },
//     }
// }

func NewLakeFormationPage() LakeFormationPageModel {
	return LakeFormationPageModel{
		Model: page.New(lakeFormationPageSpec),
	}
}

func (m *LakeFormationPageModel) FetchDataCmd(client data.Client) tea.Cmd {
    return tea.Batch(
        m.fetchDatabasesCmd(client, nil),
        m.fetchTablesCmd(client, nil),
        m.fetchLFTagsCmd(client, nil),
        m.fetchLFTagPermissionsCmd(client, nil),
        m.fetchDataLakeLocationsCmd(client, nil),
    )
}

func (m *LakeFormationPageModel) InspectFieldCmd(client data.Client) tea.Cmd {
    currRow := m.Tables[m.Tabs.CurrentTabId].GetCurrRow()

    var nextPage string
    var pageMetadata interface{}

    switch m.Tabs.CurrentTabId {
    case TAB_DATABASES:
        nextPage = "lf/databases"
        pageMetadata = databasepage.PageMetadata{
            DatabaseName: currRow[0],
        }
    case TAB_TABLES:
        nextPage = "lf/tables"
        pageMetadata = tablepage.PageMetadata{
            TableName:    currRow[0],
            DatabaseName: currRow[1],
        }
    case TAB_LF_TAGS:
        // Not implemented
        return nil
    }

    changePageCmd := func() tea.Msg {
        return page.ChangePageMsg{
            NewPage: nextPage,
            FetchData: true,
            PageMetadata: pageMetadata,
        }
    }
    return changePageCmd
}

func (m *LakeFormationPageModel) fetchDatabasesCmd(client data.Client, nextToken *string) tea.Cmd {
    return func() tea.Msg {
        rows, nextToken, _ := client.FetchDatabaseRows(nextToken)
        msg := page.NewRowsMsg{
            Page: m.Spec.Name,
            TabId: TAB_DATABASES,
            Rows: rows,
        }
        if (nextToken != nil) {
            msg.NextCmd = m.fetchDatabasesCmd(client, nextToken)
        }
        return msg
    }
}

func (m *LakeFormationPageModel) fetchTablesCmd(client data.Client, nextToken *string) tea.Cmd {
    return func() tea.Msg {
        rows, nextToken, _ := client.FetchTableRows(nextToken)
        msg := page.NewRowsMsg{
            Page: m.Spec.Name,
            TabId: TAB_TABLES,
            Rows: rows,
        }
        if (nextToken != nil) {
            msg.NextCmd = m.fetchTablesCmd(client, nextToken)
        }
        return msg
    }
}

func (m *LakeFormationPageModel) fetchLFTagsCmd(client data.Client, nextToken *string) tea.Cmd {
    return func() tea.Msg {
        rows, nextToken, _ := client.FetchLFTagRows(nextToken)
        msg := page.NewRowsMsg{
            Page: m.Spec.Name,
            TabId: TAB_LF_TAGS,
            Rows: rows,
        }
        if (nextToken != nil) {
            msg.NextCmd = m.fetchLFTagsCmd(client, nextToken)
        }
        return msg
    }
}

func (m *LakeFormationPageModel) fetchLFTagPermissionsCmd(client data.Client, nextToken *string) tea.Cmd {
    return func() tea.Msg {
        rows, nextToken, _ := client.FetchLFTagPermissionRows(nextToken)
        msg := page.NewRowsMsg{
            Page: m.Spec.Name,
            TabId: TAB_LF_TAG_PERMISSIONS,
            Rows: rows,
        }
        if (nextToken != nil) {
            msg.NextCmd = m.fetchLFTagPermissionsCmd(client, nextToken)
        }
        return msg
    }
}

func (m *LakeFormationPageModel) fetchDataLakeLocationsCmd(client data.Client, nextToken *string) tea.Cmd {
    return func() tea.Msg {
        rows, nextToken, _ := client.FetchDataLakeLocationRows(nextToken)
        msg := page.NewRowsMsg{
            Page: m.Spec.Name,
            TabId: TAB_LOCATIONS,
            Rows: rows,
        }
        if (nextToken != nil) {
            msg.NextCmd = m.fetchDataLakeLocationsCmd(client, nextToken)
        }
        return msg
    }
}
