package pane

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/ui/context"
)

type Pane interface {
	Update(tea.Msg) (model Pane, cmd tea.Cmd, consumed bool)
	View() string
	GetSpec() PaneSpec
	SetSize(width int, height int)
	Hide() // Called when the pane is not visible
}

type PaneModel struct {
	Spec PaneSpec
}

type PaneSpec interface {
	NewFromSpec(ctx *context.ProgramContext, spec PaneSpec) Pane
	GetName() string
	GetIcon() string
}

type BaseSpec struct {
	Name string
	Icon string
}

func (s BaseSpec) NewFromSpec(ctx *context.ProgramContext, spec PaneSpec) Pane {
	return New(spec)
}

func New(spec PaneSpec) Pane {
	return PaneModel{
		Spec: spec,
	}
}

func (m PaneModel) Init() tea.Cmd {
	return nil
}

func (m PaneModel) Update(msg tea.Msg) (Pane, tea.Cmd, bool) {
	return m, nil, false
}

func (m PaneModel) View() string {
	return ""
}

func (m PaneModel) SetSize(width int, height int) {}

func (m PaneModel) Hide() {}

func (m PaneModel) GetSpec() PaneSpec {
	return m.Spec
}

func (s BaseSpec) GetName() string {
	return s.Name
}

func (s BaseSpec) GetIcon() string {
	return s.Icon
}
