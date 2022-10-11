package code

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/lexers"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	gansi "github.com/charmbracelet/glamour/ansi"
	"github.com/charmbracelet/lipgloss"
	"github.com/danielcmessias/sawsy/ui/components/pane"
	"github.com/danielcmessias/sawsy/ui/context"
	"github.com/muesli/termenv"
)

var (
	lineDigitStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("239"))
	lineBarStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("236"))
)

type Model struct {
	ctx            *context.ProgramContext
	viewport       viewport.Model
	renderContext  gansi.RenderContext
	showLineNumber bool
	rawContent     string
	filepath       string

	pane.Pane
}

type CodeSpec struct {
	pane.BaseSpec
}

type NewCodeContentMsg struct {
	Page     string
	PaneId   int
	Content  string
	Filepath string
}

func (s CodeSpec) NewFromSpec(ctx *context.ProgramContext, spec pane.PaneSpec) pane.Pane {
	codeSpec, ok := spec.(CodeSpec)
	if !ok {
		log.Fatal("invalid spec type, expected CodeSpec")
	}
	return New(ctx, codeSpec)
}

func New(ctx *context.ProgramContext, spec CodeSpec) *Model {
	return &Model{
		Pane: pane.New(spec.BaseSpec),

		ctx:      ctx,
		viewport: viewport.Model{},
		renderContext: gansi.NewRenderContext(gansi.Options{
			ColorProfile: termenv.TrueColor,
			Styles:       glamour.DraculaStyleConfig,
		}),
		showLineNumber: true,
	}
}

func (m *Model) Update(msg tea.Msg) (pane.Pane, tea.Cmd, bool) {
	var cmds []tea.Cmd
	vp, cmd := m.viewport.Update(msg)
	m.viewport = vp
	cmds = append(cmds, cmd)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.ctx.Keys.Down):
			m.viewport.LineDown(1)
			return m, nil, false
		case key.Matches(msg, m.ctx.Keys.Up):
			m.viewport.LineUp(1)
			return m, nil, false
		}
	}
	return m, tea.Batch(cmds...), false
}

func (m *Model) View() string {
	return m.viewport.View()
}

func (m *Model) SetSize(width int, height int) {
	m.viewport.Width = width
	m.viewport.Height = height
	m.SetContent(m.rawContent, m.filepath)
}

func (m *Model) Hide() {}

func (m *Model) SetContent(content string, filepath string) {
	m.rawContent = content
	m.filepath = filepath
	// Generiously 'borrowed' from soft-serve
	// https://github.com/charmbracelet/soft-serve/blob/main/ui/components/code/code.go#L185
	width := m.viewport.Width
	content = strings.ReplaceAll(content, "\t", strings.Repeat(" ", 4))
	var lexer chroma.Lexer
	if filepath != "" {
		lexer = lexers.Match(filepath)
	}
	if lexer == nil {
		// The analysis is broken and I have no idea why! :(
		lexer = lexers.Analyse(content)
	}
	if lexer == nil {
		lexer = lexers.Fallback
	}

	formatter := &gansi.CodeBlockElement{
		Code:     content,
		Language: lexer.Config().Name,
	}
	s := strings.Builder{}
	rc := m.renderContext
	if m.showLineNumber {
		st := glamour.DraculaStyleConfig
		var m uint
		st.CodeBlock.Margin = &m
		rc = gansi.NewRenderContext(gansi.Options{
			ColorProfile: termenv.TrueColor,
			Styles:       st,
		})
	}
	err := formatter.Render(&s, rc)
	if err != nil {
		log.Fatal(err)
	}
	c := s.String()
	if m.showLineNumber {
		var ml int
		c, ml = withLineNumber(c)
		width -= ml
	}

	m.viewport.SetContent(lipgloss.NewStyle().Width(width).Render(c))
}

func withLineNumber(s string) (string, int) {
	lines := strings.Split(s, "\n")
	// NB: len() is not a particularly safe way to count string width (because
	// it's counting bytes instead of runes) but in this case it's okay
	// because we're only dealing with digits, which are one byte each.
	mll := len(fmt.Sprintf("%d", len(lines)))
	for i, l := range lines {
		digit := fmt.Sprintf("%*d", mll, i+1)
		bar := "│"
		digit = lineDigitStyle.Render(digit)
		bar = lineBarStyle.Render(bar)
		if i < len(lines)-1 || len(l) != 0 {
			// If the final line was a newline we'll get an empty string for
			// the final line, so drop the newline altogether.
			lines[i] = fmt.Sprintf(" %s %s %s", digit, bar, l)
		}
	}
	return strings.Join(lines, "\n"), mll
}
