package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/lfq/data"
	"github.com/danielcmessias/lfq/ui"
)

func main() {
	c := data.NewClient()
	c.FetchTableRows(nil)

	m := ui.NewModel()
    p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
