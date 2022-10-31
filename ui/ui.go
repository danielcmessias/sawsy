package ui

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/config"
	"github.com/danielcmessias/sawsy/data"

	"github.com/danielcmessias/sawsy/ui/components/help"
	"github.com/danielcmessias/sawsy/ui/components/page"
	"github.com/danielcmessias/sawsy/ui/context"
	"github.com/danielcmessias/sawsy/ui/pages/glue"
	"github.com/danielcmessias/sawsy/ui/pages/iam"
	"github.com/danielcmessias/sawsy/ui/pages/lakeformation"
	"github.com/danielcmessias/sawsy/ui/pages/lambda"
	"github.com/danielcmessias/sawsy/ui/pages/rds"
	"github.com/danielcmessias/sawsy/ui/pages/s3"
	"github.com/danielcmessias/sawsy/ui/pages/services"
	"github.com/danielcmessias/sawsy/utils"
)

type Model struct {
	config config.Config
	client data.Client
	ctx    *context.ProgramContext
	help   help.Model
	keys   utils.KeyMap

	pages        map[string]page.Page
	currentPage  string
	visitedPages []PageVisit
}

type PageVisit struct {
	PageName string
	Context  interface{}
}

func NewModel(config config.Config, firstPage string) Model {
	client := data.NewClient()
	ctx := &context.ProgramContext{
		Config:       &config,
		AwsAccountId: client.GetCurrentAWSAccountId(),
		AwsService:   firstPage,
		Keys:         utils.Keys,
	}

	var allPages = []page.Page{
		services.NewServicesPage(ctx),
		glue.NewGluePage(ctx),
		glue.NewJobsPage(ctx),
		iam.NewIAMPage(ctx),
		iam.NewUserPage(ctx),
		iam.NewRolePage(ctx),
		iam.NewPolicyPage(ctx),
		rds.NewRDSPage(ctx),
		rds.NewInstancePage(ctx),
		lakeformation.NewLakeFormationPage(ctx),
		lakeformation.NewDatabasePage(ctx),
		lakeformation.NewTablePage(ctx),
		lambda.NewLambdaPage(ctx),
		lambda.NewFunctionPage(ctx),
		s3.NewS3Page(ctx),
		s3.NewBucketPage(ctx),
		s3.NewObjectPage(ctx),
	}
	pages := map[string]page.Page{}
	for _, p := range allPages {
		pages[p.GetSpec().Name] = p
	}

	return Model{
		config:      config,
		client:      client,
		ctx:         ctx,
		help:        help.NewModel(ctx),
		keys:        utils.Keys,
		pages:       pages,
		currentPage: firstPage,
	}
}

type initMsg struct{}

func initScreen() tea.Msg {
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

	cmd, consumed := m.getCurrentPage().Update(m.client, msg)
	cmds = append(cmds, cmd)
	if consumed {
		return m, tea.Batch(cmds...)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Inspect) && !m.ctx.LockKeyboardCapture:
			cmds = append(cmds, m.getCurrentPage().Inspect(m.client))

		case key.Matches(msg, m.keys.Services) && !m.ctx.LockKeyboardCapture:
			cmds = append(cmds, m.changePage("services", nil, true))

		case key.Matches(msg, m.keys.PrevPage) && !m.ctx.LockKeyboardCapture:
			l := len(m.visitedPages)
			if l == 0 {
				break
			}
			prev := m.visitedPages[l-1]
			m.visitedPages = m.visitedPages[:l-1]
			cmds = append(cmds, m.changePage(prev.PageName, prev.Context, true))

		case key.Matches(msg, m.keys.Refresh) && !m.ctx.LockKeyboardCapture:
			// Not implemented yet!

		case key.Matches(msg, m.keys.Quit):
			if !(m.ctx.LockKeyboardCapture && msg.String() == "q") {
				return m, tea.Quit
			}
		}

	case page.NewRowsMsg:
		cmds = append(cmds, m.parseNewRowsMsg(msg))

	case page.BatchedNewRowsMsg:
		for _, _msg := range msg.Msgs {
			cmds = append(cmds, m.parseNewRowsMsg(_msg))
		}

	case page.ChangePageMsg:
		m.visitedPages = append(m.visitedPages, PageVisit{
			PageName: m.currentPage,
			Context:  m.getCurrentPage().GetPageContext(),
		})
		cmds = append(cmds, m.changePage(msg.NewPage, msg.PageContext, msg.FetchData))

	case initMsg:
		cmds = append(cmds, m.getCurrentPage().FetchData(m.client))

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
		m.help.View(),
	)
}

func (m *Model) onWindowSizeChanged(msg tea.WindowSizeMsg) {
	m.ctx.ScreenWidth = msg.Width
	m.ctx.ScreenHeight = msg.Height
	for _, p := range m.pages {
		p.SetSize(msg.Width, msg.Height)
	}
	m.help.SetWidth(msg.Width)
}

func (m *Model) getCurrentPage() page.Page {
	if page, ok := m.pages[m.currentPage]; ok {
		return page
	} else {
		log.Fatalf("No page of name %s", m.currentPage)
		return nil
	}
}

func (m *Model) parseNewRowsMsg(msg page.NewRowsMsg) tea.Cmd {
	var cmds []tea.Cmd
	if msg.Overwrite {
		m.pages[msg.Page].ClearRows(msg.PaneId)
	}
	m.pages[msg.Page].AppendRows(msg.PaneId, msg.Rows)
	// Uncomment this to fetch ALL rows
	if msg.NextCmd != nil {
		cmds = append(cmds, msg.NextCmd)
	}
	return tea.Batch(cmds...)
}

func (m *Model) changePage(pageName string, context interface{}, fetchData bool) tea.Cmd {
	_, ok := m.pages[pageName]
	if !ok {
		log.Fatalf("No page with name %s", pageName)
	}

	m.currentPage = pageName
	m.getCurrentPage().SetSize(m.ctx.ScreenWidth, m.ctx.ScreenHeight)
	m.getCurrentPage().SetPageContext(context)

	m.ctx.AwsService = m.getCurrentPage().GetSpec().Name

	var cmd tea.Cmd
	if fetchData {
		m.getCurrentPage().ClearData()
		cmd = m.getCurrentPage().FetchData(m.client)
	}
	return cmd
}
