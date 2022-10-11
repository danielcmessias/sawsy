package ui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui/components/help"
	"github.com/danielcmessias/lfq/ui/components/page"
	"github.com/danielcmessias/lfq/ui/context"
	"github.com/danielcmessias/lfq/ui/pages/databasepage"
	"github.com/danielcmessias/lfq/ui/pages/homepage"
	"github.com/danielcmessias/lfq/ui/pages/s3"
	"github.com/danielcmessias/lfq/ui/pages/services"
	"github.com/danielcmessias/lfq/ui/pages/tablepage"
	"github.com/danielcmessias/lfq/utils"
)

type Model struct {
    client        data.Client
    ctx           context.ProgramContext
    help          help.Model
    isSearching   bool
    keys          utils.KeyMap

    pages         map[string]page.Page
    currentPage   string
}


func NewModel() Model {
    client := data.NewClient()

    ctx := context.ProgramContext{
        // ScreenWidth: 171,
        // ScreenHeight: 24,
        AwsAccountId: client.FetchCurrentAWSAccountId(),
    }

    search := textinput.New()
    search.Placeholder = "..."
    search.Width = 40
    
    homePage := homepage.NewLakeFormationPage()
    databasePage := databasepage.NewDatabasesPage()
    tablePage := tablepage.NewTablePage()

    s3BucketsPage := s3.NewBucketsPage()
    s3ObjectsPage := s3.NewObjectsPage()


    pages := map[string]page.Page{
        "services": services.NewServicesPage(),
        "lf/home": &homePage,
        "lf/databases": &databasePage,
        "lf/tables": &tablePage,
        "s3/buckets": &s3BucketsPage,
        "s3/objects": &s3ObjectsPage,
    }

    return Model {
        client:        client,
        ctx:           ctx,
        help:          help.NewModel(),
        keys:          utils.Keys,

        pages:         pages, 
        currentPage:   "s3/buckets",
    }
}

type initMsg struct {
    // Config config.Config
}

func initScreen() tea.Msg {
    // config, err := config.ParseConfig()
    // if err != nil {
    // 	return errMsg{err}
    // }
    return initMsg{}
}

func (m Model) Init() tea.Cmd {
    cmds := []tea.Cmd{
        initScreen,
        tea.EnterAltScreen,
    }
    for _, p := range m.pages {
        cmds = append(cmds, p.Init())
    }
    return tea.Batch(cmds...)
}


func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var cmds []tea.Cmd

    for _, p := range m.pages {
        cmds = append(cmds, p.UpdateSearch(msg))
    }

    switch msg := msg.(type) {
        case tea.KeyMsg:
            switch {
                case key.Matches(msg, m.keys.Down):
                    m.getCurrentPage().NextItem()

                case key.Matches(msg, m.keys.Up):
                    m.getCurrentPage().PrevItem()

                case key.Matches(msg, m.keys.FirstLine) && !m.isSearching:
                    m.getCurrentPage().FirstItem()
        
                case key.Matches(msg, m.keys.LastLine) && !m.isSearching:
                    m.getCurrentPage().LastItem()

                case key.Matches(msg, m.keys.NextTab) && !m.isSearching:
                    m.getCurrentPage().NextTab()

                case key.Matches(msg, m.keys.PrevTab) && !m.isSearching:
                    m.getCurrentPage().PrevTab()

                case key.Matches(msg, m.keys.StartSearch) && !m.isSearching:
                    m.getCurrentPage().StartSearch()
                    m.isSearching = true

                case key.Matches(msg, m.keys.EndSearch) && m.isSearching:
                    m.getCurrentPage().StopSearch()
                    m.isSearching = false

                case key.Matches(msg, m.keys.InspectField) && !m.isSearching:
                    cmds = append(cmds, m.getCurrentPage().InspectFieldCmd(m.client))
                    cmds = append(cmds, m.getCurrentPage().InspectRow(m.client))
                    
                case key.Matches(msg, m.keys.Home) && !m.isSearching:
                    // m.currentPage = "lf/home"
                    m.currentPage = "s3/buckets"

                case key.Matches(msg, m.keys.Refresh) && !m.isSearching:
                    // m.getCurrentTableModel().ClearRows()
                    // m.getCurrentPage().Refresh()


                case key.Matches(msg, m.keys.Quit):
                    if (!(m.isSearching && msg.String() == "q")) {
                        return m, tea.Quit
                    }

                case key.Matches(msg, m.keys.NextCol):
                    m.getCurrentPage().NextCol()

                case key.Matches(msg, m.keys.PrevCol):
                    m.getCurrentPage().PrevCol()
                
            }

        case page.NewRowsMsg:
            cmds = append(cmds, m.parseNewRowsMsg(msg))

        case page.BatchedNewRowsMsg:
            for _, _msg := range msg.Msgs {
                cmds = append(cmds, m.parseNewRowsMsg(_msg))
            }

        case page.ChangePageMsg:
            m.currentPage = msg.NewPage
            if msg.PageMetadata != nil {
                m.getCurrentPage().SetPageMetadata(msg.PageMetadata)
            }
            if msg.InspectedRow != nil {
                m.getCurrentPage().FromInspectedRow(msg.InspectedRow)
            }

            if msg.FetchData {
                m.getCurrentPage().ClearData()
                cmds = append(cmds, m.getCurrentPage().FetchDataCmd(m.client))
            }
            // This whole ProgramContext needs to be reworked
            m.pages[m.currentPage].UpdateProgramContext(&m.ctx)

        case initMsg:
            cmds = append(cmds, m.getCurrentPage().FetchDataCmd(m.client))
            
        case tea.WindowSizeMsg:
            m.onWindowSizeChanged(msg)
    }

    var helpCmd tea.Cmd
    m.help, helpCmd = m.help.Update(msg)
    cmds = append(cmds, helpCmd)

    return m, tea.Batch(cmds...)
}

func (m Model) View() string {
    return lipgloss.JoinVertical(
		lipgloss.Left,
		m.getCurrentPage().View(),
		m.help.View(m.ctx),
	)
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
    m.ctx.ScreenWidth = msg.Width
    m.ctx.ScreenHeight = msg.Height
    for _, p := range m.pages {
        p.UpdateProgramContext(&m.ctx)
    }
    m.help.SetWidth(msg.Width)
}

func (m *Model) getCurrentPage() page.Page {
    return m.pages[m.currentPage]
}

func (m *Model) parseNewRowsMsg(msg page.NewRowsMsg) tea.Cmd {
    var cmds []tea.Cmd
    if msg.Overwrite {
        m.pages[msg.Page].ClearRows(msg.TabId)
    }
    m.pages[msg.Page].AppendRows(msg.TabId, msg.Rows)
    // Uncomment this to fetch ALL rows
    if (msg.NextCmd != nil) {
        cmds = append(cmds, msg.NextCmd)
    }
    return tea.Batch(cmds...)
}
