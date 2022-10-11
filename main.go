package main

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/danielcmessias/sawsy/config"
	"github.com/danielcmessias/sawsy/ui"
)

func main() {
	firstPage := "services"
	args := os.Args[1:]
	if len(args) > 0 {
		firstPage = args[0]
	}

	config, _ := config.ReadConfig()

	m := ui.NewModel(config, firstPage)
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),
	)
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}
