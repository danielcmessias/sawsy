package help

import (
	bbHelp "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/ui/context"
	"github.com/danielcmessias/sawsy/ui/styles"
	"github.com/danielcmessias/sawsy/utils"
)

type Model struct {
	ctx  *context.ProgramContext
	help bbHelp.Model
}

func NewModel(ctx *context.ProgramContext) Model {
	help := bbHelp.NewModel()
	help.Styles = bbHelp.Styles{
		ShortDesc:      helpTextStyle.Copy(),
		FullDesc:       helpTextStyle.Copy(),
		ShortSeparator: helpTextStyle.Copy(),
		FullSeparator:  helpTextStyle.Copy(),
		FullKey:        helpTextStyle.Copy(),
		ShortKey:       helpTextStyle.Copy(),
		Ellipsis:       helpTextStyle.Copy(),
	}

	return Model{
		ctx:  ctx,
		help: help,
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, utils.Keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}
	}

	return m, nil
}

func (m Model) View() string {
	return styles.FooterStyle.Copy().
		Width(m.ctx.ScreenWidth).
		Render(m.help.View(utils.Keys))
}

func (m *Model) SetWidth(width int) {
	m.help.Width = width
}
